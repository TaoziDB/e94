package apimodel

import (
	"background/newmovie/model"
	"background/common/logger"
	"github.com/jinzhu/gorm"
	"background/common/constant"
	"background/newmovie/service"
)

type PlayUrl struct {
	Id             uint32  `json:"id"`
	Provider       uint32  `json:"provider"`
	Url            string  `json:"url"`

}
type Video struct {
	Id       uint32  `json:"id"`
	Title    string  `json:"title"`
	Description    uint32  `json:"descriptiono"`
	Score    uint32  `json:"score"`
	ThumbX      string  `json:"thumb_x"`
	ThumbY uint32  `json:"thumb_y"`
	PublishDate    string  `json:"publish_date"`
	Year bool    `json:"year"`
	Language bool    `json:"language"`
	Country bool    `json:"country"`
	Directors bool    `json:"directors"`
	Actors bool    `json:"actors"`
	Tags bool    `json:"tags"`
	Urls     []*PlayUrl `json:"urls"`
}

func VideoFromDb(jsCode string,src *model.Video,db *gorm.DB) *Video {
	dst := Video{}
	dst.Id = src.Id
	dst.Title = src.Title
	dst.Description = src.Description
	dst.Score = src.Score
	dst.ThumbX = src.ThumbX
	dst.ThumbY = src.ThumbY
	dst.PublishDate = src.PublishDate
	dst.Year = src.Year
	dst.Language = src.Language
	dst.Country = src.Country
	dst.Directors = src.Directors
	dst.Actors = src.Actors
	dst.Tags = src.Tags

	var episode model.Episode
	if err := db.Where("video_id = ?",src.Id).First(&episode).Error ; err != nil{
		logger.Error(err)
		return nil
	}

	var playUrls []model.PlayUrl
	if err := db.Where("content_type = ? and content_id = ?",constant.MediaTypeEpisode,episode.Id).First(&playUrls).Error ; err != nil{
		logger.Error(err)
		return nil
	}

	for _,playUrl := range playUrls{
		var pUrl PlayUrl
		pUrl.Id = playUrl.Id
		pUrl.Provider = playUrl.Provider
		pUrl.Url = service.GetRealUrl(playUrl.Provider,playUrl.Url,jsCode)
		dst.Urls = append(dst.Urls,&pUrl)
	}
	return &dst
}
