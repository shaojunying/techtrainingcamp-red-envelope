package redenvelope

import "github.com/gin-gonic/gin"

func RegisterRedEnvelopeRouter(r *gin.RouterGroup) {
	RedEnvelopeAPI := r.Group("")
	RedEnvelopeAPI.POST("/snatch", SnatchRedEnvelope)      //抢红包
	RedEnvelopeAPI.POST("/open", OpenRedEnvelope)          // 拆红包
	RedEnvelopeAPI.POST("/get_wallet_list", GetWalletList) // 获取红包列表
}
