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
	r.GET("/login", controllers.Login)
	r.PUT("/usuario", controllers.UpdateUser)
	r.DELETE("/usuario", controllers.DeleteUser)

	// Puntos Control
	r.POST("/punto", controllers.CreatePuntoControl)
	r.GET("/punto", controllers.GetPuntosControlByOperador)
	r.GET("/puntos", controllers.GetPuntosControl)
	r.PUT("/punto", controllers.UpdatePuntoControl)

	// Rutas
	r.POST("/ruta", controllers.CreateRuta)
	r.POST("/ruta/punto", controllers.AddPuntoControl)
	r.PUT("/ruta", controllers.UpdateRuta)
	r.GET("/rutas", controllers.GetRutas)
	r.GET("/ruta/puntos", controllers.GetPuntosControlRuta)
	r.DELETE("/ruta/punto", controllers.RemovePuntoControl)

	r.Run()
}
