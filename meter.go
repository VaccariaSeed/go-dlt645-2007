package go_dlt645_2007

import (
	"errors"
	"strings"
	"time"
)

// NewMeter 创建一个新的电表
// prefix 唤醒符
// address 表计地址
func NewMeter(prefix, address string) *Meter {
	return &Meter{prefix: strings.TrimSpace(prefix), address: strings.TrimSpace(address)}
}

// Meter 表计结构体
type Meter struct {
	prefix  string //唤醒符
	address string //表计地址
}

// BuildMasterReadRequest 创建一个读数据/主站请求帧
// ident 数据标识
// block 负荷记录块数
// ts 给定时间
func (m *Meter) BuildMasterReadRequest(ident []byte, block byte, ts *time.Time) ([]byte, error) {
	return BuildMasterReadRequest(m.prefix, m.address, ident, block, ts)
}

// BuildMasterReadResponse 创建一个读数据/主站请求帧的正常应答
// ident 数据标识
// value 值
// valueLength 数据长度
// hasNext 是否存在后续帧，true-存在， false-不存在
func (m *Meter) BuildMasterReadResponse(ident []byte, value interface{}, valueLength byte, hasNext bool) ([]byte, error) {
	switch v := value.(type) {
	case int64:
		return BuildMasterReadResponse[int64](m.prefix, m.address, ident, &MeterData[int64]{Value: v, Length: valueLength}, hasNext)
	case uint64:
		return BuildMasterReadResponse[uint64](m.prefix, m.address, ident, &MeterData[uint64]{Value: v, Length: valueLength}, hasNext)
	case string:
		return BuildMasterReadResponse[string](m.prefix, m.address, ident, &MeterData[string]{Value: v, Length: valueLength}, hasNext)
	case []byte:
		return BuildMasterReadResponse[[]byte](m.prefix, m.address, ident, &MeterData[[]byte]{Value: v, Length: valueLength}, hasNext)
	case []int64:
		return BuildMasterReadResponse[[]int64](m.prefix, m.address, ident, &MeterData[[]int64]{Value: v, Length: valueLength}, hasNext)
	case []uint64:
		return BuildMasterReadResponse[[]uint64](m.prefix, m.address, ident, &MeterData[[]uint64]{Value: v, Length: valueLength}, hasNext)
	case []string:
		return BuildMasterReadResponse[[]string](m.prefix, m.address, ident, &MeterData[[]string]{Value: v, Length: valueLength}, hasNext)
	default:
		return nil, errors.New("invalid type")
	}
}

// BuildMeterAbnormalResponse 创建一个读数据/主站请求帧的从站异常应答
// errCode 错误码
func (m *Meter) BuildMeterAbnormalResponse(errCode byte) ([]byte, error) {
	return BuildMeterAbnormalResponse(m.prefix, m.address, errCode)
}

// BuildMasterReadNextDataRequest 创建一个主站读后续数据的请求帧
// ident 数据标识
// seq 帧序号
func (m *Meter) BuildMasterReadNextDataRequest(ident []byte, seq byte) ([]byte, error) {
	return BuildMasterReadNextDataRequest(m.prefix, m.address, ident, seq)
}

// BuildMeterReadNextDataResponse 从站正常回复后续帧的应答
// ident 数据标识
// value 数据
// valueLength 数据长度
// seq 帧序号
// hasNext 是否存在后续帧
func (m *Meter) BuildMeterReadNextDataResponse(ident []byte, value interface{}, valueLength byte, seq byte, hasNext bool) ([]byte, error) {
	switch v := value.(type) {
	case int64:
		return BuildMeterReadNextDataResponse[int64](m.prefix, m.address, ident, &MeterData[int64]{Value: v, Length: valueLength}, seq, hasNext)
	case uint64:
		return BuildMeterReadNextDataResponse[uint64](m.prefix, m.address, ident, &MeterData[uint64]{Value: v, Length: valueLength}, seq, hasNext)
	case string:
		return BuildMeterReadNextDataResponse[string](m.prefix, m.address, ident, &MeterData[string]{Value: v, Length: valueLength}, seq, hasNext)
	case []byte:
		return BuildMeterReadNextDataResponse[[]byte](m.prefix, m.address, ident, &MeterData[[]byte]{Value: v, Length: valueLength}, seq, hasNext)
	case []int64:
		return BuildMeterReadNextDataResponse[[]int64](m.prefix, m.address, ident, &MeterData[[]int64]{Value: v, Length: valueLength}, seq, hasNext)
	case []uint64:
		return BuildMeterReadNextDataResponse[[]uint64](m.prefix, m.address, ident, &MeterData[[]uint64]{Value: v, Length: valueLength}, seq, hasNext)
	case []string:
		return BuildMeterReadNextDataResponse[[]string](m.prefix, m.address, ident, &MeterData[[]string]{Value: v, Length: valueLength}, seq, hasNext)
	default:
		return nil, errors.New("invalid type")
	}
}

