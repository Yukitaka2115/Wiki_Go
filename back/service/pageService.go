package service

import (
	"encoding/json"
	"errors"
	"log"
	"wiki/config"
	"wiki/dao"
	"wiki/middleware"
	"wiki/model"
)

func AddPage(page model.Page) {
	err := dao.Db.AutoMigrate(&model.Page{})
	if err != nil {
		return
	}

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

func UpdatePage(id int, newPage model.Page) (model.Page, error) {
	var page model.Page
	if err := dao.Db.First(&page, id).Error; err != nil {
		return model.Page{}, err
	}
	if page.ID == 0 {
		return model.Page{}, errors.New("页面不存在")
	}
	page = newPage
	if err := dao.Db.Model(&page).Save(&page).Error; err != nil {
		return model.Page{}, err
	}
	return page, nil
}

func GetPagesWithPagination(page, size int) ([]model.Page, error) {
	var pages []model.Page

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
	var page model.Page
	dao.Db.First(&page, id)
	dao.Db.Delete(&page)
}

func GetPageByTitle(title string) model.Page {
	var page model.Page
	// 查询数据库
	result := dao.Db.Where("title = ?", title).First(&page)
	if result.Error != nil {
		log.Printf("未找到标题为 %s 的词条", title)
		return page
	}

	// 只有在查到数据且 ID 合法时才增加点击量
	if page.ID > 0 {
		go func(pid int) {
			ranking := middleware.NewPageRanking()
			// 增加判断：如果 ranking.client 是 nil，直接返回，不调方法
			if ranking.Client == nil {
				log.Println("Redis 客户端未初始化，跳过点击量增加")
				return
			}
			_ = ranking.IncreasePageVisit(pid)
		}(page.ID)
	}

	// 3. 原有的 JSON 解析逻辑
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
