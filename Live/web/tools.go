package web

import (
	"Live/live_core"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetPushUrl 获取推流地址
func GetPushUrl(ctx *gin.Context) {
	streamName := ctx.Query("stream_name")
	if streamName == "" {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": errors.New("直播流名称不能为空").Error(),
		})
		return
	}
	lv := live_core.NewLive()
	url := lv.GetPushUrl(streamName)
	ctx.JSON(http.StatusOK, url)
}

// GetPullUrl 获取拉流地址
func GetPullUrl(ctx *gin.Context) {
	streamName := ctx.Query("stream_name")
	if streamName == "" {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": errors.New("直播流名称不能为空").Error(),
		})
		return
	}
	lv := live_core.NewLive()
	url := lv.GetPullUrl(streamName)
	ctx.JSON(http.StatusOK, url)
}
