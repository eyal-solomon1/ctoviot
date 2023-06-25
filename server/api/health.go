package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Health handles the healthcheck requrests
func (server *Server) Health(ctx *gin.Context) {

	ctx.JSON(http.StatusOK, gin.H{"ok": true})
}
