package telegram

import (
	"bytes"
	"encoding/json"
	"github.com/cherryReptile/Todo/internal/database"
	"io"
	"log"
	"os"
)

type Service struct {
	DB       *database.SqlLite
	BotToken string
	BotUrl   string
}

type Updates struct {
	OK     bool             `json:"ok"`
	Result []MessageWrapper `json:"result"`
}

type BotMessage struct {
	Ok     bool    `json:"ok"`
	Result Message `json:"result"`
}

type Entity struct {
	Offset int    `json:"offset"`
	Length int    `json:"length"`
	Type   string `json:"type"`
}

type Chat struct {
	Id        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	Type      string `json:"type"`
}

type From struct {
	Id           uint   `json:"id"`
	IsBot        bool   `json:"is_bot"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Username     string `json:"username"`
	LanguageCode string `json:"language_code"`
}

type Message struct {
	MessageId int      `json:"message_id"`
	From      From     `json:"from"`
	Chat      Chat     `json:"chat"`
	Date      int      `json:"date"`
	Text      string   `json:"text"`
	Entities  []Entity `json:"entities"`
}

type MessageWrapper struct {
	UpdateId int     `json:"update_id"`
	Message  Message `json:"message"`
}

type ToMessage struct {
	ChatId uint   `json:"chat_id"`
	Text   string `json:"text"`
}

func (s *Service) Init(db *database.SqlLite) {
	s.DB = db
	s.BotToken = os.Getenv("TG_TOKEN")
	s.BotUrl = "https://api.telegram.org/bot" + s.BotToken
}

func (t *ToMessage) ToReader() io.Reader {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(t)
	if err != nil {
		log.Fatal(err)
	}

	return &buf
}