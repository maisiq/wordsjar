package repository

import (
	"strings"
	"time"
)

type Word struct {
	ID            string   `db:"id"`
	EN            string   `db:"en"`
	RU            string   `db:"ru"`
	Transcription string   `db:"transcription"`
	Examples      []string `db:"examples"`
}

func (w Word) RuTranslations() []string {
	var parsed []string
	trimmed := strings.Trim(w.RU, "{}")
	if trimmed != "" {
		splitted := strings.Split(trimmed, ",")
		for _, v := range splitted {
			parsed = append(parsed, strings.Trim(v, `"`))
		}
	}
	return parsed
}

type UserWord struct {
	WordID      string    `db:"word_id"`
	Username    string    `db:"username"`
	Rating      float32   `db:"knowledge_rating"`
	Attempts    int       `db:"consecutive_success_attempts"`
	LastAttempt time.Time `db:"last_attempt"`
}
