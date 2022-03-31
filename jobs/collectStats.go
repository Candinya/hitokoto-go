package jobs

import (
	"context"
	jsoniter "github.com/json-iterator/go"
	"hitokoto-go/consts"
	"hitokoto-go/global"
	"hitokoto-go/models"
	"hitokoto-go/types"
	"log"
)

func CollectStats() {

	ctx := context.TODO()
	json := jsoniter.ConfigCompatibleWithStandardLibrary

	statsCount, err := global.Redis.LLen(ctx, consts.VisitListName).Result()
	if err != nil {
		log.Println("[ERROR] Failed to get stats count with error: ", err)
		return
	}

	counts := make(map[string]map[uint]uint)

	for i := int64(0); i < statsCount; i++ {
		// Get bytes (LPush -> RPop)
		rBytes, err := global.Redis.RPop(ctx, consts.VisitListName).Bytes()
		if err != nil {
			log.Println("[ERROR] Failed to pop stats record with error: ", err)
			return
		}

		// Parse record
		var r types.VisitHitokotoRecord
		if err = json.Unmarshal(rBytes, &r); err != nil {
			log.Println("[ERROR] Failed to unmarshal stats record with error: ", err)
			return
		}

		// Initialize count
		if _, ok := counts[r.Category]; !ok {
			counts[r.Category] = make(map[uint]uint)
		}
		if _, ok := counts[r.Category][r.ID]; !ok {
			counts[r.Category][r.ID] = 0
		}
		// Add to count
		counts[r.Category][r.ID]++
	}

	// Update stats into database
	for category, ids := range counts {
		for id, count := range ids {
			// Get old count
			var countInDB uint
			if err := global.DB.
				Scopes(models.SentenceTable(models.Sentence{Type: category})).
				Model(&models.Sentence{}).
				Where("id = ?", id).
				Select("visit_count").
				First(&countInDB).Error; err != nil {
				log.Println("[ERROR] Failed to get old visit count with error: ", err)
			}

			// Calc new count
			countInDB += count

			// Save new count

			if err := global.DB.
				Scopes(models.SentenceTable(models.Sentence{Type: category})).
				Model(&models.Sentence{}).
				Where("id = ?", id).
				Update("visit_count", countInDB).Error; err != nil {
				log.Println("[ERROR] Failed to save new visit count with error: ", err)
			}
		}
	}

}
