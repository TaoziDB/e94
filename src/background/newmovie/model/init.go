package model

import (
	"background/common/logger"
	"github.com/jinzhu/gorm"
)

func InitModel(db *gorm.DB) error {
	var err error

	err = initMovie(db)
	if err != nil {
		logger.Fatal("Init db movie failed, ", err)
		return err
	}
	err = initTopSearch(db)
	if err != nil {
		logger.Fatal("Init db top_search failed, ", err)
		return err
	}

	err = initAdmin(db)
	if err != nil {
		logger.Fatal("Init db admin failed, ", err)
		return err
	}

	err = initInstallation(db)
	if err != nil {
		logger.Fatal("Init db installation failed, ", err)
		return err
	}

	err = initKvStore(db)
	if err != nil {
		logger.Fatal("Init db kv_store failed, ", err)
		return err
	}

	err = initVideo(db)
	if err != nil {
		logger.Fatal("Init db video failed, ", err)
		return err
	}

	err = initEpisode(db)
	if err != nil {
		logger.Fatal("Init db episode failed, ", err)
		return err
	}

	err = initPlayUrl(db)
	if err != nil {
		logger.Fatal("Init db play_url failed, ", err)
		return err
	}

	err = initRecommend(db)
	if err != nil {
		logger.Fatal("Init db recommend failed, ", err)
		return err
	}

	return err
}

// Do not call this method!!!!
func rebuildModel(db *gorm.DB) {
	dropMovie(db)
	dropTopSearch(db)
	dropAdmin(db)
	dropInstallation(db)
	dropVideo(db)
	dropEpisode(db)
	dropPlayUrl(db)
	dropRecommend(db)

	InitModel(db)
}
