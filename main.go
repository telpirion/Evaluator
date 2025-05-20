package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

type EvalRequest struct {
	Candidate    string `json:"candidate" binding:"required"`
	Test         string `json:"candidate_test" binding:"required"`
	Library      string `json:"library" binding:"required"`
	Language     string `json:"language" binding:"required"`
	UserPrompt   string `json:"user_prompt" binding:"required"`
	Instructions string `json:"instructions"`
}

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "ping!",
		})
	})

	r.POST("/", Eval)

	log.Fatal(r.Run(":8080"))
}

func Eval(c *gin.Context) {
	var req EvalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Fatal(err)
	}

}
