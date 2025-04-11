package main

import (
	core "github.com/root9464/Go_GamlerDefi/core/head"
	_ "github.com/root9464/Go_GamlerDefi/docs"
)

// @title			GamlerDefi API
// @version		1.0
// @description	API for GamlerDefi
// @host			localhost:6069
// @BasePath		/
func main() {
	app := core.InitApp()
	app.Start()
}
