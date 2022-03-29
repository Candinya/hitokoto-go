package inits

import (
	"hitokoto-go/global"
	"hitokoto-go/models"
	"hitokoto-go/types"
	"hitokoto-go/utils"
	"log"
)

func Meta() error {
	var categories []models.Category
	if err := global.DB.Find(&categories).Error; err != nil {
		return err
	}

	global.Meta.AllCount = 0

	for _, c := range categories {
		// Prepare counts metadata for later random usages
		if count, err := utils.CacheCategory(c.Key); err != nil {
			// warn and skip
			log.Println("[WARN] Failed to initialize category ", c.Key, " with error: ", err)
			continue
		} else {
			global.Meta.AllCount += uint(count)
			global.Meta.Categories = append(global.Meta.Categories, types.MetaCategory{
				Key:    c.Key,
				Counts: uint(count),
			})
		}
	}

	// No error
	return nil

}
