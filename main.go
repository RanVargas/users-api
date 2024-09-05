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
			controllers.GetAllUsers(ctx) //Tested
		} else {
			controllers.FindUsersByQueryParams(ctx, search, l, orderBy) //Tested
		}

	})
	r.GET("/users/:uid", middlewares.RequireAuth, controllers.GetUserAndRoleByUid)         //Tested
	r.PUT("/users/:uid", middlewares.RequireAuth, controllers.UpdateUser)                  //Broken
	r.DELETE("/users/:uid", middlewares.RequireAuth, controllers.DeleteUser)               //Tested
	r.GET("/users/:uid/groups", middlewares.RequireAuth, controllers.GetGroupsOfUser)      //Broken
	r.PUT("/users/:uid/password", middlewares.RequireAuth, controllers.UpdateUserPassword) //Tested

	r.POST("/groups", middlewares.RequireAuth, controllers.CreateGroup)        //Tested
	r.GET("/groups", middlewares.RequireAuth, controllers.GetAllGroups)        //Tested
	r.GET("/groups/:uid", middlewares.RequireAuth, controllers.GetGroup)       //Tested
	r.PUT("/groups/:uid", middlewares.RequireAuth, controllers.UpdateGroup)    //Broken
	r.DELETE("/groups/:uid", middlewares.RequireAuth, controllers.DeleteGroup) //Tested

	r.POST("/roles", middlewares.RequireAuth, controllers.CreateRole)               //Tested
	r.GET("/roles", middlewares.RequireAuth, controllers.GetAllRoles)               //Tested
	r.GET("/roles/:uid", middlewares.RequireAuth, controllers.GetRole)              //Tested
	r.PUT("/roles/:uid", middlewares.RequireAuth, controllers.UpdateRole)           //Broken
	r.DELETE("/roles/:uid", middlewares.RequireAuth, controllers.DeleteRole)        //Tested
	r.GET("/roles/:uid/users", middlewares.RequireAuth, controllers.GetUsersByRole) //Tested

	r.POST("/login", controllers.Login)   //Tested
	r.POST("/signup", controllers.Signup) //Tested
	r.POST("/logout", middlewares.RequireAuth, controllers.Logout)

	if err := r.Run(os.Getenv("API_SERVING_ADDRESS") + ":" + os.Getenv("API_PORT")); err != nil {
		log.Fatal(err)
	}
}
