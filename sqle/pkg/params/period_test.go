package params

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPeriods_ScanValue(t *testing.T) {
	ps := Periods{
		&Period{
			StartHour:   1,
			StartMinute: 2,
			EndHour:     3,
			EndMinute:   4,
		},
		&Period{
			StartHour:   1,
			StartMinute: 3,
			EndHour:     2,
			EndMinute:   4,
		},
		&Period{
			StartHour:   2,
			StartMinute: 4,
			EndHour:     3,
			EndMinute:   1,
		},
	}

	data, err := ps.Value()
	assert.NoError(t, err)

	var ps2 Periods
	err = ps2.Scan(data)
	assert.NoError(t, err)
	assert.Equal(t, ps, ps2)

	ps3 := Periods{
		&Period{
			StartHour:   0,
			StartMinute: 2,
			EndHour:     3,
			EndMinute:   4,
		},
		&Period{
			StartHour:   0,
			StartMinute: 4,
			EndHour:     2,
			EndMinute:   4,
		},
		&Period{
			StartHour:   0,
			StartMinute: 4,
			EndHour:     3,
			EndMinute:   1,
		},
	}
	assert.NotEqual(t, ps3, ps2)

	var emptyPs Periods
	data, err = emptyPs.Value()
	assert.NoError(t, err)

	var emptyPs2 Periods
	err = emptyPs2.Scan(data)
	assert.NoError(t, err)
	assert.Equal(t, emptyPs, emptyPs2)

	data = []byte("this is test scan fail")
	var failPs Periods
	err = failPs.Scan(data)
	assert.Error(t, err)
}
