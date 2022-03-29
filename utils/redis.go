package utils

import (
	"context"
	"encoding/json"
	"hitokoto-go/consts"
	"hitokoto-go/global"
	"hitokoto-go/models"
	"hitokoto-go/types"
	"log"
)

func encodeCacheKey(c string) string {
	return "hitokoto:cache:sentence_" + c
}

// CacheCategory : Pre-random sentences to prevent the database from being overloaded
func CacheCategory(c string) (int, error) {
	if !global.Config.IsProdMode {
		log.Println("[Info] Create cache for category ", c)
	}
	// Get random table from database
	var mss []models.Sentence // Model sentences
	global.DB.
		Scopes(
			models.SentenceTable(
				models.Sentence{Type: c},
			),
		).
		Order("RANDOM()").
		// Limit(consts.RandTableSize). // Remove limit so that we can cache all sentences
		Find(&mss)

	// Transform to type
	var ss []*types.Sentence // Sentences
	for _, ms := range mss {
		ss = append(ss, ms.ToJSON(c))
	}

	err := SaveHitokotosIntoCache(c, ss)

	return len(ss), err
}

func SaveHitokotosIntoCache(c string, ss []*types.Sentence) error {

	ctx := context.TODO() // Temporarily remove context time limit

	var cacheByteSs []interface{} // ???
	for _, s := range ss {
		cacheBytes, err := json.Marshal(&s)
		if err != nil {
			return err
		}
		cacheByteSs = append(cacheByteSs, cacheBytes)
	}
	if err := global.Redis.LPush(ctx, encodeCacheKey(c), cacheByteSs...).Err(); err != nil {
		log.Println("[ERROR] Failed to save hitokoto into cache with error: ", err)
		return err
	}

	return nil
}

type SelectLimitation struct {
	Category string
	MinLen   uint
	MaxLen   uint
}

func SelectHitokotoFromCache(limit *SelectLimitation) *types.Sentence {

	ctx := context.TODO() // Temporarily remove context time limit

	// Check cache
	sCachedCount, err := global.Redis.LLen(ctx, encodeCacheKey(limit.Category)).Result()
	if sCachedCount == 0 || err != nil {
		// Cache is invalid, rebuild it
		if !global.Config.IsProdMode {
			log.Println("[Info] Cache is invalid, rebuild it")
		}
		nscc, err := CacheCategory(limit.Category)
		if err != nil {
			log.Println("[ERROR] Fail to cache category ", limit.Category, " with error: ", err)
			return nil
		}
		sCachedCount = int64(nscc)
	}

	var selectedSentence *types.Sentence
	for i := int64(0); i < sCachedCount; i++ {
		var s *types.Sentence
		cachedBytes, err := global.Redis.LPop(ctx, encodeCacheKey(limit.Category)).Bytes()
		if err != nil {
			log.Println("[ERROR] Fail to get cached category ", limit.Category, " index ", i, " with error: ", err)
			return nil
		}
		if err = json.Unmarshal(cachedBytes, &s); err != nil {
			log.Println("[ERROR] Fail to parse cached category ", limit.Category, " index ", i, " with error: ", err)
			return nil
		}

		// Check if the sentence is valid
		if s.Length >= limit.MinLen && s.Length <= limit.MaxLen {
			// The one
			// Mark it as selected
			selectedSentence = s

			// Found!
			break
		} else {
			// Push back to the end of the list
			if err = global.Redis.RPush(ctx, encodeCacheKey(limit.Category), cachedBytes).Err(); err != nil {
				log.Println("[ERROR] Fail to push back cached category ", limit.Category, " with error: ", err)
			}
		}

	}

	if sCachedCount < consts.RandTableRecreateThreshold {
		// Recreate cache
		// global.Redis.Del(ctx, encodeCacheKey(limit.Category))
		//if !global.Config.IsProdMode {
		//	log.Println("[Info] Recreate threshold cache for category ", limit.Category)
		//}
		go CacheCategory(limit.Category) // Just ignore error
	}
	return selectedSentence
}
