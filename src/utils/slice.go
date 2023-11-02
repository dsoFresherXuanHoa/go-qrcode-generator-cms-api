package utils

import (
	"errors"
	"fmt"
)

var (
	ErrNotSliceValueGreaterThan = errors.New("all value of slice less than target value, last index not found")
)

type sliceUtil struct{}

func NewSliceUtil() *sliceUtil {
	return &sliceUtil{}
}

func (sliceUtil) LastIndexLessThanSliceValue(slice []int, value int) (*int, error) {
	for i, s := range slice {
		if s > value {
			return &i, nil
		}
	}
	fmt.Println("Error while get last index of slice greater than target value: all value of slice less than target value")
	return nil, ErrNotSliceValueGreaterThan
}
