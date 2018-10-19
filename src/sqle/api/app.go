package api

import (
	"fmt"
	"sqle/api/controller"

	"github.com/labstack/echo"
	"github.com/swaggo/echo-swagger"
	_ "sqle/docs"
)

// @title Sqle API Docs
// @version 1.0
// @description This is a sample server for dev.
// @BasePath /
func StartApi(port int, exitChan chan struct{}) {
	e := echo.New()
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.POST("/instances", controller.CreateInst)
	e.POST("/instances/:inst_id/connection", controller.PingInst)
	e.GET("/instances", controller.GetInsts)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%v", port)))
	close(exitChan)
}
