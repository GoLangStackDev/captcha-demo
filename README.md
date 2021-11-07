# 基于 Gin + captcha 实现的图片验证码Demo

## 代码实现解说

### 安装库

我们用到了两个库：

```bash
github.com/dchest/captcha
github.com/gin-gonic/gin
```

### 实现了一个工具类

这个工具类我们专门用来处理验证码：

```go
// Captcha 方便后期扩展
type Captcha struct {}

// 单例
var captchaInstance *Captcha
func Instance() *Captcha {
	if captchaInstance==nil {
		captchaInstance = &Captcha{}
	}
	return captchaInstance
}
```

我们声明了一个结构体，方便后期在 captcha 这个库上进行扩展。

```go
// CreateImage 创建图片验证码
func (this *Captcha) CreateImage() string {
	length := captcha.DefaultLen
	captchaId := captcha.NewLen(length)
	return captchaId
}
```

创建验证码也很容易，我们这里直接全部使用他默认的配置，生产6位数的数字验证码，后期有需要可以参考 captcha 库进行调整配置。

这里会返回一个 ID 给我们，这个 ID 就是刚我画的流程图里面的 key，他关联了一个随机数，也就是图片的数字。

**这里他存放在哪里的呢？**

默认是内存，所以重启程序后就可能找不到已经生成的验证码了，但你可以修改他存放在哪里。

```go
// Reload 重载
func (this *Captcha) Reload(captchaId string) bool {
	return captcha.Reload(captchaId)
}
```

因为不可能用户每次都能输对，所以有些时候用户不能识别的情况下就需要进行重新生成随机数，也就是重新生成一张图片，但是 key 也就是 ID 是不能变的，此时就要用到重载。

```go
// Verify 验证
func (this *Captcha) Verify(captchaId,val string) bool {
	return captcha.VerifyString(captchaId, val)
}
```

这就是验证了，传入 ID 和 用户输入的值就可验证了。

```go
// GetImageByte 获取图片二进制流
func (this *Captcha) GetImageByte(captchaId string) []byte {
	var content bytes.Buffer
	err := captcha.WriteImage(&content, captchaId, captcha.StdWidth, captcha.StdHeight)
	if err!=nil {
		log.Println(err)
		return nil
	}
	return content.Bytes()
}
```

最后就是关键了，怎么把图片输出给用户，captcha 库他会生成一个图片的二进制流，你只需要把这个二进制流返回回去即可得到图片。

### Gin部分的代码

这里都只展示关键部分的代码：

```go
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
```

首先是创建的接口，这里直接调用我们工具类的 CreateImage 方法拿到 key 即可。

这里的 URL 和下面这个现实的 API 关联。

```go
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
```

我们每请求一次这个 key 就重载刷新一下他的 Code，方便前端刷新。

前端只需要在这个地址后面加上随机参数即可实现刷新验证码。

最关键的地方就是要设置客户端的请求头里面不能让他缓存。

```go
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
```

最后就是校验了，正常来说这个接口是不能放出来了的，因为：

1、 captcha 库，只要校验一次，不管成功失败他的 ID 就失效了。

2、我们一般都只在业务里面去校验。


### 硬广告
想看更多与 Go 语言相关的资料， 欢迎关注我们的官方公众号：
![GoLang全栈](https://static.golangstack.com/%20upload/qrcode_for_gh_e41ae96a4b33_258.jpg)