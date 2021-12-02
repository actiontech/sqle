package executor

import (
	"context"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/sirupsen/logrus"
)

// NewMockExecutor returns a new mock executor.
func NewMockExecutor() (*Executor, sqlmock.Sqlmock, error) {
	mockDB, handler, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}
	mockConn, err := mockDB.Conn(context.TODO())
	if err != nil {
		return nil, nil, err
	}

	var executor = &Executor{}
	executor.Db = &BaseConn{
		log:  logrus.WithField("unittest", "unittest"),
		host: "mockhost",
		port: "mockport",
		user: "mockuser",
		db:   mockDB,
		conn: mockConn,
	}
	return executor, handler, nil
}
