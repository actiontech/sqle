package service

// import (
// 	"fmt"
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// )

// func TestCBVersion_LessThan(t *testing.T) {
// 	// 判断V1是否比V2小
// 	versions := []struct {
// 		V1     CBVersion
// 		V2     CBVersion
// 		result bool
// 	}{
// 		// 应当只比较前三位
// 		{
// 			V1: CBVersion{
// 				[]int{2, 2, 2, 2},
// 			},
// 			V2: CBVersion{
// 				[]int{2, 2, 2, 3},
// 			},
// 			result: false,
// 		},
// 		// 前两位相同
// 		{
// 			V1: CBVersion{
// 				[]int{2, 2, 2},
// 			},
// 			V2: CBVersion{
// 				[]int{2, 2, 2},
// 			},
// 			result: false,
// 		}, {
// 			V1: CBVersion{
// 				[]int{2, 2, 3},
// 			},
// 			V2: CBVersion{
// 				[]int{2, 2, 2},
// 			},
// 			result: false,
// 		}, {
// 			V1: CBVersion{
// 				[]int{2, 2, 1},
// 			},
// 			V2: CBVersion{
// 				[]int{2, 2, 2},
// 			},
// 			result: true,
// 		},
// 		// 第一位相同
// 		{
// 			V1: CBVersion{
// 				[]int{2, 2, 2},
// 			},
// 			V2: CBVersion{
// 				[]int{2, 1, 2},
// 			},
// 			result: false,
// 		}, {
// 			V1: CBVersion{
// 				[]int{2, 1, 2},
// 			},
// 			V2: CBVersion{
// 				[]int{2, 2, 2},
// 			},
// 			result: true,
// 		},
// 		// 第一位不同
// 		{
// 			V1: CBVersion{
// 				[]int{2, 2, 2},
// 			},
// 			V2: CBVersion{
// 				[]int{1, 2, 2},
// 			},
// 			result: false,
// 		}, {
// 			V1: CBVersion{
// 				[]int{1, 2, 2},
// 			},
// 			V2: CBVersion{
// 				[]int{2, 2, 2},
// 			},
// 			result: true,
// 		},
// 	}

// 	for _, version := range versions {
// 		assert.Equal(t, version.V1.LessThan(version.V2), version.result,
// 			fmt.Sprintf("V1: %v || V2: %v || result: %v", version.V1, version.V2, version.result))
// 	}
// }
