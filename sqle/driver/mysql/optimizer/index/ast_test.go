package index

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_selectAST(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input      string
		orderBy    []string
		whereEqual []string
		projection []string
	}{
		{"select id, a from t where a = 1 and b = 1 order by b", []string{"b"}, []string{"a", "b"}, []string{"id", "a"}},
		{"select * from t order by a desc", []string{"a desc"}, nil, nil},
		{"select * from t", nil, nil, nil},
		{"select t.id = 1 from t where a > 1", nil, nil, nil},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			sa, err := newSelectAST(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.orderBy, sa.ColumnsInOrderBy())
			assert.Equal(t, tt.whereEqual, sa.EqualPredicateColumnsInWhere())
			assert.Equal(t, tt.projection, sa.ColumnsInProjection())
		})
	}
}
