package models

import (
	"fmt"
	"time"
)

type User struct {
	ID             string `json:"id" db:"id"`
	Username       string `json:"username" db:"username"`
	HashedPassword string `json:"-" db:"password"`
	IsAdmin        bool   `json:"-" db:"is_admin"`
}

type Word struct {
	ID            string   `db:"id" json:"id"`
	EN            string   `db:"en" json:"en"`
	RU            []string `db:"ru" json:"ru"`
	Transcription string   `db:"transcription" json:"transcription"`
	Examples      []string `db:"examples" json:"examples"`
}

func (w *Word) SetExamples(enRuExamples ...string) error {
	if len(enRuExamples)%2 != 0 {
		return fmt.Errorf("examples should contain en and ru pairs")
	}

	w.Examples = enRuExamples
	return nil
}

type UserWord struct {
	WordID      string    `db:"word_id"`
	Username    string    `db:"username"`
	Rating      float32   `db:"knowledge_rating"`
	Attempts    int       `db:"consecutive_success_attempts"`
	LastAttempt time.Time `db:"last_attempt"`
}
