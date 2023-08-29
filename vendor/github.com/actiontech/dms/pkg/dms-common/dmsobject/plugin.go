package dmsobject

import (
	"context"
	"reflect"

	"github.com/iancoleman/strcase"
)

var OperateHandlers map[string]OperationHanler = make(map[string]OperationHanler)

// NOTE:
// The implemented structure must be named[CamelCase] by the combination of DataResourceType, OperationType, and OperationTimingType
type OperationHanler interface {
	Hanle(ctx context.Context, currentUserId string, dataResourceId string) error
}

func InitOperateHandlers(operationHandlers []OperationHanler) {
	for _, v := range operationHandlers {
		structName := strcase.ToSnake(reflect.TypeOf(v).Name())
		OperateHandlers[structName] = v
	}
}
