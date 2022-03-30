package utils

import (
	"context"
	"encoding/json"
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

	err := saveHitokotosIntoCache(c, ss)

	return len(ss), err
}

func saveHitokotosIntoCache(c string, ss []*types.Sentence) error {

	//ctx := context.TODO() // Temporarily remove context time limit
	ctx, _ := context.WithTimeout(context.TODO(), 3*time.Second)

	// Clear old cache
	global.Redis.LTrim(ctx, encodeCacheKey(c), 1, 0)

	var cacheSentenceBytes []interface{} // ???
	for _, s := range ss {
		sentenceBytes, err := json.Marshal(&s)
		if err != nil {
			return err
		}
		cacheSentenceBytes = append(cacheSentenceBytes, sentenceBytes)
	}
	if len(cacheSentenceBytes) == 0 {
		// Nothing to cache
		log.Println("[WARN] Nothing to cache!")
		return nil
	}
	if n, err := global.Redis.LPush(ctx, encodeCacheKey(c), cacheSentenceBytes...).Result(); err != nil {
		log.Println("[ERROR] Failed to save hitokoto into cache with error: ", err)
		return err
	} else {
		log.Println("[Info] Successfully saved ", n, " hitokotos of type ", c, " into cache")
	}

	return nil
}

type SelectLimitation struct {
	Category string
	MinLen   uint
	MaxLen   uint
}

func SelectHitokotoFromCache(limit *SelectLimitation) *types.Sentence {

	ctx, _ := context.WithTimeout(context.TODO(), 3*time.Second)

	// Check cache
	sCachedCount, err := global.Redis.LLen(ctx, encodeCacheKey(limit.Category)).Result()
	if sCachedCount == 0 || err != nil {
		log.Println("[WARN] No hitokoto found in cache for category ", limit.Category)
		return nil
	}

	isMaxLengthLimitInvalid := limit.MaxLen < limit.MinLen

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

		// Push back to the end of the list
		if err = global.Redis.RPush(ctx, encodeCacheKey(limit.Category), cachedBytes).Err(); err != nil {
			log.Println("[ERROR] Fail to push back cached category ", limit.Category, " with error: ", err)
		}

		// Check if the sentence is valid
		if s.Length >= limit.MinLen && // Check min length
			(isMaxLengthLimitInvalid || // If max length limitation is valid
				s.Length <= limit.MaxLen) { // Check max length
			// The one
			// Mark it as selected
			selectedSentence = s

			// Found!
			break
		}

	}
	return selectedSentence
}
