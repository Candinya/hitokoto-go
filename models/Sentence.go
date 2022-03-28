package models

import (
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
	"hitokoto-go/types"
	"strconv"
	"time"
)

type Sentence struct {
	gorm.Model

	UUID       uuid.UUID `gorm:"type:uuid;index;not null"`
	Hitokoto   string
	Type       string `gorm:"-"`
	From       string
	FromWho    string
	Creator    string
	CreatorUID uint
	Reviewer   uint
	CommitFrom string
	Length     uint
}

const prefix = "sentence_"

func SentenceTable(s Sentence) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		return tx.Table(prefix + s.Type)
	}
}

// How to use:
// db.Scopes(SentenceTable(Sentence{Type: "x"}))
//   .Find(&sentences)

func (s *Sentence) ToJSON(t string) *types.Sentence {
	return &types.Sentence{
		ID:         s.ID,
		UUID:       s.UUID.String(),
		Hitokoto:   s.Hitokoto,
		Type:       t,
		From:       s.From,
		FromWho:    s.FromWho,
		Creator:    s.Creator,
		CreatorUID: s.CreatorUID,
		Reviewer:   s.Reviewer,
		CommitFrom: s.CommitFrom,
		CreatedAt:  strconv.FormatInt(s.CreatedAt.Unix(), 10),
		Length:     s.Length,
	}
}

func (s *Sentence) FromJSON(sentence *types.Sentence) {
	// Set create time
	created, err := strconv.ParseInt(sentence.CreatedAt, 10, 64)
	if err == nil {
		s.CreatedAt = time.Unix(created, 0)
	}

	s.UUID = uuid.FromStringOrNil(sentence.UUID)
	s.Hitokoto = sentence.Hitokoto
	s.Type = sentence.Type
	s.From = sentence.From
	s.FromWho = sentence.FromWho
	s.Creator = sentence.Creator
	s.CreatorUID = sentence.CreatorUID
	s.Reviewer = sentence.Reviewer
	s.CommitFrom = sentence.CommitFrom
	s.Length = sentence.Length
}
