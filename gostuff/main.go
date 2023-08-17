package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorcon/rcon"
)

type RconCommand struct {
	Command  string `json:"command"`
	Host     string `json:"host"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {
	r := gin.Default()
	r.POST("/rcon", func(ctx *gin.Context) {
		var command RconCommand
		err := ctx.BindJSON(&command)
		if err != nil {
			ctx.String(http.StatusInternalServerError, "failed to parse received json")
			return
		}
		conn, err := rcon.Dial(command.Host, command.Password)
		if err != nil {
			ctx.String(http.StatusInternalServerError, "failed to connect to rcon host %s", command.Host)
			return
		}
		defer conn.Close()

		response, err := conn.Execute(command.Command)
		if err != nil {
			ctx.String(http.StatusInternalServerError, "failed to execute command %s on rcon host %s", command.Command, command.Host)
		}
		ctx.String(http.StatusOK, "%s", response)

	})
	r.Run("0.0.0.0:8080")
}
