package grpc

import (
	"actiontech/ucommon/log"
	"fmt"

	"google.golang.org/grpc/grpclog"
)

func init() {
	// set grpc logger
	grpclog.SetLogger(myGrpcLogger)
}

type grpcLogger struct{}

// myGrpcLogger implement grpc's Logger interface. It will write logs to detail.log.
var myGrpcLogger = newGrpcLogger()

func newGrpcLogger() grpclog.Logger {
	return new(grpcLogger)
}

func (logger *grpcLogger) Fatal(args ...interface{}) {
	log.DetailDilute2(log.NewStage().Enter("grpc"), fmt.Sprint(args...), fmt.Sprint(args...))
}

func (logger *grpcLogger) Fatalf(format string, args ...interface{}) {
	log.DetailDilute2(log.NewStage().Enter("grpc"), fmt.Sprintf(format, args...), format, args...)
}

func (logger *grpcLogger) Fatalln(args ...interface{}) {
	log.DetailDilute2(log.NewStage().Enter("grpc"), fmt.Sprint(args...), fmt.Sprint(args...))
}

func (logger *grpcLogger) Print(args ...interface{}) {
	log.DetailDilute2(log.NewStage().Enter("grpc"), fmt.Sprint(args...), fmt.Sprint(args...))
}

func (logger *grpcLogger) Printf(format string, args ...interface{}) {
	log.DetailDilute2(log.NewStage().Enter("grpc"), fmt.Sprintf(format, args...), format, args...)
}

func (logger *grpcLogger) Println(args ...interface{}) {
	log.DetailDilute2(log.NewStage().Enter("grpc"), fmt.Sprint(args...), fmt.Sprint(args...))
}
