package go_dlt645_2007

import (
	"encoding/hex"
	"errors"
	"fmt"
)

// toLittleEndianBytes 将 MeterData 类型转换为小端字节序的 []byte
func toLittleEndianBytes[T ScalarOrVector](value *MeterData[T]) ([]byte, error) {
	var result []byte
	var err error
	// 使用类型断言处理各种类型
	switch v := any(value.Value).(type) {
	case int64:
		result, err = int64ToBytes(v, value.Length)
	case uint64:
		result, err = uint64ToBytes(v, value.Length)
	case string:
		result, err = hex.DecodeString(v)
		if err == nil {
			result = reverseBytes(result)
		}
	case []int64:
		result, err = int64ArrayToBytes(v, value.Length)
	case []uint64:
		result, err = unt64ArrayToBytes(v, value.Length)
	case []string:
		result, err = strArrayToBytes(v)
	case []byte:
		result = v
	default:
		err = fmt.Errorf("toLittleEndianBytes: unknown type %T", value)
	}
	return result, err
}

func strArrayToBytes(v []string) ([]byte, error) {
	var result []byte
	for _, str := range v {
		b, err := hex.DecodeString(str)
		if err != nil {
			return nil, err
		}
		result = append(result, reverseBytes(b)...)
	}
	return result, nil
}

func unt64ArrayToBytes(v []uint64, length byte) ([]byte, error) {
	var result []byte
	for i := 0; i < len(v); i++ {
		d, err := uint64ToBytes(v[i], length)
		if err != nil {
			return nil, err
		}
		result = append(result, d...)
	}
	return result, nil
}

func int64ArrayToBytes(v []int64, length byte) ([]byte, error) {
	var result []byte
	for i := 0; i < len(v); i++ {
		d, err := int64ToBytes(v[i], length)
		if err != nil {
			return nil, err
		}
		result = append(result, d...)
	}
	return result, nil
}

func uint64ToBytes(value uint64, length byte) ([]byte, error) {
	if length == 0 {
		return nil, errors.New("length must be greater than 0")
	}
	per := fmt.Sprintf("%d", value)
	for len(per) < int(length*2) {
		per = "0" + per
	}
	result, err := hex.DecodeString(per)
	if err != nil {
		return nil, err
	}
	return reverseBytes(result), nil
}

func int64ToBytes(value int64, length byte) ([]byte, error) {
	return uint64ToBytes(uint64(value), length)
}
