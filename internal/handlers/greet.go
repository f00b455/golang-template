package handlers

import (
	"net/http"

	"github.com/f00b455/golang-template/pkg/shared"
	"github.com/gin-gonic/gin"
)

// GreetHandler handles greeting requests.
type GreetHandler struct{}

// NewGreetHandler creates a new GreetHandler.
func NewGreetHandler() *GreetHandler {
	return &GreetHandler{}
}

// GreetResponse represents the response for the greet endpoint.
type GreetResponse struct {
	Message string `json:"message" example:"Hello, World!"`
}

// Greet handles GET /api/greet
// @Summary      Greet endpoint
// @Description  Returns a greeting message
// @Tags         greet
// @Accept       json
// @Produce      json
// @Param        name    query     string  false  "Name to greet" default(World)
// @Success      200     {object}  GreetResponse
// @Router       /greet [get]
func (h *GreetHandler) Greet(c *gin.Context) {
	name := c.DefaultQuery("name", "World")
	message := shared.Greet(name)

	c.JSON(http.StatusOK, GreetResponse{
		Message: message,
	})
}
