package ast

import (
	"errors"
	"fmt"
)

type Mappers struct {
	mappers []*Mapper
}

func NewMappers() *Mappers {
	return &Mappers{}
}

func (s *Mappers) AddMapper(ms ...*Mapper) error {
	for _, m := range ms {
		if m == nil {
			return errors.New("can not add null mapper to mappers")
		}
		s.mappers = append(s.mappers, m)
	}
	return nil
}

type StmtInfo struct {
	FilePath  string
	StartLine uint64
	SQL       string
}

func (s *Mappers) GetStmts(skipErrorQuery bool) ([]StmtInfo, error) {
	ctx := NewContext()
	stmts := []StmtInfo{}
	for _, m := range s.mappers {
		for id, node := range m.SqlNodes {
			ctx.Sqls[fmt.Sprintf("%v.%v", m.NameSpace, id)] = node
		}
	}

	for _, m := range s.mappers {
		ctx.DefaultNamespace = m.NameSpace
		stmt, err := m.GetStmts(ctx, skipErrorQuery)
		if err != nil {
			return nil, fmt.Errorf("get sqls from mapper failed, namespace: %v, err: %v", m.NameSpace, err)
		}
		for _, info := range stmt {
			stmts = append(stmts, StmtInfo{
				FilePath:  m.FilePath,
				StartLine: info.StartLine,
				SQL:       info.SQL,
			})
		}
	}
	return stmts, nil
}
