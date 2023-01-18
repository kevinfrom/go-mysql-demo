package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type JsonApiResponse[T any] struct {
	Data T `json:"data"`
}

func SetupServer() *gin.Engine {
	server := gin.Default()
	server.SetTrustedProxies([]string{})

	server.GET("/products", func(c *gin.Context) {
		c.IndentedJSON(http.StatusOK, JsonApiResponse[[]Product]{
			Data: GetProducts(),
		})
	})

	server.POST("/products", func(c *gin.Context) {
		var postData struct {
			Name  *string `json:"name"`
			Price *int64  `json:"price"`
		}

		c.BindJSON(&postData)

		if postData.Name == nil {
			c.IndentedJSON(http.StatusBadRequest, JsonApiResponse[string]{
				Data: "\"name\" is required",
			})
			return
		}

		if postData.Price == nil {
			c.IndentedJSON(http.StatusBadRequest, JsonApiResponse[string]{
				Data: "\"price\" is required",
			})
			return
		}

		product := CreateProduct(*postData.Name, *postData.Price)

		c.IndentedJSON(http.StatusCreated, JsonApiResponse[Product]{
			Data: product,
		})
	})

	server.GET("/products/:id", func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 0)
		CheckError(err)

		product, err := GetProduct(id)

		if err != nil {
			c.IndentedJSON(http.StatusNotFound, JsonApiResponse[string]{
				Data: "product not found",
			})
		} else {
			c.IndentedJSON(http.StatusOK, JsonApiResponse[Product]{
				Data: product,
			})
		}
	})

	server.PATCH("/products/:id", func(c *gin.Context) {
		var postData struct {
			Name  *string `json:"name,omitempty"`
			Price *int64  `json:"price,omitempty"`
		}

		c.BindJSON(&postData)

		id, err := strconv.ParseUint(c.Param("id"), 10, 0)
		CheckError(err)

		product, err := GetProduct(id)

		if err != nil {
			c.IndentedJSON(http.StatusNotFound, JsonApiResponse[string]{
				Data: "product not found",
			})
		} else {
			if postData.Name != nil {
				product.Name = *postData.Name
			} else if postData.Price != nil {
				product.Price = *postData.Price
			} else {
				c.IndentedJSON(http.StatusBadRequest, JsonApiResponse[string]{
					Data: "\"name\" or \"price\" is required",
				})
			}

			_, err := SaveProduct(product)
			CheckError(err)

			c.IndentedJSON(http.StatusOK, JsonApiResponse[Product]{
				Data: product,
			})
		}
	})

	server.DELETE("/products/:id", func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 0)
		CheckError(err)

		_, err = DeleteProduct(id)

		if err != nil {
			c.IndentedJSON(http.StatusNotFound, JsonApiResponse[string]{
				Data: "product not found",
			})
		} else {
			c.IndentedJSON(http.StatusOK, JsonApiResponse[string]{
				Data: "product deleted",
			})
		}
	})

	return server
}
