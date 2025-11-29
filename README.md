# golang dlt645-2007规约

## 创建一个电表对象
```go
	meter := NewMeter("", "00013310")
```

#### 1.创建一个读数据的报文
```go
frame, err := meter.BuildMasterReadRequest([]byte{0x02, 0x01, 0x01, 0x00}, 0, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println("读取一个数据")
	fmt.Println(hex.EncodeToString(frame))
```

## 作为主站
#### 1.创建结果接收器
```go
var _ MeterDataReceiver = (*TestMeterParper)(nil)
```
#### 2.构建解码器
```go
codec := NewMeterDataCodec(&TestMeterParper{})
```

#### 3.构建解析规则并注册解析规则（以A相电压为例）
```go
//构建
dataParser, err := NewMeterDataParser(2, nil, 0.1, 0, "A")
if err != nil {
panic(err)
}
//注册
codec.Register([]byte{0x02, 0x01, 0x01, 0x00}, dataParser)
```

#### 解析报文
```go
pro := &MeterDlt645Protocol{}
	frame, _ := hex.DecodeString(strings.ReplaceAll("FE FE FE FE6800 51 44 18 11 1768910635 33 B3 35 36 834516", " ", ""))
	err = pro.Decode(frame)
	if err != nil {
		panic(err)
	}
```

#### 解析结果
```go
codec.ParseData(pro.ControlChar, pro.Data)
```

## 作为电表
#### 1.实现解析器
```go
var _ MasterDataReceiver = (*MasterAnalyzer)(nil)
```

#### 2. 解析报文
```
pro := &MeterDlt645Protocol{}
frame, _ := hex.DecodeString(strings.ReplaceAll("FE FE FE FE6800 51 44 18 11 1768910635 33 B3 35 36 834516", " ", ""))
err = pro.Decode(frame)
if err != nil {
panic(err)
}
```

#### 3.构建解析器
```go
codec := NewMasterDataCodec(MasterAnalyzer)
```

#### 4. 解析
```go
codec.ParseData(pro.ControlChar, pro.Data)
```