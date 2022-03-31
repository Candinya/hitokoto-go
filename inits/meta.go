package inits

import (
	"hitokoto-go/global"
	"hitokoto-go/models"
	"hitokoto-go/types"
	"log"
)

func Meta() error {
	var categories []models.Category
	if err := global.DB.Find(&categories).Error; err != nil {
		return err
	}

	global.Meta.AllCount = 0

	for _, c := range categories {
		var count int64 // Prepare counts metadata for later random usages
		if err := global.DB.Scopes(models.SentenceTable(models.Sentence{Type: c.Key})).Count(&count).Error; err != nil {
			// warn and skip
			log.Println("[WARN] Failed to initialize category ", c.Key, " with error: ", err)
			continue
		}
		uCount := uint(count)
		global.Meta.Categories = append(global.Meta.Categories, types.MetaCategory{
			Key:    c.Key,
			Counts: uCount,
		})
		global.Meta.AllCount += uCount
	}

	// No error
	return nil

}
