package main

import(
	"em/db"
	"em/config"
	// "net/http"
	"em/handler"
	"github.com/gin-gonic/gin"
)

func main(){
	cfg := config.LoadConfig()
	db.InitDB(cfg)
	router := gin.Default()
	router.POST("/add",handler.AddPersonHandler)
	router.GET("/person", handler.ViewPerson)
	router.POST("/update",handler.UpdatePersonHandler)
	router.POST("/delete",handler.DeletePersonHandler)
	router.Run(":8084")

}