// BuildMeterReadNextErrResponse 从站异常回复后续帧的应答
// errCode 错误码
func (m *Meter) BuildMeterReadNextErrResponse(errCode byte) ([]byte, error) {
	return BuildMeterReadNextErrResponse(m.prefix, m.address, errCode)
}

/*---------------写数据------------------*/

// BuildMasterSetRequest 构建一个主站向从站请求设置数据(或编程)的报文
// ident 数据标识
// pwd 密码
// operatorCode 操作者代码
// value 设定值
// valueLength 数据长度
func (m *Meter) BuildMasterSetRequest(ident, pwd, operatorCode []byte, value interface{}, valueLength byte) ([]byte, error) {
	switch v := value.(type) {
	case int64:
		return BuildMasterSetRequest[int64](m.prefix, m.address, ident, pwd, operatorCode, &MeterData[int64]{Value: v, Length: valueLength})
	case uint64:
		return BuildMasterSetRequest[uint64](m.prefix, m.address, ident, pwd, operatorCode, &MeterData[uint64]{Value: v, Length: valueLength})
	case string:
		return BuildMasterSetRequest[string](m.prefix, m.address, ident, pwd, operatorCode, &MeterData[string]{Value: v, Length: valueLength})
	case []byte:
		return BuildMasterSetRequest[[]byte](m.prefix, m.address, ident, pwd, operatorCode, &MeterData[[]byte]{Value: v, Length: valueLength})
	case []int64:
		return BuildMasterSetRequest[[]int64](m.prefix, m.address, ident, pwd, operatorCode, &MeterData[[]int64]{Value: v, Length: valueLength})
	case []uint64:
		return BuildMasterSetRequest[[]uint64](m.prefix, m.address, ident, pwd, operatorCode, &MeterData[[]uint64]{Value: v, Length: valueLength})
	case []string:
		return BuildMasterSetRequest[[]string](m.prefix, m.address, ident, pwd, operatorCode, &MeterData[[]string]{Value: v, Length: valueLength})
	default:
		return nil, errors.New("invalid type")
	}
}

// BuildMeterSetResponse 构建一个回复主站向从站请求设置数据(或编程)的正常应答报文
func (m *Meter) BuildMeterSetResponse() ([]byte, error) {
	return BuildMeterSetResponse(m.prefix, m.address)
}

// BuildMeterSetErrResponse 构建一个回复主站向从站请求设置数据(或编程)的异常应答报文
// errorCode 错误码
func (m *Meter) BuildMeterSetErrResponse(errorCode byte) ([]byte, error) {
	return BuildMeterSetErrResponse(m.prefix, m.address, errorCode)
}

// BuildMasterReadMeterAddrResponse 构建一个回复读取电表通讯地址的报文
func (m *Meter) BuildMasterReadMeterAddrResponse() ([]byte, error) {
	return BuildMasterReadMeterAddrResponse(m.prefix, m.address)
}

// BuildMasterSetMeterAddrRequest 构建一个设置电表通讯地址的报文
func (m *Meter) BuildMasterSetMeterAddrRequest() ([]byte, error) {
	return BuildMasterSetMeterAddrRequest(m.prefix, m.address)
}

// BuildMeterSetMeterAddrResponse 构建一个回复设置电表通讯地址的报文
func (m *Meter) BuildMeterSetMeterAddrResponse() ([]byte, error) {
	return BuildMeterSetMeterAddrResponse(m.prefix, m.address)
}
