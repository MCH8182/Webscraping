package api

import (
	"goquery-example/api/controller"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RunAPI() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.Use(unseenPanicHandler)

	configureSubRouters(r)

	err := r.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}

func configureSubRouters(engine *gin.Engine) {
	controller.StartUserRouter(engine)
	// TODO: tambah router news
	controller.StartNewsRouter(engine)
}

func unseenPanicHandler(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			er, _ := err.(error)
			log.Println(er)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unexpected error occured"})
		}
	}()
	c.Next()
}
