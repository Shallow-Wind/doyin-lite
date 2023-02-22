package main

import (
	favoriteController "github.com/ds_my/api/favorite/controller"
	relationController "github.com/ds_my/api/follow/controller"
	userController "github.com/ds_my/api/user/controller"
	videoController "github.com/ds_my/api/video/controller"
	"github.com/ds_my/common/msg"
	"github.com/ds_my/utils"

	"net/http"

	commentController "github.com/ds_my/api/comment/controller"
	"github.com/ds_my/common"
	"github.com/gin-gonic/gin"
	zhs "github.com/go-playground/validator/v10/translations/zh"
	"github.com/golang-jwt/jwt/v4"
)

var mySecret = []byte(common.MySecret)

// JwtMiddleware jwt中间件
func JwtMiddleware(method string) gin.HandlerFunc {
	return func(c *gin.Context) {
		//从请求头中获取token
		var tokenStr string
		if method == "query" {
			tokenStr = c.Query("token")
		} else if method == "form-data" {
			tokenStr = c.PostForm("token")
		} else if method == "feed" {
			tokenStr = c.Query("token")
		}

		token, err := jwt.ParseWithClaims(tokenStr, &utils.MyClaims{}, func(token *jwt.Token) (interface{}, error) {
			return mySecret, nil
		})
		if err != nil {
			if method == "feed" {
				c.Next()
				return
			}
			if ve, ok := err.(*jwt.ValidationError); ok {
				if ve.Errors&jwt.ValidationErrorMalformed != 0 { //token格式错误
					c.JSON(http.StatusOK, common.Response{StatusCode: -1, StatusMsg: msg.TokenValidationErrorMalformed})
					c.Abort() //阻止执行
					return
				} else if ve.Errors&jwt.ValidationErrorExpired != 0 { //token过期
					c.JSON(http.StatusOK, common.Response{StatusCode: -1, StatusMsg: msg.TokenValidationErrorExpired})
					c.Abort() //阻止执行
					return
				} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 { //token未激活
					c.JSON(http.StatusOK, common.Response{StatusCode: -1, StatusMsg: msg.TokenValidationErrorNotValidYet})
					c.Abort() //阻止执行
					return
				} else {
					c.JSON(http.StatusOK, common.Response{StatusCode: -1, StatusMsg: msg.TokenHandleFailed})
					c.Abort() //阻止执行
					return
				}
			}
		}

		if claims, ok := token.Claims.(*utils.MyClaims); ok && token.Valid {
			id := claims.ID
			c.Set("user_id", id)

			c.Next()
			return
		}
		//失效的token
		c.JSON(http.StatusOK, common.Response{StatusCode: -1, StatusMsg: msg.TokenValid})
		c.Abort() //阻止执行
		return
	}
}

func initRouter(r *gin.Engine) {
	// GRoute总路由组
	GRoute := r.Group("/douyin")
	{
		// user路由组
		user := GRoute.Group("/user")
		{
			user.POST("/register/", userController.RegisterUser)
			user.POST("/login/", userController.LoginUser)
			user.GET("/", JwtMiddleware("query"), userController.GetUserInfo)
		}
		//follow路由组
		relation := GRoute.Group("relation").Use(JwtMiddleware("query"))
		{
			relation.POST("/action/", relationController.RelationAction)
			relation.GET("/follow/list/", relationController.FollowList)
			relation.GET("/follower/list/", relationController.FollowerList)
		}
		//favorite路由组
		favorite := GRoute.Group("/favorite").Use(JwtMiddleware("query"))
		{
			favorite.POST("/action/", favoriteController.FavoriteAction)
			favorite.GET("/list/", favoriteController.FavoriteList)
		}
		//feed获取视频流接口
		GRoute.GET("/feed/", JwtMiddleware("feed"), videoController.GetVideoFeed)
		//publish路由组
		publish := GRoute.Group("/publish")
		{
			publish.POST("/action/", JwtMiddleware("form-data"), videoController.PublishVideo)
			publish.GET("/list/", JwtMiddleware("query"), videoController.PublicList)
		}
		//comment路由组
		comment := GRoute.Group("/comment")
		{
			comment.POST("/action/", JwtMiddleware("query"), commentController.CommentAction)
			comment.GET("/list/", commentController.CommentList)
		}
	}
	// 注册翻译器
	err := zhs.RegisterDefaultTranslations(common.Validate, common.Trans)
	if err != nil {
		utils.Log.Error("翻译器注册错误")
	}
}
