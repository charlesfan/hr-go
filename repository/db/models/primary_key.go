package models

import (
	"github.com/charlesfan/hr-go/utils/strs"
)

type ValueType interface {
	string | int | int32 | int64 | float32 | float64
}

type PrimaryKey[T ValueType] map[string]T

func (pk PrimaryKey[ValueType]) Validate() bool {
	for k, v := range pk {
		if strs.IsBlank(k) {
			return false
		}

		switch t := any(v).(type) {
		case string:
			if strs.IsBlank(t) {
				return false
			}
		}
	}
	return true
}
