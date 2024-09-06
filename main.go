package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"strconv"
	"users-api/controllers"
	"users-api/database"
	"users-api/middlewares"
)

func init() {
	database.DBInit()
	controllers.InitializeUsersRepo()
	controllers.InitializeRolesRepo()
	controllers.InitializeGroupsRepo()
	controllers.InitializeUserPasswordRepo()
}

func main() {

	r := gin.Default()

	r.POST("/users", middlewares.RequireAuth, controllers.CreateUser)
	r.GET("/users", middlewares.RequireAuth, func(ctx *gin.Context) {
		search := ctx.Query("searchTerm")
		limitParam := ctx.DefaultQuery("limit", "10")
		l, _ := strconv.Atoi(limitParam)
		orderBy := ctx.DefaultQuery("orderBy", "id")
		if search == "" {
			controllers.GetAllUsers(ctx)
		} else {
			controllers.FindUsersByQueryParams(ctx, search, l, orderBy)
		}

	})
	r.GET("/users/:uid", middlewares.RequireAuth, controllers.GetUserAndRoleByUid)
	r.PUT("/users/:uid", middlewares.RequireAuth, controllers.UpdateUser)
	r.DELETE("/users/:uid", middlewares.RequireAuth, controllers.DeleteUser)
	r.GET("/users/:uid/groups", middlewares.RequireAuth, controllers.GetGroupsOfUser)
	r.PUT("/users/:uid/password", middlewares.RequireAuth, controllers.UpdateUserPassword)

	r.POST("/groups", middlewares.RequireAuth, controllers.CreateGroup)
	r.GET("/groups", middlewares.RequireAuth, controllers.GetAllGroups)
	r.GET("/groups/:uid", middlewares.RequireAuth, controllers.GetGroup)
	r.PUT("/groups/:uid", middlewares.RequireAuth, controllers.UpdateGroup)
	r.DELETE("/groups/:uid", middlewares.RequireAuth, controllers.DeleteGroup)

	r.POST("/roles", middlewares.RequireAuth, controllers.CreateRole)
	r.GET("/roles", middlewares.RequireAuth, controllers.GetAllRoles)
	r.GET("/roles/:uid", middlewares.RequireAuth, controllers.GetRole)
	r.PUT("/roles/:uid", middlewares.RequireAuth, controllers.UpdateRole)
	r.DELETE("/roles/:uid", middlewares.RequireAuth, controllers.DeleteRole)
	r.GET("/roles/:uid/users", middlewares.RequireAuth, controllers.GetUsersByRole)

	r.POST("/login", controllers.Login)
	r.POST("/signup", controllers.Signup)
	r.POST("/logout", middlewares.RequireAuth, controllers.Logout)

	if err := r.Run(os.Getenv("API_SERVING_ADDRESS") + ":" + os.Getenv("API_PORT")); err != nil {
		log.Printf("Failed to start server: %v", err)
		log.Fatal(err)
	}
}
