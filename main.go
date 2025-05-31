package main

import(
	"em/db"
	"em/config"
	"em/handler"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/files"
	_ "em/docs"
)
// @title           People API
// @version         1.0
// @description     REST API for managing people
// @host            localhost:8084
// @BasePath        /
func main(){
	cfg := config.LoadConfig()
	db.InitDB(cfg)
	router := gin.Default()
	router.POST("/add",handler.AddPersonHandler)
	router.GET("/person", handler.ViewPerson)
	router.POST("/update",handler.UpdatePersonHandler)
	router.POST("/delete",handler.DeletePersonHandler)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Run(":8084")

}
