package service

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"wiki/config"
	"wiki/dao"
)

type Page struct {
	ID           int        `gorm:"primary_key"`
	Title        string     `json:"Title"`
	Brief        string     `json:"Brief"`
	Background   string     `json:"Background"`
	History      mapHandler `json:"History" gorm:"type:json"`
	Mainchara    Chara
	MaincharaStr string
	Relatives    Chara
	RelativesStr string
	Team         Group
	TeamStr      string
	Comments     Comment
}

type Group struct {
	Chara string
	Grade string
	Group string
}

type Chara struct {
	Chara string
	Cast  string
	Info  string
}

type Comment struct {
	Username string
	Comment  string
}

type mapHandler map[string]interface{}

func (mh mapHandler) Value() (driver.Value, error) {
	str, err := json.Marshal(mh)
	if err != nil {
		return nil, err
	}
	return string(str), nil
}

func (mh *mapHandler) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	b, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("value is not []byte, value: %v", value)
	}

	return json.Unmarshal(b, &mh)
}

func AddPage(page Page) {
	dao.Db.AutoMigrate(&Page{})

	maincharaStr, err := config.ConvertStructToJSONStr(page.Mainchara)
	page.MaincharaStr = maincharaStr
	relativesStr, err := config.ConvertStructToJSONStr(page.Relatives)
	page.RelativesStr = relativesStr
	teamStr, err := config.ConvertStructToJSONStr(page.Team)
	page.TeamStr = teamStr

	// 尝试创建记录
	if err := dao.Db.Create(&page).Error; err != nil {
		log.Println("Failed to create page:", err)
		return
	}

	// 打印存储的数据
	jsonData, err := json.Marshal(page)
	if err != nil {
		log.Println("Failed to marshal page to JSON:", err)
		return
	}
	log.Println("Page stored in database:", string(jsonData))
}

func UpdatePage(id int, newPage Page) (Page, error) {
	var page Page
	if err := dao.Db.First(&page, id).Error; err != nil {
		return Page{}, err
	}
	if page.ID == 0 {
		return Page{}, errors.New("页面不存在")
	}
	page = newPage
	if err := dao.Db.Model(&page).Save(&page).Error; err != nil {
		return Page{}, err
	}
	return page, nil
}

func GetPagesWithPagination(page, size int) ([]Page, error) {
	var pages []Page

	// 计算偏移量
	offset := (page - 1) * size
	res := dao.Db.Offset(offset).Limit(size).Find(&pages)
	if res.Error != nil {
		log.Println("Failed to find pages:", res.Error)
		return nil, res.Error
	}

	for i := range pages {
		err := config.UnmarshalJSONStr(pages[i].MaincharaStr, &pages[i].Mainchara)
		if err != nil {
			log.Println("Failed to unmarshal Mainchara:", err)
		}
		err = config.UnmarshalJSONStr(pages[i].RelativesStr, &pages[i].Relatives)
		if err != nil {
			log.Println("Failed to unmarshal Relatives:", err)
		}
		err = config.UnmarshalJSONStr(pages[i].TeamStr, &pages[i].Team)
		if err != nil {
			log.Println("Failed to unmarshal Team:", err)
		}
	}

	return pages, nil
}

func DeletePageByID(id int) {
	var page Page
	dao.Db.First(&page, id)
	dao.Db.Delete(&page)
}

func GetPageByTitle(title string) Page {
	var page Page
	dao.Db.Where("title = ?", title).First(&page)

	if err := config.UnmarshalJSONStr(page.MaincharaStr, &page.Mainchara); err != nil {
		log.Println("Failed to unmarshal Mainchara:", err)
	}
	if err := config.UnmarshalJSONStr(page.RelativesStr, &page.Relatives); err != nil {
		log.Println("Failed to unmarshal Relatives:", err)
	}
	if err := config.UnmarshalJSONStr(page.TeamStr, &page.Team); err != nil {
		log.Println("Failed to unmarshal Team:", err)
	}

	return page
}
