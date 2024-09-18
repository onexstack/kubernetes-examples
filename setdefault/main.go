package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// 定义请求结构体
type Request struct {
	Limit  int    `form:"limit" json:"limit"`
	Offset int    `form:"offset" json:"offset"`
	Filter string `form:"filter" json:"filter"`
}

// 校验参数并设置默认值的函数
func (req *Request) Validate() error {
	// 设置默认值
	if req.Limit <= 0 {
		req.Limit = 10 // 默认值
	}
	if req.Offset < 0 {
		req.Offset = 0 // 默认值
	}
	if req.Filter == "" {
		req.Filter = "all" // 默认值
	}

	return nil
}

func main() {
	r := gin.Default()

	r.GET("/api/resource", func(c *gin.Context) {
		var req Request

		// 使用 gin.Bind 来解析参数
		if err := c.ShouldBindQuery(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 调用校验函数设置默认值
		if err := req.Validate(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 返回响应
		c.JSON(http.StatusOK, gin.H{"limit": req.Limit, "offset": req.Offset, "filter": req.Filter})
	})

	// 启动服务器
	r.Run(":8080")
}
