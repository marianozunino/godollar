package route

import (
	"github.com/gin-gonic/gin"
	"github.com/marianozunino/godollar/internal/app/server"
	"github.com/marianozunino/godollar/internal/app/service"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Invoke(registerRoutes),
)

type Todo struct {
	UserId, Id int
	Completed  bool
	Title      string
}

func validateDateFormat(date string) gin.H {
	if len(date) != 10 {
		return gin.H{"error": "invalid date format, expected YYYY-MM-DD"}
	}

	if date[4] != '-' || date[7] != '-' {
		return gin.H{"error": "invalid date format, expected YYYY-MM-DD"}
	}
	return nil
}

func registerRoutes(instance *server.GinHandler, ine service.IneService) {
	instance.Gin.POST("/populate", func(c *gin.Context) {
		if err := ine.Populate(); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{})
	})

	// today rates

	instance.Gin.GET("/rates/today", func(c *gin.Context) {
		data, err := ine.GetToday()
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, data)
	})

	instance.Gin.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{})
	})

	// get specific rate
	instance.Gin.GET("/rates/:date", func(c *gin.Context) {
		date := c.Param("date")
		data, err := ine.GetByDate(date)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, data)
	})

	// get range of rates using query params "from" and "to"
	// if no params are provided, it will return all rates
	// if only "from" is provided, it will return all rates from "from" to today
	// if only "to" is provided, it will return all rates from the beginning to "to"
	instance.Gin.GET("/rates", func(c *gin.Context) {
		from := c.Query("from")
		to := c.Query("to")

		if from == "" && to == "" {
			data, err := ine.GetAll()
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
			}
			c.JSON(200, data)
			return
		}

		if from == "" {

			response := validateDateFormat(to)
			if response != nil {
				c.JSON(400, response)
				return
			}

			data, err := ine.GetByDateRange("", to)
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
			}
			c.JSON(200, data)
			return
		}

		if to == "" {
			response := validateDateFormat(from)
			if response != nil {
				c.JSON(400, response)
				return
			}

			data, err := ine.GetByDateRange(from, "")
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
			}
			c.JSON(200, data)
			return
		}

		data, err := ine.GetByDateRange(from, to)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
		}
		c.JSON(200, data)

	})

}
