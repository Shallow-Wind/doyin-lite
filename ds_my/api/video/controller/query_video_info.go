package controller

import (
	"github.com/ds_my/api/video"
	"github.com/ds_my/api/video/service"
	"github.com/ds_my/common"
	"github.com/ds_my/common/msg"
	"github.com/ds_my/utils"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"
)

// 获取视频返回值
type getVideoResponse struct {
	common.Response
	NextTime  int64                `json:"next_time"`
	VideoList []video.TheVideoInfo `json:"video_list"`
}

// 获取视频返回值
type videoListResponse struct {
	common.Response
	VideoList []video.TheVideoInfo `json:"video_list"`
}

// GetVideoFeed 获取视频流信息
func GetVideoFeed(c *gin.Context) {
	lastTime, _ := strconv.ParseInt(c.Query("last_time"), 10, 32)
	userIDInterface, success := c.Get("user_id")
	var userID int32
	if success {
		userID = int32(userIDInterface.(int))
	} // 若不存在，userID默认为0

	if lastTime == 0 {
		lastTime = time.Now().Unix()
	}
	// 需要获取NextTime、VideoList
	nextTime, videoInfo, state := service.GetVideoFeed(lastTime, userID)
	if state == 0 {
		c.JSON(http.StatusOK, &getVideoResponse{
			Response: common.Response{
				StatusCode: -1,
				StatusMsg:  msg.HasNoVideoMsg,
			}, NextTime: lastTime,
		})
	} else if state == 1 {
		c.JSON(http.StatusOK, &getVideoResponse{
			Response: common.Response{
				StatusCode: 0,
				StatusMsg:  msg.GetVideoInfoSuccessMsg,
			}, NextTime: nextTime,
			VideoList: videoInfo,
		})
	}
}

func PublishVideo(c *gin.Context) {
	title := c.PostForm("title")
	data, err := c.FormFile("data")
	userID, success := c.Get("user_id")
	if !success {
		c.JSON(http.StatusOK,
			common.Response{
				StatusCode: -1,
				StatusMsg:  msg.TokenParameterAcquisitionError,
			})
		return
	}
	if err != nil {
		c.JSON(http.StatusOK,
			common.Response{
				StatusCode: -1,
				StatusMsg:  msg.PublishVideoFailedMsg,
			})
		return
	}

	fileHandle, err1 := data.Open() //打开上传文件
	if err1 != nil {
		utils.Log.Error("打开文件失败" + err1.Error())
	}

	// 闭包处理错误
	defer func(fileHandle multipart.File) {
		err := fileHandle.Close()
		if err != nil {
			utils.Log.Error("关闭文件错误" + err.Error())
		}
	}(fileHandle)

	fileByte, err2 := ioutil.ReadAll(fileHandle)
	if err2 != nil {
		utils.Log.Error("读取文件错误" + err2.Error())
	}

	if service.PublishVideo(userID.(int), title, fileByte) {
		c.JSON(http.StatusOK,
			common.Response{
				StatusCode: 0,
				StatusMsg:  msg.PublishVideoSuccessMsg,
			})
	} else {
		c.JSON(http.StatusOK,
			common.Response{
				StatusCode: -1,
				StatusMsg:  msg.PublishVideoFailedMsg,
			})
	}
}

// PublicList 登录用户的视频发布列表，直接列出用户所有投稿过的视频
func PublicList(c *gin.Context) {

	userIDStr := c.Query("user_id")
	userID, _ := strconv.ParseInt(userIDStr, 10, 32)
	videoInfo, success2 := service.PublishList(int32(userID))
	if success2 {
		c.JSON(http.StatusOK, &videoListResponse{
			Response: common.Response{
				StatusCode: 0,
				StatusMsg:  msg.GetPublishListSuccessMsg,
			},
			VideoList: videoInfo,
		})
	} else {
		c.JSON(http.StatusOK,
			common.Response{
				StatusCode: -1,
				StatusMsg:  msg.GetPublishListFailedMsg,
			})
	}
}
