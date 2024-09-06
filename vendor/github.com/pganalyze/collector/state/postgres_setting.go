package state

import "github.com/guregu/null"

type PostgresSetting struct {
	Name         string      `json:"name"`
	CurrentValue null.String `json:"current_value"`
	Unit         null.String `json:"unit"`
	BootValue    null.String `json:"boot_value"`
	ResetValue   null.String `json:"reset_value"`
	Source       null.String `json:"source"`
	SourceFile   null.String `json:"sourcefile"`
	SourceLine   null.String `json:"sourceline"`
}
