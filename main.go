package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// type uriBinding struct {
// 	ID string `uri:"id"`
// }

type Product struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Customer struct {
	Email         string `json:"email" binding:"required,email"`
	Password      string `json:"password" binding:"required,password"`
	Role          string `json:"role" binding:"required,oneof=BASIC ADMIN"`
	StreetAddress string `json:"streetAddress"`
	StreetNumber  string `json:"streetNumber" binding:"required_with=StreetAddress"`
}

func verifyPassword(fl validator.FieldLevel) bool {
	regex := regexp.MustCompile(`\\w{8,}`)
	password := fl.Field().String()
	return regex.MatchString(password)
}

func main() {
	router := gin.Default()
	address := ":3004"

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("password", verifyPassword)
	}

	// v1 := router.Group("/api/v1")

	accounts := map[string]string{
		"john": "doe",
		"foo":  "bar",
	}

	authMiddleware := gin.BasicAuth(accounts)

	// router.Use(func(c *gin.Context) {
	// 	fmt.Println("logging...")
	// })

	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "ping")
	})

	router.GET("/products/:id", func(c *gin.Context) {
		id := c.Param("id")
		fmt.Println("ID: ", id)
		c.String(http.StatusOK, "hello world")
	})

	router.POST("/products", authMiddleware, func(c *gin.Context) {
		var product Product

		if e := c.ShouldBindJSON(&product); e != nil {
			c.String(http.StatusBadRequest, e.Error())
			return
		}

		fmt.Println("Binding: ", product)
		c.String(http.StatusOK, product.ID)
	})

	router.POST("/customers", func(c *gin.Context) {
		var customer Customer

		if e := c.ShouldBindJSON(&customer); e != nil {
			c.String(http.StatusBadRequest, e.Error())
			return
		}

		fmt.Println("Binding: ", customer)
		c.JSON(http.StatusOK, customer)
	})

	log.Fatalln(router.Run((address)))
}
