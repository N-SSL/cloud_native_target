package main

import (
	"github.com/N-SSL/container-target/MySQL"
	"github.com/N-SSL/container-target/k8s"
	"github.com/N-SSL/container-target/routes"
	"github.com/gin-gonic/gin"
	"log"
)




func main() {
	k8s.InitClient()
	err := MySQL.ConnectToServer()
	if err != nil {
		return
	}
	r := gin.Default()
	routes.InitRoutes(r)
	log.Fatal(r.Run(":8080"))
}
