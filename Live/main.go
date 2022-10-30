package main

import (
	"Live/web"
	"github.com/gin-gonic/gin"
	"log"
)

// 基于腾讯云的直播：https://cloud.tencent.com/document/product/267/59019
// API 在线工具：https://console.cloud.tencent.com/api/explorer?Product=cvm&Version=2017-03-12&Action=DescribeZones
// 需要的 API:
//  1. 获取推流地址（鉴权）
//  2. 获取播放地址（鉴权）
//  3. 查询直播中的流
//  4. 禁推流
//  5. 断开流
//  6. 恢复直播流
//  7. 查询流状态
//  8. 获取禁推流列表

// webrtc:webrtc://livepush.brook-w.com/live/123456?txSecret=74d9ac3ca8584643fa8a974ccb6410b4&txTime=635e73b0
func main() {
	r := gin.Default()

	gLvie := r.Group("/v1/live")
	gTools := r.Group("/v1/tools")
	web.InitRouterLive(gLvie)
	web.InitRouterToools(gTools)

	err := r.Run("0.0.0.0:3000")
	if err != nil {
		log.Fatalln(err.Error())
	}
}
