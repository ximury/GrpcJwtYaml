//package main
package main

import (
	"log"
	"net/http"
	"os"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

//type login struct {
//	Username string `form:"username" json:"username" binding:"required"`
//	Password string `form:"password" json:"password" binding:"required"`
//}

type supportedProtocols struct {
	Protocols []string `json:"protocols"`
}

var identityKey = "id"

func CreateRouter(
	logFile string,
	payloadFunc func(data interface{}) jwt.MapClaims,
	authenticator func(c *gin.Context) (interface{}, error),
	identityHandler func(c *gin.Context) interface{},
	authorizator func(data interface{}, c *gin.Context) bool) (*gin.Engine, *jwt.GinJWTMiddleware) {

	gin.SetMode(gin.ReleaseMode)
	gin.DisableConsoleColor()
	file, _ := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	gin.DefaultWriter = file
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// the jwt middleware
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		//中间件名称
		Realm:       "test zone",
		//Secret key used for signing.
		Key:         []byte("secret key"),
		//token过期时间
		Timeout:     time.Hour,
		//token刷新最大时间
		MaxRefresh:  time.Hour,
		//身份验证的key值
		IdentityKey: identityKey,
		//登录期间的回调的函数
		PayloadFunc: payloadFunc,
		//解析并设置用户身份信息
		IdentityHandler: identityHandler,
		//根据登录信息对用户进行身份验证的回调函数
		Authenticator: authenticator,
		//接收用户信息并编写授权规则
		Authorizator: authorizator,
		//自定义处理未进行授权的逻辑
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"status_code": code,
				"reason": message,
			})
		},
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		TokenLookup: "header: Authorization, query: token, cookie: jwt",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})

	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	r.POST("/login", authMiddleware.LoginHandler)

	r.NoRoute(authMiddleware.MiddlewareFunc(), func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		log.Printf("NoRoute claims: %#v\n", claims)
		c.JSON(404, gin.H{
			"status_code": "PAGE_NOT_FOUND", "reason": "Page not found"})
	})

	auth := r.Group("/")
	// Refresh time can be longer than token timeout
	auth.GET("/refresh_token", authMiddleware.RefreshHandler)
	auth.Use(authMiddleware.MiddlewareFunc())
	{
		auth.GET("/supportedprotocols", func(c *gin.Context) {
			response := supportedProtocols{}
			response.Protocols = append(response.Protocols, "TCP")
			c.IndentedJSON(http.StatusOK, response)
		})
	}
	return r, authMiddleware
}

//func RunServer(addr string, r *gin.Engine) {
//	if err := http.ListenAndServe(addr, r); err != nil {
//		log.Fatal(err)
//	}
//}
