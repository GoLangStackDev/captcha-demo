package main

import (
	captcha "captcha-demo/lib"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	r := gin.New()
	// 创建
	// 这里方便看到效果 我用的 GET 请求，实际生产最好不要用 GET
	r.Handle("GET", "/captcha/create", func(c *gin.Context) {
		imgId := captcha.Instance().CreateImage()
		c.JSON(http.StatusOK,
			gin.H{
				"code": 200,
				"key": imgId,
				"url": "/captcha/img/"+imgId,
			})
	})
	// 现实图片
	r.Handle("GET", "/captcha/img/:key", func(c *gin.Context) {
		captchaId := c.Param("key")
		c.Writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		c.Writer.Header().Set("Pragma", "no-cache")
		c.Writer.Header().Set("Expires", "0")
		c.Writer.Header().Set("Content-Type", "image/png")
		// 重载一次
		captcha.Instance().Reload(captchaId)
		// 输出图片
		c.Writer.Write(captcha.Instance().GetImageByte(captchaId))
	})
	// 校验
	r.Handle("GET", "/captcha/verify/:key/:val", func(c *gin.Context) {
		captchaId := c.Param("key")
		val := c.Param("val")
		if captcha.Instance().Verify(captchaId,val) {
			c.JSON(http.StatusOK, gin.H{"code": 200})
		}else{
			c.JSON(http.StatusOK, gin.H{"code": 400})
		}
	})

	r.Run(":8083")
}


