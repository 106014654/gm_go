package config

import (
	"encoding/json"
	"fmt"
	"google.golang.org/protobuf/proto"
	"reflect"
	"strconv"
	"sync/atomic"
)

type Value interface {
	Bool() (bool, error)
	Int() (int64, error)
	Float() (float64, error)
	String() (string, error)
	Scan(interface{}) error
	Load() interface{}
	Store(interface{})
}

type atomicValue struct {
	atomic.Value
}

func (v *atomicValue) typeAssertError() error {
	return fmt.Errorf("type assert to %v failed", reflect.TypeOf(v.Load()))
}

func (v *atomicValue) Scan(obj interface{}) error {
	data, err := json.Marshal(v.Load())
	if err != nil {
		return err
	}
	if pb, ok := obj.(proto.Message); ok {
		return json.Unmarshal(data, pb)
	}
	return json.Unmarshal(data, obj)
}

type errValue struct {
	err error
}

func (v *atomicValue) Bool() (bool, error) {
	switch val := v.Load().(type) {
	case bool:
		return val, nil
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, string:
		return strconv.ParseBool(fmt.Sprint(val))
	}
	return false, v.typeAssertError()
}

func (v *atomicValue) Int() (int64, error) {
	switch val := v.Load().(type) {
	case int:
		return int64(val), nil
	case int8:
		return int64(val), nil
	case int16:
		return int64(val), nil
	case int32:
		return int64(val), nil
	case int64:
		return val, nil
	case uint:
		return int64(val), nil
	case uint8:
		return int64(val), nil
	case uint16:
		return int64(val), nil
	case uint32:
		return int64(val), nil
	case uint64:
		return int64(val), nil
	case float32:
		return int64(val), nil
	case float64:
		return int64(val), nil
	case string:
		return strconv.ParseInt(val, 10, 64)
	}
	return 0, v.typeAssertError()
}

func (v *atomicValue) Float() (float64, error) {
	switch val := v.Load().(type) {
	case int:
		return float64(val), nil
	case int8:
		return float64(val), nil
	case int16:
		return float64(val), nil
	case int32:
		return float64(val), nil
	case int64:
		return float64(val), nil
	case uint:
		return float64(val), nil
	case uint8:
		return float64(val), nil
	case uint16:
		return float64(val), nil
	case uint32:
		return float64(val), nil
	case uint64:
		return float64(val), nil
	case float32:
		return float64(val), nil
	case float64:
		return val, nil
	case string:
		return strconv.ParseFloat(val, 64)
	}
	return 0.0, v.typeAssertError()
}

func (v *atomicValue) String() (string, error) {
	switch val := v.Load().(type) {
	case string:
		return val, nil
	case bool, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return fmt.Sprint(val), nil
	case []byte:
		return string(val), nil
	case fmt.Stringer:
		return val.String(), nil
	}
	return "", v.typeAssertError()
}

func (v errValue) Bool() (bool, error)     { return false, v.err }
func (v errValue) Int() (int64, error)     { return 0, v.err }
func (v errValue) Float() (float64, error) { return 0.0, v.err }
func (v errValue) String() (string, error) { return "", v.err }
func (v errValue) Scan(interface{}) error  { return v.err }
func (v errValue) Load() interface{}       { return nil }
func (v errValue) Store(interface{})       {}
