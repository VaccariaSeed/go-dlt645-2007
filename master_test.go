package go_dlt645_2007

import (
	"encoding/hex"
	"fmt"
	"strings"
	"testing"
)

var _ MeterDataReceiver = (*TestMeterParper)(nil)

type TestMeterParper struct{}

func (t *TestMeterParper) MeterDefaultReadResponse(funcCode byte, data []byte) {
	//TODO implement me
	panic("implement me")
}

func (t *TestMeterParper) MeterReadResponse(ident []byte, parser *MeterDataParser, hasNext bool, seq byte) {
	fmt.Println("MeterReadResponse", hex.EncodeToString(ident), parser.ObtainValue(), hasNext, seq)

}

func (t *TestMeterParper) MeterReadErrorResponse(funcCode byte, errCode byte) {
	//TODO implement me
	panic("implement me")
}

func (t *TestMeterParper) MeterReqMasterSet(isSuccess bool, errCode byte) {
	//TODO implement me
	panic("implement me")
}

func (t *TestMeterParper) MeterAddress(addr string) {
	//TODO implement me
	panic("implement me")
}

func (t *TestMeterParper) FreezeCommandResponse(isSuccess bool, errCode byte) {
	//TODO implement me
	panic("implement me")
}

func (t *TestMeterParper) ErrorData(funcCode byte, data []byte, err error) {
	//TODO implement me
	panic("implement me")
}

func TestMaster(t *testing.T) {
	//meter := NewMeter("", "00013310")
	//frame, err := meter.BuildMasterReadRequest([]byte{0x02, 0x01, 0x01, 0x00}, 0, nil)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println("读取一个数据")
	//fmt.Println(hex.EncodeToString(frame))
	//frame, err = meter.BuildMasterReadResponse([]byte{0x02, 0x01, 0x01, 0x00}, uint64(2219), 2, false)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println("返回A相电压")
	//fmt.Println(hex.EncodeToString(frame))
	codec := NewMeterDataCodec(&TestMeterParper{})
	dataParser, err := NewMeterDataParser(2, nil, 0.1, 0, "A")
	if err != nil {
		panic(err)
	}
	codec.Register([]byte{0x02, 0x01, 0x01, 0x00}, dataParser)
	dataParser, err = NewMeterDataParser(2, nil, 0.01, 0, "HZ")
	if err != nil {
		panic(err)
	}
	codec.Register([]byte{0x02, 0x80, 0x00, 0x02}, dataParser)
	pro := &MeterDlt645Protocol{}
	frame, _ := hex.DecodeString(strings.ReplaceAll("FE FE FE FE6800 51 44 18 11 1768910635 33 B3 35 36 834516", " ", ""))
	err = pro.Decode(frame)
	if err != nil {
		panic(err)
	}
	codec.ParseData(pro.ControlChar, pro.Data)
}
