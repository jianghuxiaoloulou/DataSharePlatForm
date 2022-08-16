package routers

import (
	_ "WowjoyProject/ObjectCloudService/docs"
	v1 "WowjoyProject/ObjectCloudService/internal/routers/api/v1"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func NewRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	// 注册中间件
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	apiv1 := r.Group("/api/v1")
	{
		// 上传单文件
		apiv1.POST("/Object/File/:File", v1.UploadFile)
		// 检查号上传
		apiv1.POST("/Object/AccessNumber/:AccessNumber", v1.UploadNumbers)
		// 全量上传
		apiv1.POST("/Object/All", v1.UploadALL)
		// 单文件下载
		apiv1.GET("/Object/File/:File", v1.DownFile)
		// 检查号下载
		apiv1.GET("/Object/AccessNumber/:AccessNumber", v1.DownNumbers)
		// 全量下载
		apiv1.GET("/Object/All", v1.DownALL)
		// 单文件删除
		apiv1.DELETE("/Object/File/:File", v1.DeleteFile)
		// 检查号删除
		apiv1.DELETE("/Object/AccessNumber/:AccessNumber", v1.DeleteNumbers)
	}
	return r
}
