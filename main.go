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
	r.GET("/users/:uid", middlewares.RequireAuth, controllers.GetUserAndRoleByUid) //Tested
	r.PUT("/users/:uid", middlewares.RequireAuth, controllers.UpdateUser)
	r.DELETE("/users/:uid", middlewares.RequireAuth, controllers.DeleteUser)
	r.GET("/users/:uid/groups", middlewares.RequireAuth, controllers.GetGroupsOfUser)

	//// GetAllUsersByRoleId GetAllGroupsOfUser GetUserAndRoleByUid FindUsersByQueryParameters

	r.POST("/groups", middlewares.RequireAuth, controllers.CreateGroup)
	r.GET("/groups", middlewares.RequireAuth, controllers.GetAllGroups) // Tested
	r.GET("/groups/:id", middlewares.RequireAuth, controllers.GetGroup)
	r.PUT("/groups/:id", middlewares.RequireAuth, controllers.UpdateGroup)
	r.DELETE("/groups/:id", middlewares.RequireAuth, controllers.DeleteGroup)

	r.POST("/roles", middlewares.RequireAuth, controllers.CreateRole)
	r.GET("/roles", middlewares.RequireAuth, controllers.GetAllRoles) // Tested
	//r.GET("/roles/:id", middlewares.RequireAuth, controllers.GetRole)
	r.PUT("/roles/:id", middlewares.RequireAuth, controllers.UpdateRole)
	r.DELETE("/roles/:id", middlewares.RequireAuth, controllers.DeleteRole)
	r.GET("/roles/:roleId/users", middlewares.RequireAuth, controllers.GetUsersByRole)

	r.POST("/login", controllers.Login)
	r.POST("/signup", controllers.Signup)
	r.POST("/logout", middlewares.RequireAuth, controllers.Logout)

	if err := r.Run(os.Getenv("API-PORT")); err != nil {
		log.Fatal(err)
	}
}
