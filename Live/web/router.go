package web

import "github.com/gin-gonic/gin"

func InitRouterLive(group *gin.RouterGroup) {
	group.GET("/get/live/stream/online/list", GetLiveStreamOnlineList)
	group.POST("/forbid/live/stream", ForbidLiveStream)
}

func InitRouterToools(group *gin.RouterGroup) {
	group.GET("/get/push/url", GetPushUrl)
	group.GET("/get/pull/url", GetPullUrl)
}
