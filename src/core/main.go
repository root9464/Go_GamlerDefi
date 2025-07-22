package main

import (
	_ "github.com/root9464/Go_GamlerDefi/docs"
	core "github.com/root9464/Go_GamlerDefi/src/core/head"
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
