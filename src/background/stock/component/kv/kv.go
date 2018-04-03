package kv

import (
	"background/common/logger"

	"errors"
	"github.com/jinzhu/gorm"
)

var (
	ErrKeyNotFound = errors.New("KV store key Not Found")
)

func InitKv(db *gorm.DB) error {
	err := initKvStore(db)
	if err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

func GetValueForKey(appId uint32, versionId uint32, key string, db *gorm.DB) (string, error) {
	var kvs []*KvStore
	if err := db.Where("`key`= ? and app_id = ? and version_id = ?", key, appId, versionId).Find(&kvs).Error; err != nil {
		return "", err
	}
	if len(kvs) <= 0 {
		return "", nil
	}
	return kvs[0].Value, nil
}

func SetValueForKey(appId uint32, versionId uint32, key string, value string, db *gorm.DB) error {
	var kvs []*KvStore
	if err := db.Where("`key`= ? and app_id = ? and version_id = ?", key, appId, versionId).Find(&kvs).Error; err != nil {
		return err
	}
	if len(kvs) <= 0 {
		var kv KvStore
		kv.AppId = appId
		kv.VersionId = versionId
		kv.Key = key
		kv.Value = value
		if err := db.Create(&kv).Error; err != nil {
			return err
		}
	} else {
		if err := db.Table(KvStore{}.TableName()).Where("`key`= ? and app_id = ? and version_id = ?", key, appId, versionId).Update("value", value).Error; err != nil {
			return err
		}
	}

	return nil
}

func DeleteValueForKey(appId uint32, versionId uint32, key string, db *gorm.DB) error {
	if err := db.Where("`key`= ? and app_id = ? and version_id = ?", key, appId, versionId).Delete(KvStore{}).Error; err != nil {
		return err
	}
	return nil
}

func DeleteValueForKeyByApp(appId uint32, key string, db *gorm.DB) error {
	if err := db.Where("`key`= ? and app_id = ?", key, appId).Delete(KvStore{}).Error; err != nil {
		return err
	}
	return nil
}

func DeleteValueForKeyByVersion(versionId uint32, key string, db *gorm.DB) error {
	if err := db.Where("`key`= ? and version_id = ?", key, versionId).Delete(KvStore{}).Error; err != nil {
		return err
	}
	return nil
}
