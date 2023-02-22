package service

import (
	"github.com/ds_my/api/video"
	"github.com/ds_my/api/video/repository"
	"github.com/ds_my/common"
	"github.com/ds_my/utils"
	"bytes"
	"fmt"
	"os/exec"
	"sync"
)

func GetVideoFeed(lastTime int64, userID int32) (nextTime int64, videoInfo []video.TheVideoInfo, state int) {
	// state 0:已经没有视频了  1:获取成功  -1:获取失败

	allVideoInfoData, isExist := repository.VideoDao.GetVideoFeed(int32(lastTime))

	if !isExist {
		// 已经没有视频了
		return nextTime, videoInfo, 0
	}

	nextTime = int64(allVideoInfoData[len(allVideoInfoData)-1].Time)
	videoInfo = make([]video.TheVideoInfo, len(allVideoInfoData))

	wg := sync.WaitGroup{}
	wg.Add(len(allVideoInfoData))

	for index, videoInfoData := range allVideoInfoData {
		go func(index int, videoInfo []video.TheVideoInfo, videoInfoData repository.VideoInfo, userID int32) {
			followerCount, followCount, commentCount, favoriteCount := repository.VideoDao.GetVideoInfo(videoInfoData.UserID, videoInfoData.VideoID)
			var isFollow, isFavorite bool
			if userID != 0 {
				isFavorite = repository.VideoDao.QueryIsFavorite(userID, videoInfoData.VideoID)
				isFollow = repository.VideoDao.QueryIsFollow(userID, videoInfoData.UserID)
			}

			videoInfo[index] = video.TheVideoInfo{
				ID: videoInfoData.VideoID,
				Author: video.AuthorInfo{
					ID:            videoInfoData.UserID,
					Name:          videoInfoData.Username,
					FollowCount:   int(followCount),
					FollowerCount: int(followerCount),
					IsFollow:      isFollow,
				},
				PlayURL:       common.OSSPreURL + videoInfoData.PlayURL + ".mp4",
				CoverURL:      common.OSSPreURL + videoInfoData.CoverURL + ".jpg",
				FavoriteCount: int(favoriteCount),
				CommentCount:  int(commentCount),
				IsFavorite:    isFavorite,
				Title:         videoInfoData.Title,
			}
			wg.Done()
		}(index, videoInfo, videoInfoData, userID)
	}
	wg.Wait()
	return nextTime, videoInfo, 1
}

func PublishVideo(userID int, title string, fileBytes []byte) (success bool) {
	node, _ := utils.NewWorker(1)
	randomId := node.GetId()
	fileID := fmt.Sprintf("%v", randomId)
	success = true

	if !utils.UploadFile(fileBytes, fileID, "video") {
		success = false
	}

	// 通过ffmpeg截取视频第一帧为视频封面
	videoURL := common.OSSPreURL + fileID + ".mp4"
	cmd := exec.Command("./ffmpeg", "-i", videoURL, "-vframes", "1", "-f", "image2", "-")
	buf := new(bytes.Buffer)
	cmd.Stdout = buf
	err := cmd.Run()
	if err != nil {
		utils.Log.Error("ffmpeg运行错误 " + err.Error())
		success = false
	}
	// 将视频封面上传至OSS中
	if !utils.UploadFile(buf.Bytes(), fileID, "picture") {
		success = false
	}
	if success {
		if repository.VideoDao.PublishVideo(userID, title, fileID) {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

func PublishList(userID int32) (videoInfo []video.TheVideoInfo, success bool) {

	videoInfoDataList, _ := repository.VideoDao.GetVideoList(userID)

	videoInfo = make([]video.TheVideoInfo, len(videoInfoDataList))

	wg := sync.WaitGroup{}
	wg.Add(len(videoInfoDataList))

	for index, videoInfoData := range videoInfoDataList {
		go func(index int, videoInfo []video.TheVideoInfo, videoInfoData repository.VideoInfo) {
			followerCount, followCount, commentCount, favoriteCount := repository.VideoDao.GetVideoInfo(videoInfoData.UserID, videoInfoData.VideoID)
			isFavorite := repository.VideoDao.QueryIsFavorite(videoInfoData.UserID, videoInfoData.VideoID)
			videoInfo[index] = video.TheVideoInfo{
				ID: videoInfoData.VideoID,
				Author: video.AuthorInfo{
					ID:            videoInfoData.UserID,
					Name:          videoInfoData.Username,
					FollowCount:   int(followCount),
					FollowerCount: int(followerCount),
					IsFollow:      false,
				},
				PlayURL:       common.OSSPreURL + videoInfoData.PlayURL + ".mp4",
				CoverURL:      common.OSSPreURL + videoInfoData.CoverURL + ".jpg",
				FavoriteCount: int(favoriteCount),
				CommentCount:  int(commentCount),
				IsFavorite:    isFavorite,
				Title:         videoInfoData.Title,
			}
			wg.Done()
		}(index, videoInfo, videoInfoData)
	}
	wg.Wait()
	return videoInfo, true
}
