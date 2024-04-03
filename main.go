package main

import (
	"backend/controllers"
	"backend/initializers"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
}

func main() {
	r := gin.Default()

	// Users
	r.POST("/usuario", controllers.CreateUsuario)
	r.GET("/usuarios", controllers.GetUsuarios)
	r.GET("/usuario", controllers.GetUsuarioByUsername)
	r.PUT("/usuario", controllers.UpdateUser)
	r.DELETE("/usuario", controllers.DeleteUser)

	// Puntos Control
	r.POST("/punto", controllers.CreatePuntoControl)
	r.GET("/punto", controllers.GetPuntosControlByOperador)
	r.GET("/puntos", controllers.GetPuntosControl)
	r.PUT("/punto", controllers.UpdatePuntoControl)

	r.Run()
}
