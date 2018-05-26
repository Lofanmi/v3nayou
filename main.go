package main

import (
	"github.com/Lofanmi/v3nayou/cfg"
	"github.com/Lofanmi/v3nayou/routes"
	// "github.com/gin-gonic/gin"
)

func main() {
	// gin.DisableConsoleColor()

	cfg.LoadEnv()
	cfg.InitTmpl()
	cfg.InitDB()

	r := routes.SetupRouter()

	r.Run("0.0.0.0:5666")
}
