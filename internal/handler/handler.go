package handler

import (
	"github.com/gin-gonic/gin"
)

type CheckUpdateRequest struct {
	Version  string `json:"version"`
	Platform string `json:"platform"`
}

type CheckUpdateResponse struct {
	Code       int    `json:"code"`
	Message    string `json:"message"`
	FullCDNURL string `json:"full_cdn_url"`
	DiffCDNURL string `json:"diff_cdn_url"`
}

func CheckUpdateHandler(c *gin.Context) {
	var req CheckUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	// Here you would typically check the version against your database or service.
	// For demonstration, we return a mock response.
	resp := CheckUpdateResponse{
		Code:       0,
		Message:    "Success",
		FullCDNURL: "https://example.com/full_update.zip",
		DiffCDNURL: "https://example.com/diff_update.zip",
	}

	c.JSON(200, resp)
}

type FilterIdsRequest struct {
	Version  string `json:"version"`
	Platform string `json:"platform"`
	Channel  string `json:"channel"`
	UserId   string `json:"user_id"`
}

type FilterIdsResponse struct {
	MatchedIdsMap map[string]uint `json:"matched_ids_map"`
}

func FilterIdsHandler(c *gin.Context) {
	var req FilterIdsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	// Here you would typically filter IDs based on the request parameters.
	// For demonstration, we return a mock response.
	resp := FilterIdsResponse{
		MatchedIdsMap: map[string]uint{
			"example_id_1": 1,
			"example_id_2": 2,
		},
	}

	c.JSON(200, resp)
}
