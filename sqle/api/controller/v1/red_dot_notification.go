package v1

import (
	"github.com/labstack/echo/v4"
)

type RedDotModule interface {
	Name() string
	HasRedDot(ctx echo.Context) (bool, error)
}

type RedDot struct {
	ModuleName string
	HasRedDot  bool
}

var redDotList []RedDotModule

func RegisterRedDotModules(redDotModule ...RedDotModule) {
	redDotList = append(redDotList, redDotModule...)
}

func GetSystemModuleRedDotsList(ctx echo.Context) ([]*RedDot, error) {
	redDots := make([]*RedDot, len(redDotList))
	for i, rd := range redDotList {
		hasRedDot, err := rd.HasRedDot(ctx)
		if err != nil {
			return nil, err
		}
		redDots[i] = &RedDot{
			ModuleName: rd.Name(),
			HasRedDot:  hasRedDot,
		}
	}
	return redDots, nil
}
