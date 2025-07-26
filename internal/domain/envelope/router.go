package envelope

import "github.com/gin-gonic/gin"

func RegisterRedEnvelopeRouter(r *gin.RouterGroup) {
	RedEnvelopeAPI := r.Group("")
	RedEnvelopeAPI.POST("/snatch", SnatchRedEnvelope)      //抢红包
	RedEnvelopeAPI.POST("/open", OpenRedEnvelope)          // 拆红包
	RedEnvelopeAPI.POST("/get_wallet_list", GetWalletList) // 获取红包列表
}
// 其他接口 不走防作弊逻辑
func RegisterOtherRouter(r *gin.RouterGroup) {
	RedEnvelopeAPI := r.Group("")
	RedEnvelopeAPI.GET("/wrktest", WrkTest)                // 压测
	RedEnvelopeAPI.POST("/config", SetRedEnvelopeConfig)   // 设置红包全局配置
	RedEnvelopeAPI.POST("/catch_koi", CatchKoi)            // 抓锦鲤
}
