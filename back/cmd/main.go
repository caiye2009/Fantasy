package main

import (
	"back/config"
	"log"
)

// @title           后台管理系统 API
// @version         1.0
// @description     后台管理系统的API文档
// @BasePath        /api/v1
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description 请输入 Bearer {token}
func main() {
	log.Println("Starting application...")

	if err := config.Init(); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
