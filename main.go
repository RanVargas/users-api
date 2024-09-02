package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"users-api/controllers"
	"users-api/database"
	"users-api/middlewares"
)

func init() {
	database.DBInit()
}

func main() {

	r := gin.Default()

	r.POST("/users", middlewares.RequireAuth, controllers.CreateUser)
	r.GET("/users", middlewares.RequireAuth, controllers.GetAllUsers)
	r.GET("/users/:id", middlewares.RequireAuth, controllers.GetUser)
	r.PUT("/users/:id", middlewares.RequireAuth, controllers.UpdateUser)
	r.DELETE("/users/:id", middlewares.RequireAuth, controllers.DeleteUser)

	r.POST("/groups", middlewares.RequireAuth, controllers.CreateGroup)
	r.GET("/groups", middlewares.RequireAuth, controllers.GetAllGroups)
	r.GET("/groups/:id", middlewares.RequireAuth, controllers.GetGroup)
	r.PUT("/groups/:id", middlewares.RequireAuth, controllers.UpdateGroup)
	r.DELETE("/groups/:id", middlewares.RequireAuth, controllers.DeleteGroup)

	r.POST("/roles", middlewares.RequireAuth, controllers.CreateRole)
	r.GET("/roles", middlewares.RequireAuth, controllers.GetAllRoles)
	r.GET("/roles/:id", middlewares.RequireAuth, controllers.GetRole)
	r.PUT("/roles/:id", middlewares.RequireAuth, controllers.UpdateRole)
	r.DELETE("/roles/:id", middlewares.RequireAuth, controllers.DeleteRole)

	r.POST("/login", controllers.Login)
	r.POST("/signup", controllers.Signup)
	r.POST("/logout", middlewares.RequireAuth, controllers.Logout)

	if err := r.Run(os.Getenv("API-PORT")); err != nil {
		log.Fatal(err)
	}
}
