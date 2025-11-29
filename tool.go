package go_dlt645_2007

import (
	"encoding/hex"
	"errors"
	"strings"
	"time"
)

type MeterData[T ScalarOrVector] struct {
	Value  T
	Length byte
}

type ScalarOrVector interface {
	~int64 | ~uint64 | ~string |
		~[]int64 | ~[]uint64 | ~[]string | ~[]byte
}

// BuildMasterReadRequest 创建一个读数据/主站请求帧
// prefix 通配唤醒前缀
// meterId 表地址
// ident 数据标识
// block 负荷记录块数
// ts 给定时间
func BuildMasterReadRequest(prefix, meterId string, ident []byte, block byte, ts *time.Time) ([]byte, error) {
	if ident == nil || len(ident) != 4 {
		return nil, errors.New("data ident length error")
	}
	data := reverseBytes(ident)
	if block > 0x00 {
		data = append(data, block)
	}
	if ts != nil {
		mm := ts.Minute()     // 分钟 (0-59)
		hh := ts.Hour()       // 小时 (0-23)
		DD := ts.Day()        // 日 (1-31)
		MM := ts.Month()      // 月 (January-December)
		YY := ts.Year() % 100 // 年 (四位数)
		data = append(data, byte(mm), byte(hh), byte(DD), byte(MM), byte(YY))
	}
	statute := &MeterDlt645Protocol{prefix: prefix, Address: meterId, Data: data, ControlChar: MainStationRequestFrame}
	return statute.Encode()
}

// BuildMasterReadResponse 创建一个读数据/主站请求帧的正常应答
// prefix 通配唤醒前缀
// meterId 表地址
// ident 数据标识
// Value 值
// hasNext 是否存在后续帧，true-存在， false-不存在
func BuildMasterReadResponse[T ScalarOrVector](prefix, meterId string, ident []byte, value *MeterData[T], hasNext bool) ([]byte, error) {
	if ident == nil || len(ident) != 4 {
		return nil, errors.New("data ident length error")
	}
	if value == nil {
		statute := &MeterDlt645Protocol{prefix: prefix, Address: meterId, Data: ident, ControlChar: RespondingNormallyNoNext}
		return statute.Encode()
	}
	controlCode := RespondingNormallyNoNext
	if hasNext {
		controlCode = RespondingNormallyHasNext
	}
	var data = reverseBytes(ident)
	valArr, err := toLittleEndianBytes(value)
	if err != nil {
		return nil, err
	}
	data = append(data, valArr...)
	statute := &MeterDlt645Protocol{prefix: prefix, Address: meterId, Data: data, ControlChar: controlCode}
	return statute.Encode()
}

// BuildMeterAbnormalResponse 创建一个从站异常应答
// prefix 通配唤醒前缀
// meterId 表地址
// errCode 错误码
func BuildMeterAbnormalResponse(prefix, meterId string, errCode byte) ([]byte, error) {
	statute := &MeterDlt645Protocol{prefix: prefix, Address: meterId, Data: []byte{errCode}, ControlChar: SlaveErrResponse}
	return statute.Encode()
}

// BuildMasterReadNextDataRequest 创建一个主站读后续数据的请求帧
// prefix 通配唤醒前缀
// meterId 表地址
// ident 数据标识
// seq 帧序号 1～255。
func BuildMasterReadNextDataRequest(prefix, meterId string, ident []byte, seq byte) ([]byte, error) {
	if ident == nil || len(ident) != 4 {
		return nil, errors.New("data ident length error")
	}
	data := append(reverseBytes(ident), seq)
	statute := &MeterDlt645Protocol{prefix: prefix, Address: meterId, Data: data, ControlChar: ReadNextFrame}
	return statute.Encode()
}

