// Code generated by ent, DO NOT EDIT.

package ent

import (
	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

func convert(val interface{}, typ string) (interface{}, error) {
	switch typ {
	case "string", "json.RawMessage":
		return cast.ToStringE(val)
	case "[]byte":
		s, err := cast.ToStringE(val)
		if err != nil {
			return nil, err
		}
		return []byte(s), err
	case "[16]byte":
		var byteArr [16]byte
		s, err := cast.ToStringE(val)
		if err != nil {
			return nil, err
		}
		copy(byteArr[:], s)
		return byteArr, err
	case "time.Time":
		return cast.ToTimeE(val)
	case "bool":
		return cast.ToBoolE(val)
	case "float64":
		return cast.ToFloat64E(val)
	case "float32":
		return cast.ToFloat32E(val)
	case "int64":
		return cast.ToInt64E(val)
	case "int32":
		return cast.ToInt32E(val)
	case "int16":
		return cast.ToInt16E(val)
	case "int8":
		return cast.ToInt8E(val)
	case "int":
		return cast.ToIntE(val)
	case "uint64":
		return cast.ToUint64E(val)
	case "uint32":
		return cast.ToUint32E(val)
	case "uint16":
		return cast.ToUint16E(val)
	case "uint8":
		return cast.ToUint8E(val)
	case "uint":
		return cast.ToUintE(val)
	}
	return nil, errors.Errorf("type '%s' not support", typ)
}
