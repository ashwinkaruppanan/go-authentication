package main

import (
	"fmt"
	"log"
	"os"

	"ashwin.com/go-auth/routers"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	router := gin.Default()

	routers.AuthRoutes(router)
	routers.ServiceRoutes(router)
	err := godotenv.Load(".env")
	if err != nil {
		log.Panic(err)
	}
	port := os.Getenv("PORT")
	fmt.Println(port)
	router.Run(":" + port)
}