// BuildMeterReadNextDataResponse 从站正常回复后续帧的应答
// prefix 通配唤醒前缀
// meterId 表地址
// ident 数据标识
// Value 数据
// seq 帧序号
// hasNext 是否存在后续帧
func BuildMeterReadNextDataResponse[T ScalarOrVector](prefix, meterId string, ident []byte, value *MeterData[T], seq byte, hasNext bool) ([]byte, error) {
	if ident == nil || len(ident) != 4 {
		return nil, errors.New("data ident length error")
	}
	if value == nil {
		statute := &MeterDlt645Protocol{prefix: prefix, Address: meterId, Data: ident, ControlChar: NextRespondingNormallyNoNext}
		return statute.Encode()
	}
	conctrlCode := NextRespondingNormallyNoNext
	if hasNext {
		conctrlCode = NextRespondingNormallyHasNext
	}
	data := reverseBytes(ident)
	valArr, err := toLittleEndianBytes(value)
	if err != nil {
		return nil, err
	}
	data = append(data, valArr...)
	data = append(data, seq)
	statute := &MeterDlt645Protocol{prefix: prefix, Address: meterId, Data: data, ControlChar: conctrlCode}
	return statute.Encode()
}

// BuildMeterReadNextErrResponse 从站异常回复后续帧的应答
// prefix 通配唤醒前缀
// meterId 表地址
// errCode 错误码
func BuildMeterReadNextErrResponse(prefix, meterId string, errCode byte) ([]byte, error) {
	statute := &MeterDlt645Protocol{prefix: prefix, Address: meterId, Data: []byte{errCode}, ControlChar: NextSlaveErrResponse}
	return statute.Encode()
}

/*---------------写数据------------------*/

// BuildMasterSetRequest 构建一个主站向从站请求设置数据(或编程)的报文
// prefix 通配唤醒前缀
// meterId 表地址
// ident 数据标识
// pwd 密码
// operatorCode 操作者代码
// Value 设定值
func BuildMasterSetRequest[T ScalarOrVector](prefix, meterId string, ident, pwd, operatorCode []byte, value *MeterData[T]) ([]byte, error) {
	if ident == nil || len(ident) != 4 {
		return nil, errors.New("data ident length error")
	}
	if pwd == nil || len(pwd) != 4 {
		return nil, errors.New("pwd length error")
	}
	if operatorCode == nil || len(operatorCode) != 4 {
		return nil, errors.New("operatorCode length error")
	}
	data := append(ident, pwd...)
	data = append(data, operatorCode...)
	valArr, err := toLittleEndianBytes(value)
	if err != nil {
		return nil, err
	}
	data = append(data, byte(len(valArr)))
	data = append(data, valArr...)
	statute := &MeterDlt645Protocol{prefix: prefix, Address: meterId, Data: data, ControlChar: MasterSetRequest}
	return statute.Encode()
}

// BuildMeterSetResponse 构建一个回复主站向从站请求设置数据(或编程)的正常应答报文
// prefix 通配唤醒前缀
// meterId 表地址
func BuildMeterSetResponse(prefix, meterId string) ([]byte, error) {
	statute := &MeterDlt645Protocol{prefix: prefix, Address: meterId, Data: []byte{0x00}, ControlChar: MeterSetResponse}
	return statute.Encode()
}

// BuildMeterSetErrResponse 构建一个回复主站向从站请求设置数据(或编程)的异常应答报文
// prefix 通配唤醒前缀
// meterId 表地址
// errorCode 错误码
func BuildMeterSetErrResponse(prefix, meterId string, errorCode byte) ([]byte, error) {
	statute := &MeterDlt645Protocol{prefix: prefix, Address: meterId, Data: []byte{errorCode}, ControlChar: MeterSetErrResponse}
	return statute.Encode()
}

// BuildMasterReadMeterAddrRequest 构建一个读取电表通讯地址的报文
// prefix 通配唤醒前缀
func BuildMasterReadMeterAddrRequest(prefix string) ([]byte, error) {
	statute := &MeterDlt645Protocol{prefix: prefix, Address: substitute + substitute + substitute + substitute + substitute + substitute, ControlChar: MasterReadMeterAddrRequest}
	return statute.Encode()
}

// BuildMasterReadMeterAddrResponse 构建一个回复读取电表通讯地址的报文
// prefix 通配唤醒前缀
// address 电表地址
func BuildMasterReadMeterAddrResponse(prefix string, address string) ([]byte, error) {
	addr, err := hex.DecodeString(address)
	if err != nil {
		return nil, errors.New("invalid address")
	}
	statute := &MeterDlt645Protocol{prefix: prefix, Address: address, ControlChar: MeterAddrResponse, Data: reverseBytes(addr)}
	return statute.Encode()
}

