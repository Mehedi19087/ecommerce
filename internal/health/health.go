package health

import (
	"encoding/json"
	"fmt"
	"net/http"

	//"time"

	"github.com/gin-gonic/gin"
	//"gorm.io/gorm"
	//"json"
)

func VisitorDivision(ctx *gin.Context) {
	ip := ctx.ClientIP()
	if ip == "" || ip == "127.0.0.1" || ip == "::1" {
		ctx.JSON(500, gin.H{"error": "Failed to get location"})
	}
	resp, err := http.Get(fmt.Sprintf("http://ip-api.com/json/%s", ip))
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to get location"})
		return
	}
	defer resp.Body.Close()

	var data struct {
		Status      string  `json:"status"`
		Country     string  `json:"country"`
		CountryCode string  `json:"countryCode"`
		Region      string  `json:"region"`
		RegionName  string  `json:"regionName"`
		City        string  `json:"city"`
		Zip         string  `json:"zip"`
		Lat         float64 `json:"lat"`
		Lon         float64 `json:"lon"`
		Timezone    string  `json:"timezone"`
		ISP         string  `json:"isp"`
		Org         string  `json:"org"`
		AS          string  `json:"as"`
		Query       string  `json:"query"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to parse location"})
		return
	}

	if data.Status != "success" {
		ctx.JSON(500, gin.H{"error": "Location lookup failed"})
		return
	}

	ctx.JSON(200, gin.H{
		"country":     data.Country,
		"countryCode": data.CountryCode,
		"region":      data.Region,
		"regionName":  data.RegionName,
		"city":        data.City,
		"zip":         data.Zip,
		"lat":         data.Lat,
		"lon":         data.Lon,
		"timezone":    data.Timezone,
		"isp":         data.ISP,
		"org":         data.Org,
		"as":          data.AS,
		"ip":          data.Query,
	})
}

// type HealthController struct {
// 	db *gorm.DB
// }

// func NewHealthController(db *gorm.DB) *HealthController {
// 	return &HealthController{db: db}
// }

// func (h *HealthController) Health(ctx *gin.Context) {
// 	ctx.JSON(http.StatusOK, gin.H{
// 		"status":    "UP",
// 		"timestamp": time.Now().Unix(),
// 	})
// }
