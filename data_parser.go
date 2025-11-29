package go_dlt645_2007

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"strconv"
)

// LengthMismatchError 长度不匹配
var LengthMismatchError = errors.New("dlt645_2007 meterDataParser: length mismatch")

// NewMeterDataParser 创建一个数据解析器
func NewMeterDataParser(size int, order binary.ByteOrder, ratio, offset float64, uint string) (*MeterDataParser, error) {
	if size <= 0 {
		return nil, errors.New("dlt645_2007 meterDataParser: size must be positive")
	}
	if ratio == 0 {
		return nil, errors.New("ratio must be is not zero")
	}
	if order == nil {
		order = binary.LittleEndian
	}
	return &MeterDataParser{size: size, order: order, ratio: ratio, offset: offset, unit: uint}, nil
}

// MeterDataParser 数据解析器
type MeterDataParser struct {
	size   int              //数据长度
	order  binary.ByteOrder //大小端序， 默认小端序
	ratio  float64          //倍率
	offset float64          //偏移量，偏移量是减法运算
	unit   string           //单位
	data   []float64
}

// ObtainValue 获取解析结果
func (p *MeterDataParser) ObtainValue() float64 {
	return p.data[0]
}

// ObtainValues 获取解析结果
func (p *MeterDataParser) ObtainValues() []float64 {
	return p.data
}

func (p *MeterDataParser) decode(data []byte) error {
	if len(data)%p.size != 0 {
		return LengthMismatchError
	}
	for i := 0; i < len(data); i += p.size {
		end := i + p.size
		val := data[i:end]
		value, err := p.parser(val)
		if err != nil {
			return err
		}
		p.data = append(p.data, value*p.ratio-p.offset)
	}
	return nil
}

func (p *MeterDataParser) parser(data []byte) (float64, error) {
	if len(data) != p.size {
		return 0, LengthMismatchError
	}
	original := hex.EncodeToString(reverseBytes(data))
	if p.order == binary.BigEndian {
		original = hex.EncodeToString(data)
	}
	val, err := strconv.Atoi(original)
	if err != nil {
		return 0, err
	}
	return float64(val), nil
}
