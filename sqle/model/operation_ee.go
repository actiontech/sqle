//go:build enterprise
// +build enterprise

package model

func getConfigurableOperationCodeListForEE() []uint {
	return []uint{}
}

func additionalOperationForEE(opCode uint) string {
	return "未知动作"
}
