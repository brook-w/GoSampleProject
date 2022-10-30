package web

import (
	"Live/live_core"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

func GetLiveStreamOnlineList(ctx *gin.Context) {
	pageNum, _ := strconv.ParseUint(ctx.DefaultQuery("page_num", "1"), 10, 64)
	pageSize, _ := strconv.ParseUint(ctx.DefaultQuery("page_size", "10"), 10, 64)
	streamName := ctx.Query("stream_name")

	streamNameList := make([]string, 0)
	if streamName != "" {
		streamNameList = append(streamNameList, streamName)
	}

	lv := live_core.NewLive()
	res, err := lv.GetLiveStreamOnlineList(pageNum, pageSize, streamNameList...)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, res)
}

type ForbidLiveStreamReq struct {
	StreamName string    `form:"stream_name" binding:"required"`
	ResumeTime time.Time `form:"resume_time" time_format:"2006-01-02 15:04:05" binding:"required"`
	Reason     string    `form:"reason" binding:"required"`
}

func ForbidLiveStream(ctx *gin.Context) {
	req := &ForbidLiveStreamReq{}
	err := ctx.ShouldBind(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}

	res, err := live_core.NewLive().ForbidLiveStream(req.StreamName, req.ResumeTime, req.Reason)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, res)
}
