//go:build release
// +build release

package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLicense_AfterFind(t *testing.T) {
	content := "This license is for: &{WorkDurationDay:60 Version:演示环境 UserCount:10 NumberOfInstanceOfEachType:map[custom:{DBType:custom Count:3} mysql:{DBType:mysql Count:3}]};;1_XBm2N8t7coUEuhg7J5V8o9AYlhUfq2AmndctDHCxz9u~GyOKyJW0e~sVDuQVbkaKzAZQvpsGBqB~liD7svsTvbzD3ZHfdvEtSPkoYSnk2nxrYJLrW0wmzTVIicDWg1Dp2MICEK9T09Od3Xn1u4XWO7e182mzrHqncLOGKXJKlSrCsL_kWY6o6w8pWKL1Xdzduyq4uLdXuL9E6oOzyUMF3rYlnOhvoOwdoE;;9S~ViK_ZoRx8045cLM5pTZXCCpDEY_yxjfaLYGBMMOKyWpgc"
	l := &License{
		WorkDurationHour: 1,
		Content:          content,
	}
	assert.NoError(t, l.BeforeSave())
	l.Content = ""
	assert.NoError(t, l.AfterFind())
	assert.Equal(t, content, l.Content)
}
