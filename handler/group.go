package handler

import (
	"net/http"

	"github.com/FrosTiK-SD/auth/constants"
	"github.com/FrosTiK-SD/auth/controller"
	"github.com/FrosTiK-SD/auth/interfaces"
	"github.com/FrosTiK-SD/auth/util"
	"github.com/gin-gonic/gin"
)

func (h *Handler) GetAllGroups(ctx *gin.Context) {
	noCache := util.GetNoCache(ctx)
	groups, err := controller.GetAllGroups(h.MongikClient, noCache)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": err,
			"data":  nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":  groups,
		"error": nil,
	})
}

func (h *Handler) BatchEditGroup(ctx *gin.Context) {
	noCache := util.GetNoCache(ctx)

	assignRequests := []interfaces.AssignRequest{}

	if errBinding := ctx.BindJSON(&assignRequests); errBinding != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   constants.ERROR_INCORRENT_BODY,
			"message": errBinding,
		})
		return
	}

	addResult, removeResult, errors := controller.BatchEditGroup(h.MongikClient, assignRequests, noCache)

	if len(*errors) != 0 {
		ctx.JSON(http.StatusPartialContent, gin.H{
			"data": gin.H{
				"addList":    addResult,
				"removeList": removeResult,
			},
			"error": errors,
		})
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"data": gin.H{
				"addList":    addResult,
				"removeList": removeResult,
			},
			"error": nil,
		})
	}

}
