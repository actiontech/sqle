package model

import (
	"testing"
	"time"

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

func TestPeriods_SelfCheck(t *testing.T) {
	ps1 := Periods{
		{
			StartHour:   0,
			StartMinute: 0,
			EndHour:     23,
			EndMinute:   59,
		}, {
			StartHour:   10,
			StartMinute: 20,
			EndHour:     10,
			EndMinute:   40,
		},
	}
	assert.Equal(t, ps1.SelfCheck(), true)
	ps2 := Periods{
		{
			StartHour:   1,
			StartMinute: 2,
			EndHour:     0,
			EndMinute:   4,
		}, {
			StartHour:   10,
			StartMinute: 20,
			EndHour:     30,
			EndMinute:   40,
		},
	}
	assert.Equal(t, ps2.SelfCheck(), false)
	ps3 := Periods{
		{
			StartHour:   1,
			StartMinute: 2,
			EndHour:     3,
			EndMinute:   0,
		}, {
			StartHour:   10,
			StartMinute: 20,
			EndHour:     30,
			EndMinute:   40,
		},
	}
	assert.Equal(t, ps3.SelfCheck(), false)
	ps4 := Periods{
		{
			StartHour:   25,
			StartMinute: 2,
			EndHour:     3,
			EndMinute:   4,
		}, {
			StartHour:   10,
			StartMinute: 20,
			EndHour:     30,
			EndMinute:   40,
		},
	}
	assert.Equal(t, ps4.SelfCheck(), false)
	ps5 := Periods{
		{
			StartHour:   1,
			StartMinute: 60,
			EndHour:     3,
			EndMinute:   4,
		}, {
			StartHour:   10,
			StartMinute: 20,
			EndHour:     30,
			EndMinute:   40,
		},
	}
	assert.Equal(t, ps5.SelfCheck(), false)
	ps6 := Periods{
		{
			StartHour:   -1,
			StartMinute: 2,
			EndHour:     3,
			EndMinute:   4,
		}, {
			StartHour:   10,
			StartMinute: 20,
			EndHour:     30,
			EndMinute:   40,
		},
	}
	assert.Equal(t, ps6.SelfCheck(), false)
	ps7 := Periods{
		{
			StartHour:   1,
			StartMinute: 2,
			EndHour:     3,
			EndMinute:   -4,
		}, {
			StartHour:   10,
			StartMinute: 20,
			EndHour:     30,
			EndMinute:   40,
		},
	}
	assert.Equal(t, ps7.SelfCheck(), false)

}

func TestPeriods_IsWithinScope(t *testing.T) {
	ps := Periods{
		{
			StartHour:   2,
			StartMinute: 1,
			EndHour:     3,
			EndMinute:   2,
		}, {
			StartHour:   4,
			StartMinute: 3,
			EndHour:     5,
			EndMinute:   4,
		},
	}

	t0, err := time.Parse("2006-01-02 15:04:05", "2017-12-08 05:04:03")
	assert.NoError(t, err)
	assert.Equal(t, ps.IsWithinScope(t0), true)

	t1, err := time.Parse("2006-01-02 15:04:05", "2017-12-08 02:01:03")
	assert.NoError(t, err)
	assert.Equal(t, ps.IsWithinScope(t1), true)

	t2, err := time.Parse("2006-01-02 15:04:05", "2017-12-08 03:01:53")
	assert.NoError(t, err)
	assert.Equal(t, ps.IsWithinScope(t2), true)

	t3, err := time.Parse("2006-01-02 15:04:05", "2017-12-08 01:01:53")
	assert.NoError(t, err)
	assert.Equal(t, ps.IsWithinScope(t3), false)

	t4, err := time.Parse("2006-01-02 15:04:05", "2017-12-08 03:03:53")
	assert.NoError(t, err)
	assert.Equal(t, ps.IsWithinScope(t4), false)

	t5, err := time.Parse("2006-01-02 15:04:05", "2017-12-08 23:01:53")
	assert.NoError(t, err)
	assert.Equal(t, ps.IsWithinScope(t5), false)

}
