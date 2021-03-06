package middleware

import (
	"net/http"
	"red_envelope/api/redenvelope"

	"github.com/gin-gonic/gin"
)

func LoadConfig(c *gin.Context) {
	config, err := redenvelope.Mapper.GetConfigParameters(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    -1,
			"message": "获取配置参数信息失败",
		})
		c.Abort()
		return
	}
	//log.Printf("成功获取配置信息。max_count: %d, probability: %g, budget: %d, total_number: %d," +
	//	" max_value: %d, min_value: %d", *config.MaxCount, *config.Probability, *config.Budget, *config.TotalNumber,
	//	*config.MaxValue, *config.MinValue)
	c.Set(redenvelope.MaxCountField, *config.MaxCount)
	c.Set(redenvelope.ProbabilityField, *config.Probability)
	c.Set(redenvelope.BudgetField, *config.Budget)
	c.Set(redenvelope.TotalNumberField, *config.TotalNumber)
	c.Set(redenvelope.MaxValueField, *config.MaxValue)
	c.Set(redenvelope.MinValueField, *config.MinValue)
}

//ConfigLoadingMiddleware 获取配置参数信息
func ConfigLoadingMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {
		LoadConfig(c)
		c.Next()
	}
}
