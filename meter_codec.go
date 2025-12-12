package go_dlt645_2007

import (
	"encoding/hex"
	"errors"
)

// FuncCodeError 功能码错误或未实现这个功能码的逻辑
var FuncCodeError = errors.New("dlt645_2007: function code error or this program does not support this function code")

// DataDomainError 数据域错误
var DataDomainError = errors.New("dlt645_2007: Data domain error")

type MeterDataReceiver interface {
	// MeterReadResponse 电表正确应答的数据 ident-数据标识，parser解析的结果，hasNext是否存在后续帧, seq-帧序号,0标识最开始的帧
	MeterReadResponse(ident []byte, parser *MeterDataParser, hasNext bool, seq byte)
	MeterDefaultReadResponse(funcCode byte, data []byte) //MeterReadResponse 找不到注册器就会到这里
	// MeterReadErrorResponse 读数据后电表的异常应答，reqFrame-请求的报文， funcCode-控制码，errCode-错误信息字
	MeterReadErrorResponse(funcCode byte, errCode byte)
	MeterReqMasterSet(isSuccess bool, errCode byte)     //设置电表后的回复
	MeterAddress(addr string)                           //读电表地址的回复
	FreezeCommandResponse(isSuccess bool, errCode byte) //冻结命令回复
	// ErrorData 解析失败的数据会调用这个方法 funcCode-控制码， data-数据域， err-错误类型
	ErrorData(funcCode byte, data []byte, err error)
}

func NewMeterDataCodec(receiver MeterDataReceiver) *MeterDataCodec {
	return &MeterDataCodec{receiver: receiver, parsers: make(map[string]*MeterDataParser)}
}

// MeterDataCodec 数据解析器
type MeterDataCodec struct {
	receiver MeterDataReceiver
	parsers  map[string]*MeterDataParser
}

// Register 注册数据解析器
// ident 数据标识
// parser 数据解析器
func (m *MeterDataCodec) Register(ident []byte, parser *MeterDataParser) {
	m.parsers[hex.EncodeToString(reverseBytes(ident))] = parser
}

func (m *MeterDataCodec) ParseData(funcCode byte, data []byte) {
	if m.receiver == nil || data == nil || len(data) == 0 {
		return
	}
	switch funcCode {
	case RespondingNormallyNoNext, RespondingNormallyHasNext: //从站正常应答
		m.parseRespondingNormally(funcCode, data)
	case SlaveErrResponse: //从站异常应答
		m.receiver.MeterReadErrorResponse(funcCode, data[0])
	case NextRespondingNormallyNoNext, NextRespondingNormallyHasNext: //从站正常应答， 无后续帧,从站正常应答， 有后续帧
		m.parserReadNextResponse(funcCode, data)
	case NextSlaveErrResponse: //从站异常应答
		m.receiver.MeterReadErrorResponse(funcCode, data[0])
	case MeterSetResponse, MeterSetErrResponse: //主站设置，从站正常/异常应答
		m.receiver.MeterReqMasterSet(funcCode == MeterSetResponse, data[0])
	case MeterAddrResponse:
		m.parserMeterAddrResponse(data)
	case MeterSetMeterAddrResponse:
		m.parserSetAddrResponse(data)
	case FreezeCommandResponse, FreezeCommandErrorResponse:
		m.receiver.FreezeCommandResponse(funcCode == FreezeCommandResponse, data[0])
	default:
		m.receiver.ErrorData(funcCode, data, FuncCodeError)
	}
}

func (m *MeterDataCodec) parseRespondingNormally(funcCode byte, data []byte) {
	if len(data) < 4 {
		m.receiver.ErrorData(funcCode, data, DataDomainError)
		return
	}
	//是否存在后续帧
	hasNext := funcCode == RespondingNormallyHasNext
	//解析数据
	ident := data[:4]
	//判断是否存在数据解析器
	if parser, ok := m.parsers[hex.EncodeToString(ident)]; ok {
		parser.flush()
		if len(data) == 4 {
			m.receiver.MeterReadResponse(reverseBytes(ident), nil, hasNext, 0)
			return
		}
		err := parser.decode(data[4:])
		if err != nil {
			m.receiver.ErrorData(funcCode, data, err)
			return
		}
		m.receiver.MeterReadResponse(reverseBytes(ident), parser, hasNext, 0)
	} else {
		m.receiver.MeterDefaultReadResponse(funcCode, data)
	}
}

func (m *MeterDataCodec) parserReadNextResponse(funcCode byte, data []byte) {
	if len(data) < 5 {
		m.receiver.ErrorData(funcCode, data, DataDomainError)
		return
	}
	hasNext := funcCode == NextRespondingNormallyHasNext
	ident := data[:4]
	//判断是否存在数据解析器
	if parser, ok := m.parsers[hex.EncodeToString(ident)]; ok {
		err := parser.decode(data[4 : len(data)-1])
		if err != nil {
			m.receiver.ErrorData(funcCode, data, err)
			return
		}
		m.receiver.MeterReadResponse(reverseBytes(ident), parser, hasNext, data[len(data)-1])
	} else {
		m.receiver.MeterDefaultReadResponse(funcCode, data)
	}
}

func (m *MeterDataCodec) parserMeterAddrResponse(data []byte) {
	if len(data) < 6 {
		m.receiver.ErrorData(MeterAddrResponse, data, DataDomainError)
		return
	}
	m.receiver.MeterAddress(hex.EncodeToString(reverseBytes(data)))
}

func (m *MeterDataCodec) parserSetAddrResponse(data []byte) {
	if len(data) < 6 {
		m.receiver.ErrorData(MeterSetMeterAddrResponse, data, DataDomainError)
		return
	}
	m.receiver.MeterAddress(hex.EncodeToString(reverseBytes(data)))
}
