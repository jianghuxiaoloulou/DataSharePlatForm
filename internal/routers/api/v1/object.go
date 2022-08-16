package v1

import (
	"WowjoyProject/ObjectCloudService/global"
	"WowjoyProject/ObjectCloudService/internal/model"
	"WowjoyProject/ObjectCloudService/pkg/app"
	"WowjoyProject/ObjectCloudService/pkg/errcode"

	"github.com/gin-gonic/gin"
)

// 单文件上传
func UploadFile(c *gin.Context) {
	// global.AutoUploadFlag = false
	// global.AutoDowndFlag = false
	// // 成功：
	// app.NewResponse(c).ToResponse(nil)
	// // 失败：
	// app.NewResponse(c).ToErrorResponse(errcode.ServerError)
}

// 检查号上传
func UploadNumbers(c *gin.Context) {
	global.AutoUploadFlag = false
	global.AutoDowndFlag = false

	id := c.Param("AccessNumber")
	global.Logger.Info("需要上传的检查号是：", id)
	if id != "" {
		// 成功：
		app.NewResponse(c).ToResponse(nil)
		// 获取上传任务：
		model.GetObjectData(id, global.UPLOAD)
		// test
		// model.TestUploadeData(global.UPLOAD)
	} else {
		// 失败：
		app.NewResponse(c).ToErrorResponse(errcode.ServerError)
	}
}

// 全量上传
func UploadALL(c *gin.Context) {
	global.Logger.Info("******后台开始全量上传影像******")
	// 成功：
	app.NewResponse(c).ToResponse(nil)
	global.AutoUploadFlag = true
	global.AutoDowndFlag = false
	// 获取上传任务：
	model.AutoUploadObjectData()
	// test
	// model.TestUploadeData(global.UPLOAD)
	// 失败：
	// app.NewResponse(c).ToErrorResponse(errcode.ServerError)
}

// 单文件下载
func DownFile(c *gin.Context) {
	// global.AutoDowndFlag = false
	// global.AutoUploadFlag = false
	// // 成功：
	// app.NewResponse(c).ToResponse(nil)
	// // 失败：
	// app.NewResponse(c).ToErrorResponse(errcode.ServerError)
}

// 检查号下载
func DownNumbers(c *gin.Context) {
	global.AutoDowndFlag = false
	global.AutoUploadFlag = false
	id := c.Param("AccessNumber")
	global.Logger.Info("需要下载的检查号是：", id)
	if id != "" {
		// 成功：
		app.NewResponse(c).ToResponse(nil)
		// 获取下载任务：
		model.GetObjectData(id, global.DOWNLOAD)

		// test
		//model.TestDownData(global.DOWNLOAD)
	} else {
		// 失败：
		app.NewResponse(c).ToErrorResponse(errcode.ServerError)
	}
}

// 全量下载
func DownALL(c *gin.Context) {
	global.Logger.Info("******开始全量下载文件******")
	// 成功：
	app.NewResponse(c).ToResponse(nil)
	global.AutoDowndFlag = true
	global.AutoUploadFlag = false
	// 获取下载任务：
	model.AutoDownObjectData()

	// // test
	// //model.TestDownData(global.DOWNLOAD)
	// // 失败：
	// app.NewResponse(c).ToErrorResponse(errcode.ServerError)
}

// 单文件删除
func DeleteFile(c *gin.Context) {
	// global.AutoUploadFlag = false
	// global.AutoDowndFlag = false
	// // 成功：
	// app.NewResponse(c).ToResponse(nil)
	// // 失败：
	// app.NewResponse(c).ToErrorResponse(errcode.ServerError)
}

// 检查号删除
func DeleteNumbers(c *gin.Context) {
	// id := c.Param("AccessNumber")
	// global.Logger.Info("需要删除的检查号是：", id)
	// if id != "" {
	// 	// 获取下载任务：
	// 	model.GetObjectData(id, global.DELETE)
	// 	// 成功：
	// 	app.NewResponse(c).ToResponse(nil)
	// } else {
	// 	// 失败：
	// 	app.NewResponse(c).ToErrorResponse(errcode.ServerError)
	// }
}
