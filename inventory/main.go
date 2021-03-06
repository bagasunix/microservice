package main

import (
	"inventory/router"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowMethods:    []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:    []string{"Origin", "Authorization", "Content-type"},
		AllowAllOrigins: true,
	}))

	r.GET("/", handleTest)

	r.POST("/food", router.InputFood)
	r.GET("/food", router.GetFood)
	r.PATCH("/food/:id", router.EditFood)
	r.DELETE("/food/:id", router.DeleteFood)
	r.GET("/food/:id", router.GetFoodById)

	r.Run("0.0.0.0:" + os.Getenv("BIND_ADDR")) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func handleTest(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{
		"message": "Status OK!, test berjalan",
	})
}