// BuildMasterSetMeterAddrRequest 构建一个设置电表通讯地址的报文
// prefix 通配唤醒前缀
// Address 电表地址
func BuildMasterSetMeterAddrRequest(prefix string, address string) ([]byte, error) {
	addr, err := hex.DecodeString(address)
	if err != nil {
		return nil, errors.New("invalid address")
	}
	statute := &MeterDlt645Protocol{prefix: prefix, Address: substitute + substitute + substitute + substitute + substitute + substitute, ControlChar: MasterSetMeterAddrRequest, Data: reverseBytes(addr)}
	return statute.Encode()
}

// BuildMeterSetMeterAddrResponse 构建一个回复设置电表通讯地址的报文
// prefix 通配唤醒前缀
// Address 电表地址
func BuildMeterSetMeterAddrResponse(prefix string, address string) ([]byte, error) {
	addr, err := hex.DecodeString(address)
	if err != nil {
		return nil, errors.New("invalid address")
	}
	statute := &MeterDlt645Protocol{prefix: prefix, Address: substitute + substitute + substitute + substitute + substitute + substitute, ControlChar: MeterSetMeterAddrResponse, Data: reverseBytes(addr)}
	return statute.Encode()
}

// BuildBroadcastTimeCalibration 广播校时
// prefix 通配唤醒前缀
// ti 需要设置的时间
func BuildBroadcastTimeCalibration(prefix string, ti time.Time) ([]byte, error) {
	ss := ti.Second()     //秒
	mm := ti.Minute()     // 分钟 (0-59)
	hh := ti.Hour()       // 小时 (0-23)
	DD := ti.Day()        // 日 (1-31)
	MM := ti.Month()      // 月 (January-December)
	YY := ti.Year() % 100 // 年 (四位数)
	data := []byte{byte(ss), byte(mm), byte(hh), byte(DD), byte(MM), byte(YY)}
	statute := &MeterDlt645Protocol{prefix: prefix, Address: BroadcastAddress, ControlChar: BroadcastTimeCalibration, Data: data}
	return statute.Encode()
}

// BuildFreezeCommandRequest 冻结命令
// prefix 通配唤醒前缀
// Address 电表地址
// ti 冻结时间
func BuildFreezeCommandRequest(prefix, address string, ti time.Time) ([]byte, error) {
	if strings.TrimSpace(address) == "" {
		address = BroadcastAddress
	}
	mm := ti.Minute() // 分钟 (0-59)
	hh := ti.Hour()   // 小时 (0-23)
	DD := ti.Day()    // 日 (1-31)
	MM := ti.Month()  // 月 (January-December)
	data := []byte{byte(mm), byte(hh), byte(DD), byte(MM)}
	statute := &MeterDlt645Protocol{prefix: prefix, Address: BroadcastAddress, ControlChar: FreezeCommand, Data: data}
	return statute.Encode()
}

// BuildFreezeCommandResponse 生成一个冻结命令的正确回复报文
// prefix 通配唤醒前缀
// Address 电表地址
func BuildFreezeCommandResponse(prefix, address string) ([]byte, error) {
	statute := &MeterDlt645Protocol{prefix: prefix, Address: BroadcastAddress, ControlChar: FreezeCommandResponse}
	return statute.Encode()
}

// BuildFreezeCommandErrorResponse 生成一个冻结命令的异常回复报文
// prefix 通配唤醒前缀
// Address 电表地址
func BuildFreezeCommandErrorResponse(prefix, address string) ([]byte, error) {
	statute := &MeterDlt645Protocol{prefix: prefix, Address: address, ControlChar: FreezeCommandErrorResponse}
	return statute.Encode()
}

func reverseBytes(original []byte) []byte {
	length := len(original)
	reversed := make([]byte, length)
	for i := 0; i < length; i++ {
		reversed[i] = original[length-1-i]
	}
	return reversed
}
