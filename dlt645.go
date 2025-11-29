package go_dlt645_2007

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
)

const (
	dlt645StartChar  byte   = 0x68           //帧起始符
	dlt645EndChar    byte   = 0x16           //结束符
	BroadcastAddress string = "999999999999" //广播地址
	substitute       string = "AA"           //补位
	disturb          byte   = 0x33
)

type MeterDlt645Protocol struct {
	prefix string //通配前缀
	//startChar1  byte   //帧起始符
	Address string //地址域, 低字节在前，高字节在后
	//startChar2  byte   //帧起始符
	ControlChar byte   //控制码
	Length      byte   //数据域长度,读数据时 L≤200，写数据时 L≤50，L=0 表示无数据域
	Data        []byte //数据域,传输时发送方按字节进行加33H 处理，接收方按字节进行减33H 处理
	//cs          byte   //校验码，从第一个帧起始符开始到校验码之前的所有各字节的模256的和
	endChar  byte //结束符
	original []byte
}

// Decode 解码
func (m *MeterDlt645Protocol) Decode(frame []byte) error {
	if len(frame) < 12 {
		return errors.New("not enough frame length")
	}
	buf := bufio.NewReader(bytes.NewBuffer(frame))
	return m.DecodeByBuf(buf)
}

func (m *MeterDlt645Protocol) DecodeByBuf(buf *bufio.Reader) error {
	m.original = nil
	var startChar byte
	for {
		err := binary.Read(buf, binary.BigEndian, &startChar)
		if err != nil {
			return err
		}
		if startChar != dlt645StartChar {
			continue
		}
		break
	}
	//帧起始符
	snap := []byte{dlt645StartChar}
	//地址域, 低字节在前，高字节在后
	var address = make([]byte, 6)
	err := binary.Read(buf, binary.BigEndian, &address)
	if err != nil {
		return err
	}
	snap = append(snap, address...)
	m.Address = hex.EncodeToString(reverseBytes(address))
	//帧起始符
	err = binary.Read(buf, binary.BigEndian, &startChar)
	if err != nil {
		return err
	}
	if startChar != dlt645StartChar {
		return fmt.Errorf("start char 2 != 68H")
	}
	//控制码
	err = binary.Read(buf, binary.BigEndian, &m.ControlChar)
	if err != nil {
		return err
	}
	//数据域长度,读数据时 L≤200，写数据时 L≤50，L=0 表示无数据域
	err = binary.Read(buf, binary.BigEndian, &m.Length)
	if err != nil {
		return err
	}
	snap = append(snap, dlt645StartChar, m.ControlChar, m.Length)
	//数据域,传输时发送方按字节进行加33H 处理，接收方按字节进行减33H 处理
	if m.Length > 0 {
		m.Data = make([]byte, m.Length)
		err = binary.Read(buf, binary.BigEndian, &m.Data)
		if err != nil {
			return err
		}
		snap = append(snap, m.Data...)
		for i, b := range m.Data {
			m.Data[i] = b - disturb
		}
	}
	//校验码
	var cs byte
	err = binary.Read(buf, binary.BigEndian, &cs)
	if err != nil {
		return err
	}
	//计算校验码
	if cs != m.cs(snap) {
		return errors.New("cs error")
	}
	//结束符
	var endChar byte
	err = binary.Read(buf, binary.BigEndian, &endChar)
	if err != nil {
		return err
	}
	if endChar != dlt645EndChar {
		return fmt.Errorf("end char != 16H")
	}
	m.original = append(snap, cs, endChar)
	return nil
}

func (m *MeterDlt645Protocol) cs(frame []byte) byte {
	var sum uint8
	for _, b := range frame {
		sum += b
	}
	return sum
}

func (m *MeterDlt645Protocol) Encode() ([]byte, error) {
	frame := []byte{dlt645StartChar}
	if len(m.Address) > 12 {
		return nil, fmt.Errorf("address too long")
	} else if len(m.Address) == 0 {
		return nil, fmt.Errorf("address is empty")
	} else {
		for len(m.Address) < 12 {
			m.Address = "0" + m.Address
		}
	}
	addr, err := hex.DecodeString(m.Address)
	if err != nil {
		return nil, errors.New("invalid address")
	}
	frame = append(frame, reverseBytes(addr)...)
	if m.Data == nil {
		frame = append(frame, dlt645StartChar, m.ControlChar, 0x00)
	} else {
		frame = append(frame, dlt645StartChar, m.ControlChar, byte(len(m.Data)))
		for i, b := range m.Data {
			m.Data[i] = b + disturb
		}
		frame = append(frame, m.Data...)
	}
	//计算cs
	frame = append(frame, m.cs(frame), dlt645EndChar)
	if strings.TrimSpace(m.prefix) != "" {
		pf, pfErr := hex.DecodeString(m.prefix)
		if pfErr != nil {
			return nil, err
		}
		frame = append(pf, frame...)
	}
	return frame, nil
}

// ObtainDataLen 获取数据域的长度
func (m *MeterDlt645Protocol) ObtainDataLen() int {
	return len(m.Data) + 1
}
