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

//Updates

type Updates struct {
	OK     bool             `json:"ok"`
	Result []MessageWrapper `json:"result"`
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
	MessageId   int         `json:"message_id"`
	From        From        `json:"from"`
	Chat        Chat        `json:"chat"`
	Date        int         `json:"date"`
	Text        string      `json:"text"`
	Entities    []Entity    `json:"entities"`
	ReplyMarkup ReplyMarkup `json:"reply_markup"`
}

type MessageWrapper struct {
	UpdateId      int           `json:"update_id"`
	CallbackQuery CallbackQuery `json:"callback_query"`
	Message       Message       `json:"message"`
}

//Callbacks

type CallbackQuery struct {
	Id           string `json:"id"`
	From         `json:"from"`
	Message      `json:"message"`
	ChatInstance string `json:"chat_instance"`
	Data         string `json:"data"`
}

type ReplyMarkup struct {
	InlineKeyboard [][1]InlineKeyboardBtn `json:"inline_keyboard"`
}

type InlineKeyboardBtn struct {
	Text         string `json:"text"`
	CallbackData string `json:"callback_data"`
}

type Callback struct {
	UpdateId      int `json:"update_id"`
	CallbackQuery `json:"callback_query"`
}

//BotMsgs

type BotMessage struct {
	Ok     bool    `json:"ok"`
	Result Message `json:"result"`
}

// Converting structs

type ToMessage struct {
	ChatId uint   `json:"chat_id"`
	Text   string `json:"text"`
}
type ToEditMessage struct {
	ChatId    uint   `json:"chat_id"`
	MessageId int    `json:"message_id"`
	Text      string `json:"text"`
}
type ToUpdates struct {
	//For specify update request
	//Offset         int      `json:"offset"`
	//Limit          int      `json:"limit"`
	//Timeout        int      `json:"timeout"`
	AllowedUpdates []string `json:"allowed_updates"`
}

type ToAnswerCallback struct {
	CallbackQueryId string `json:"callback_query_id"`
	Text            string `json:"text"`
}

type ToInlineKeyboardBtn struct {
	ChatId      uint        `json:"chat_id"`
	MessageId   int         `json:"message_id"`
	Text        string      `json:"text"`
	ReplyMarkup ReplyMarkup `json:"reply_markup"`
}

//Just id and type model for callback

type ModelFromCallback struct {
	Id     uint   `json:"id"`
	Model  string `json:"model"`
	Method string `json:"method"`
}

//For service struct

func (s *Service) Init(db *database.SqlLite) {
	s.DB = db
	s.BotToken = os.Getenv("TG_TOKEN")
	s.BotUrl = "https://api.telegram.org/bot" + s.BotToken
}

func ToReader(i interface{}) io.Reader {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(i)
	if err != nil {
		log.Fatal(err)
	}

	return &buf
}
