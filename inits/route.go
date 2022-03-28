package inits

import (
	"github.com/gin-gonic/gin"
	"hitokoto-go/handlers/public"
	"hitokoto-go/routers"
)

func r(e *gin.Engine) {

	// Public API
	publicEndpoint := e.Group("/")
	{
		publicEndpoint.GET("/", public.GetHitokoto)
	}

	// Commit API
	// todo

	// Admin API
	// todo

}

func Routes() *gin.Engine {
	routers.Include(r)
	return routers.Init()
}
