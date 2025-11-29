package go_dlt645_2007

import "encoding/hex"

type MasterDataReceiver interface {
	MasterReadRequest(req *MasterReadRequestModel)                           // 主站读取数据
	MasterReadNextRequest(ident []byte, seq byte)                            // 主站读取后续数据
	MasterSetRequest(ident []byte, pwd []byte, operator []byte, data []byte) //主站向从站请求设置数据
	MasterReadMeterAddrRequest()                                             //主站请求读地址
	MasterSetMeterAddrRequest(addr string)                                   //主站设置地址
	BroadcastTimeCalibration(ss, mm, hh, DD, MM, YY byte)                    //广播校时
	FreezeCommand(mm, hh, DD, MM byte)                                       //冻结命令
	ErrorData(funcCode byte, data []byte, err error)                         //解析失败的数据会调用这个方法
}

func NewMasterDataCodec(receiver MasterDataReceiver) *MasterDataCodec {
	return &MasterDataCodec{receiver: receiver}
}

type MasterDataCodec struct {
	receiver MasterDataReceiver
}

func (m *MasterDataCodec) ParseData(funcCode byte, data []byte) {
	if m.receiver == nil || data == nil || len(data) == 0 {
		return
	}
	switch funcCode {
	case MainStationRequestFrame: //读数据,主站请求帧
		m.parseMainStationRequestFrame(data)
	case ReadNextFrame:
		m.parseReadNextFrame(data)
	case MasterSetRequest: //主站向从站请求设置数据(或编程)
		m.parseMasterSetRequest(data)
	case MasterReadMeterAddrRequest: //主站读电表地址
		m.receiver.MasterReadMeterAddrRequest()
	case MasterSetMeterAddrRequest: //设置某从站的通信地址，仅支持点对点通信
		m.parseMasterSetMeterAddrRequest(data)
	case BroadcastTimeCalibration: //广播校时
		m.parseBroadcastTimeCalibration(data)
	case FreezeCommand: //冻结命令
		m.parseFreezeCommand(data)
	default:
		m.receiver.ErrorData(funcCode, data, FuncCodeError)
	}
}

// 解析读数据,主站请求帧
func (m *MasterDataCodec) parseMainStationRequestFrame(data []byte) {
	model := &MasterReadRequestModel{}
	err := model.decode(data)
	if err != nil {
		m.receiver.ErrorData(MainStationRequestFrame, data, err)
		return
	}
	m.receiver.MasterReadRequest(model)
}

// 解析请求读后续数据
func (m *MasterDataCodec) parseReadNextFrame(data []byte) {
	if len(data) < 5 {
		m.receiver.ErrorData(ReadNextFrame, data, DataDomainError)
		return
	}
	ident := data[:4]
	seq := data[4]
	m.receiver.MasterReadNextRequest(ident, seq)
}

func (m *MasterDataCodec) parseMasterSetRequest(data []byte) {
	if len(data) < 12 {
		m.receiver.ErrorData(MasterSetRequest, data, DataDomainError)
		return
	}
	ident := data[:4]
	pwd := data[4:8]
	operator := data[8:12]
	data = data[12:]
	m.receiver.MasterSetRequest(ident, pwd, operator, data)
}

func (m *MasterDataCodec) parseMasterSetMeterAddrRequest(data []byte) {
	if len(data) < 6 {
		m.receiver.ErrorData(MasterSetMeterAddrRequest, data, DataDomainError)
		return
	}
	m.receiver.MasterSetMeterAddrRequest(hex.EncodeToString(reverseBytes(data)))
}

func (m *MasterDataCodec) parseBroadcastTimeCalibration(data []byte) {
	if len(data) < 6 {
		m.receiver.ErrorData(BroadcastTimeCalibration, data, DataDomainError)
		return
	}
	ss, mm, hh, DD, MM, YY := data[0], data[1], data[2], data[3], data[4], data[5]
	m.receiver.BroadcastTimeCalibration(ss, mm, hh, DD, MM, YY)
}

func (m *MasterDataCodec) parseFreezeCommand(data []byte) {
	if len(data) < 6 {
		m.receiver.ErrorData(BroadcastTimeCalibration, data, DataDomainError)
		return
	}
	mm, hh, DD, MM := data[0], data[1], data[2], data[3]
	m.receiver.FreezeCommand(mm, hh, DD, MM)
}
