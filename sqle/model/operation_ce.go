//go:build !enterprise
// +build !enterprise

package model

import (
	"context"

	"github.com/actiontech/sqle/sqle/locale"
)

func getConfigurableOperationCodeListForEE() []uint {
	return []uint{}
}

func additionalOperationForEE(ctx context.Context, opCode uint) string {
	return locale.Bundle.LocalizeMsgByCtx(ctx, locale.OpUnknown)
}
