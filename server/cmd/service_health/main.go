package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type SystemConfig struct {
	Services []Service `json:"services"`
}

type Service struct {
	LastUpdated time.Time
	Name        string `json:"name"`
	URL         string `json:"url"`
	Status      string `json:"status"`
}

type ResponseService struct {
	Name   string `json:"name"`
	URL    string `json:"url"`
	Status string `json:"status"`
}

func (s *Service) ToResponseService() ResponseService {
	return ResponseService{s.Name, s.URL, s.Status}
}

func main() {
	services := readSystemConfig()

	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		var wg sync.WaitGroup
		for i := range services {
			wg.Add(1)
			go updateServiceHealth(&services[i], &wg)
		}

		wg.Wait()

		responseServices := make([]ResponseService, 0, len(services))
		for _, service := range services {
			responseServices = append(responseServices, service.ToResponseService())
		}
		c.JSON(http.StatusOK, gin.H{
			"services": responseServices,
		})
	})

	r.Run(":8080")
}

func updateServiceHealth(service *Service, wg *sync.WaitGroup) {
	defer wg.Done()

	if time.Since(service.LastUpdated) < 5*time.Second {
		return
	}

	client := http.Client{
		Timeout: 5 * time.Second,
	}

	url := service.URL + "/health"
	resp, err := client.Get(url)
	if err != nil {
		service.Status = "down"
		return
	}
	defer resp.Body.Close()

	service.Status = "up"
	service.LastUpdated = time.Now()
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
