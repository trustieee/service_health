package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	var services []*gin.Engine
	for i := 1; i < 5; i++ {
		go func(i int) {
			if i%2 == 0 {
				return
			}

			r := gin.Default()
			services = append(services, r)
			r.GET("/health", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"message": "pong",
				})
			})

			r.Run(fmt.Sprintf(":808%d", i))
		}(i)
	}

	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
