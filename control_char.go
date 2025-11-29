package go_dlt645_2007

const (
	MainStationRequestFrame       byte = 0x11 //读数据,主站请求帧
	RespondingNormallyNoNext      byte = 0x91 //从站正常应答， 无后续帧
	RespondingNormallyHasNext     byte = 0xB1 //从站正常应答， 有后续帧
	SlaveErrResponse              byte = 0xD1 //从站异常应答
	ReadNextFrame                 byte = 0x12 //主站读后续数据
	NextRespondingNormallyNoNext  byte = 0x92 //从站正常应答， 无后续帧
	NextRespondingNormallyHasNext byte = 0xB2 //从站正常应答， 有后续帧
	NextSlaveErrResponse          byte = 0xD2 //从站异常应答
	MasterSetRequest              byte = 0x14 //主站向从站请求设置数据(或编程)
	MeterSetResponse              byte = 0x94 //主站设置，从站正常应答
	MeterSetErrResponse           byte = 0xD4 //主站设置，从站异常应答
	MasterReadMeterAddrRequest    byte = 0x13 //主站读电表地址
	MeterAddrResponse             byte = 0x93
	MasterSetMeterAddrRequest     byte = 0x15 //设置某从站的通信地址，仅支持点对点通信
	MeterSetMeterAddrResponse     byte = 0x95
	BroadcastTimeCalibration      byte = 0x08 //广播校时
	FreezeCommand                 byte = 0x16 //冻结命令
	FreezeCommandResponse         byte = 0x96
	FreezeCommandErrorResponse    byte = 0xD6
)
