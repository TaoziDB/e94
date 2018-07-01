package main

import (
	"background/newmovie/config"
	"background/common/logger"
	"background/newmovie/model"
	"background/common/util"

	"strings"
	"background/common/constant"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"

	"flag"
	"os"
	"bufio"
	"io"
)

func main(){
	logger.SetLevel(config.GetLoggerLevel())

	configPath := flag.String("conf", "../config/config.json", "Config file path")
	flag.Parse()

	err := config.LoadConfig(*configPath)
	if err != nil {
		return
	}

	db, err := gorm.Open(config.GetDBName(), config.GetDBSource())
	if err != nil {
		logger.Fatal("Open db Failed!!!!", err)
		return
	}
	db.LogMode(true)
	model.InitModel(db)
	f, err := os.Open("/home/lyric/Git/e94/src/background/newmovie/tools/result.txt")
	if err != nil {
		logger.Error(err)
	}
	defer f.Close()

	rd := bufio.NewReader(f)
	for {
		line, err := rd.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		}
		line = strings.Replace(line, "\n", "", -1)

		logger.Debug(line)

		var title string
		var url string
		fields := strings.Split(line, "|")
		if len(fields) == 2 {
			title = fields[0]
			url = fields[1]
		} else {
			continue
		}

		var stream model.Stream
		stream.Title = title
		stream.Title = strings.Replace(stream.Title,"高清","",-1)
		stream.Title = strings.Replace(stream.Title,"-","",-1)
		//stream.Title = util.TrimChinese(stream.Title)
		stream.Pinyin = util.TitleToPinyin(stream.Title)
		stream.Title = strings.Trim(stream.Title," ")
		logger.Debug(stream.Title)

		tx := db.Begin()
		if err := tx.Where("title = ?",stream.Title).First(&stream).Error ; err == gorm.ErrRecordNotFound{
			if strings.Contains(stream.Title,"CCTV"){
				stream.Category = "央视"
			}else if strings.Contains(stream.Title,"卫视"){
				stream.Category = "卫视"
			}else{
				stream.Category = "地方"
			}

			stream.OnLine = constant.MediaStatusOnLine
			stream.Sort = 0

			if err = tx.Create(&stream).Error ; err != nil{
				tx.Rollback()
				logger.Error(err)
				return
			}
		}

		var play model.PlayUrl
		play.Url = url
		play.Provider = uint32(constant.ContentProviderSystem)
		if err := tx.Where("provider = ? and url = ?",play.Provider,play.Url).First(&play).Error ; err == gorm.ErrRecordNotFound{
			play.Title = title
			play.OnLine = constant.MediaStatusOnLine
			play.Sort = 0
			play.ContentType = uint8(constant.MediaTypeStream)
			play.ContentId = stream.Id
			play.Quality = uint8(constant.VideoQuality720p)

			if err = tx.Create(&play).Error ; err != nil{
				tx.Rollback()
				logger.Error(err)
				return
			}
		}
		tx.Commit()
	}


}

