package utils

import (
	"hitokoto-go/global"
	"hitokoto-go/models"
	"hitokoto-go/types"
)

func SelectHitokotoFromDB(limit *SelectLimitation) *types.Sentence {
	// Select hitokoto
	var ht models.Sentence
	global.DB.
		Scopes(
			models.SentenceTable(
				models.Sentence{Type: limit.Category},
			),
		).
		Where("length >= ?", limit.MinLen).
		Where("length <= ?", limit.MaxLen).
		Order("RANDOM()").
		First(&ht)
	return ht.ToJSON(limit.Category)
}
