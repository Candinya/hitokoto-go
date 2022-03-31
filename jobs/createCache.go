package jobs

import (
	"hitokoto-go/global"
	"hitokoto-go/utils"
)

func RefreshCache() {
	for _, c := range global.Meta.Categories {
		_ = utils.CacheCategory(c.Key)
	}
}
