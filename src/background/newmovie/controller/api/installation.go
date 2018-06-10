package api

import (
	"background/newmovie/model"
	"background/common/constant"
	"background/common/logger"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

/*
	POST /cms/v1.0/installation
	配置App，获取App初始化参数
	@Author: HYK
	http://localhost:2000/#!./cms/api-config.md
*/
func InstallationHandler(c *gin.Context) {
	type param struct {
		InstallationId uint64  `json:"installation_id"`
		DeviceId       string  `json:"device_id" binding:"required"`
		MacAddress     string  `json:"mac_address" binding:"required"`
		Imei           string  `json:"imei"`
		OsVersion     string  `json:"os_version"`
		Product   string  `json:"product"` //产品名称
		Model   string  `json:"model"` //设备型号
		Brand   string  `json:"brand"` //设备品牌
		Carrier   uint8  `json:"carrier"` //电话类型

	//	CarrierTypeUnknown      = 0 // 未知类型
	//CarrierTypeChinaMobile  = 1 // 中国移动
	//CarrierTypeChinaUnicom  = 2 // 中国联通
	//CarrierTypeChinaTelecom = 3 // 中国电信

	}
	var p param
	var err error

	if err = c.Bind(&p); err != nil {
		logger.Error(err)
		return
	}


	db := c.MustGet(constant.ContextDb).(*gorm.DB)

	var dbInstall model.Installation
	if p.InstallationId != 0 {
		db.Where("id=?", p.InstallationId).First(&dbInstall)
	}
	if dbInstall.Id == 0 {
		if err = db.Where(" device_id = ? AND mac_address = ? AND imei = ?", p.DeviceId, p.MacAddress, p.Imei).First(&dbInstall).Error; err != nil && err != gorm.ErrRecordNotFound {
			logger.Error(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}

	dbInstall.OsVersion = p.OsVersion
	dbInstall.DeviceId = p.DeviceId
	dbInstall.MacAddress = p.MacAddress
	dbInstall.Imei = p.Imei
	dbInstall.Carrier = p.Carrier
	dbInstall.Brand = p.Brand
	dbInstall.Product = p.Product
	dbInstall.DeviceModel = p.Model
	//dbInstall.Channel = p.Channel
	//dbInstall.PushType = p.PushType
	//dbInstall.PushToken = p.PushToken
	//dbInstall.Area = p.Area
	//dbInstall.Longitude = p.Longitude
	//dbInstall.Latitude = p.Latitude
	dbInstall.ActiveIp = c.ClientIP()

	if dbInstall.Id != 0 {
		if err = db.Save(&dbInstall).Error; err != nil {
			logger.Error(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	} else {
		dbInstall.Id, _ = strconv.ParseUint(time.Now().Format("060102150405"), 10, 64)
		dbInstall.Id = dbInstall.Id*100 + uint64(time.Now().Nanosecond()/1e7)
		dbInstall.Id = dbInstall.Id*100 + uint64(rand.Intn(100))
		if err = db.Create(&dbInstall).Error; err != nil {
			logger.Error(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success, "data": dbInstall})
}
