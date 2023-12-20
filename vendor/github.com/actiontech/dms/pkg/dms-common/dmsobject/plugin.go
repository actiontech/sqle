package dmsobject

import (
	"context"
	"reflect"

	"github.com/iancoleman/strcase"
)

var operateHandlers map[string]OperationHandler = make(map[string]OperationHandler)

// OperationHandler NOTE:
// The implemented structure must be named[CamelCase] by the combination of DataResourceType, OperationType, and OperationTimingType
type OperationHandler interface {
	Handle(ctx context.Context, currentUserId string, objId string) error
}

type DefaultOperateHandle struct {
}

func (f DefaultOperateHandle) Handle(ctx context.Context, currentUserId string, objId string) error {
	return nil
}

func InitOperateHandlers(operationHandlers []OperationHandler) {
	for _, v := range operationHandlers {
		structName := strcase.ToSnake(reflect.TypeOf(v).Name())
		operateHandlers[structName] = v
	}
}

func GetOperateHandle(name string) OperationHandler {
	handle, ok := operateHandlers[name]
	if ok {
		return handle
	}

	return DefaultOperateHandle{}
}
