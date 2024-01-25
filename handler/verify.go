package handler

import (
	"frostik.com/auth/constants"
	"frostik.com/auth/controller"
	models "github.com/FrosTiK-SD/mongik/models"
	"github.com/allegro/bigcache/v3"
	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/v2/jwk"
)

type Handler struct {
	MongikClient *models.Mongik
	JwkSet       *jwk.Set
	BigCache     *bigcache.BigCache
}

func (h *Handler) HandlerVerifyStudentIdToken(ctx *gin.Context) {
	idToken := ctx.GetHeader("token")
	noCache := false
	if ctx.GetHeader("cache-control") == constants.NO_CACHE {
		noCache = true
	}

	email, exp, err := controller.VerifyToken(h.BigCache, idToken, h.JwkSet, noCache)

	if err != nil {
		ctx.JSON(200, gin.H{
			"student": nil,
			"expire":  exp,
			"error":   err,
		})
	} else {
		student, err := controller.GetUserByEmail(h.MongikClient, email, &constants.ROLE_STUDENT, noCache)
		ctx.JSON(200, gin.H{
			"data":   student,
			"error":  err,
			"expire": exp,
		})
	}
}

func (h *Handler) InvalidateCache(ctx *gin.Context) {
	h.MongikClient.CacheClient.Delete("GCP_JWKS")
	ctx.JSON(200, gin.H{
		"message": "Successfully invalidated cache",
	})
}
