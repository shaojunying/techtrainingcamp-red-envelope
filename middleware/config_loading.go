package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"red_envelope/api/redenvelope"
)

//ConfigLoadingMiddleware 获取配置参数信息
func ConfigLoadingMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {
		config, err := redenvelope.Mapper.GetConfigParameters(c)
		if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "code":    -1,
                "message": "获取配置参数信息失败",
            })
            c.Abort()
            return
        }
		c.Set(redenvelope.MaxCountField, config.MaxCount)
		c.Set(redenvelope.ProbabilityField, config.Probability)
		c.Set(redenvelope.BudgetField, config.Budget)
		c.Set(redenvelope.TotalNumberField, config.TotalNumber)
		c.Set(redenvelope.MaxValueField, config.MaxValue)
		c.Set(redenvelope.MinValueField, config.MinValue)

		c.Next()
	}
}
