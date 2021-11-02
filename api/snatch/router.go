package snatch

import "github.com/gin-gonic/gin"

func RegisterRedEnvelopeRouter(r *gin.RouterGroup) {
	redenvelopeAPI := r.Group("")
	redenvelopeAPI.POST("", SnatchRedEnvelope) //抢红包
	//redenvelopeAPI.POST("/open", OpenRedEnvelope) //拆红包
	//redenvelopeAPI.POST("/get_wallet_list", GetWalletList) //钱包列表
}
