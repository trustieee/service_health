package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type SystemConfig struct {
	Services []Service `json:"services"`
}

type Service struct {
	Name   string `json:"name"`
	URL    string `json:"url"`
	Status string `json:"status"`
}

func main() {
	log.Println("starting service_health")
	log.Println("reading system config...")
	services := readSystemConfig()

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/health", func(c *gin.Context) {
		for i := range services {
			updateServiceHealth(&services[i])
		}

		c.JSON(http.StatusOK, gin.H{
			"services": services,
		})
	})

	r.Run()
}

func updateServiceHealth(service *Service) {
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	url := service.URL + "/health"
	log.Printf("getting service health for service %s at %s", service.Name, url)
	resp, err := client.Get(url)
	log.Printf("service health for %s is %s", service.Name, service.Status)
	log.Printf("%v \n %v", resp, err)
	if err != nil {
		service.Status = "down"
		return
	}
	defer resp.Body.Close()

	service.Status = "up"
}

func readSystemConfig() []Service {
	bytes, err := os.ReadFile("fake_system.json")
	if err != nil {
		log.Fatal(err)
	}

	var systemConfig SystemConfig
	err = json.Unmarshal(bytes, &systemConfig)
	if err != nil {
		log.Fatal(err)
	}

	return systemConfig.Services
}
