package utils

import (
	"context"
	"encoding/json"
	"hitokoto-go/consts"
	"hitokoto-go/global"
	"hitokoto-go/models"
	"hitokoto-go/types"
	"log"
	"time"
)

func encodeCacheKey(c string) string {
	return "hitokoto:cache:sentence_" + c
}

// CacheCategory : Pre-random sentences to prevent the database from being overloaded
func CacheCategory(c string) error {
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
		Limit(consts.RandTableSize).
		Find(&mss)

	// Transform to type
	var ss []*types.Sentence // Sentences
	for _, ms := range mss {
		ss = append(ss, ms.ToJSON(c))
	}

	err := SaveHitokotosIntoCache(c, ss)

	return err
}

func SaveHitokotosIntoCache(c string, ss []*types.Sentence) error {
	cacheBytes, err := json.Marshal(&ss)
	if err != nil {
		return err
	}

	ctx, _ := context.WithTimeout(context.TODO(), 1*time.Second)
	err = global.Redis.Set(ctx, encodeCacheKey(c), string(cacheBytes), consts.RandTableCacheExpire).Err()

	if err != nil {
		log.Println("[ERROR] Failed to save hitokotos into cache with error: ", err)
		return err
	}

	return err
}

type SelectLimitation struct {
	Category string
	MinLen   uint
	MaxLen   uint
}

func SelectHitokotoFromCache(limit *SelectLimitation) *types.Sentence {
	var ss []*types.Sentence

	ctx, _ := context.WithTimeout(context.TODO(), 1*time.Second)

	// Check cache
	existCount, err := global.Redis.Exists(ctx, encodeCacheKey(limit.Category)).Result()
	if existCount == 0 || err != nil {
		// Cache is invalid, rebuild it
		if !global.Config.IsProdMode {
			log.Println("[Info] Cache is invalid, rebuild it")
		}
		err = CacheCategory(limit.Category)
		if err != nil {
			log.Println("[ERROR] Fail to cache category ", limit.Category, " with error: ", err)
			return nil
		}
	}

	cachedBytes, err := global.Redis.Get(ctx, encodeCacheKey(limit.Category)).Bytes()
	if err != nil {
		log.Println("[ERROR] Fail to get cached category ", limit.Category, " with error: ", err)
		return nil
	}
	if err = json.Unmarshal(cachedBytes, &ss); err != nil {
		log.Println("[ERROR] Fail to parse cached category ", limit.Category, " with error: ", err)
		return nil
	}

	var selectedSentence *types.Sentence
	var nss []*types.Sentence // New Sentences
	for i, s := range ss {
		if s.Length >= limit.MinLen && s.Length <= limit.MaxLen {
			// The one
			// Split from cache
			nss = append(ss[:i], ss[i+1:]...)
			// Mark it as selected
			selectedSentence = s

			// Found!
			break
		}
	}

	if len(nss) < consts.RandTableRecreateThreshold {
		// Recreate cache
		// global.Redis.Del(ctx, encodeCacheKey(limit.Category))
		//if !global.Config.IsProdMode {
		//	log.Println("[Info] Recreate threshold cache for category ", limit.Category)
		//}
		go CacheCategory(limit.Category) // Just ignore error
	} else if selectedSentence != nil {
		// Update cache
		//if !global.Config.IsProdMode {
		//	log.Println("[Info] Update cache with unused sentences for category ", limit.Category)
		//}
		_ = SaveHitokotosIntoCache(limit.Category, nss) // Ignore error
	}
	return selectedSentence
}
