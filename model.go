package go_dlt645_2007

import "time"

type MasterReadRequestModel struct {
	ident    []byte    //数据标识
	block    byte      //负荷记录块数
	hasBlock bool      // 是否存在负荷记录块数
	ts       time.Time //给定时间
	hasTs    bool      //是否存在给定时间
}

// ObtainIdent 获取数据标识
func (m *MasterReadRequestModel) ObtainIdent() []byte {
	return m.ident
}

// ObtainBlock 获取给定负荷记录块数
func (m *MasterReadRequestModel) ObtainBlock() byte {
	return m.block
}

// ObtainTs 获取给定时间
func (m *MasterReadRequestModel) ObtainTs() time.Time {
	return m.ts
}

// HasBlock 是否存在负荷记录块数
func (m *MasterReadRequestModel) HasBlock() bool {
	return m.hasBlock
}

// HasTs 是否存在给定时间
func (m *MasterReadRequestModel) HasTs() bool {
	return m.hasTs
}

func (m *MasterReadRequestModel) decode(data []byte) error {
	if data == nil || len(data) < 4 {
		return DataDomainError
	}
	m.ident = data[:4]
	if len(data) == 5 {
		m.block = data[4]
		m.hasBlock = true
	}
	if len(data) >= 10 {
		m.hasTs = true
		mm, hh, DD, MM, YY := data[5], data[6], data[7], data[8], data[9]
		m.ts = time.Date(2000+int(YY), time.Month(MM), int(DD), int(hh), int(mm), 0, 0, time.Local)
	}
	return nil
}
