package models

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID            uuid.UUID  `json:"id" db:"id"`
	UserID        *uuid.UUID `json:"user_id,omitempty" db:"user_id"`
	Content       string     `json:"content" db:"content"`
	SenderHeader  string     `json:"sender_header" db:"sender_header"`
	ReceivedAt    *time.Time `json:"received_at,omitempty" db:"received_at"`
	MessageType   string     `json:"message_type" db:"message_type"` // SMS, WhatsApp, Email
	PhoneNumber   *string    `json:"phone_number,omitempty" db:"phone_number"`
	HasLinks      bool       `json:"has_links" db:"has_links"`
	LinkCount     int        `json:"link_count" db:"link_count"`
	ExtractedURLs []string   `json:"extracted_urls,omitempty" db:"extracted_urls"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`
}

type MessageInput struct {
	Content      string     `json:"content" binding:"required,max=1000"`
	SenderHeader string     `json:"sender_header" binding:"required"`
	ReceivedAt   *time.Time `json:"received_at,omitempty"`
	MessageType  string     `json:"message_type" binding:"omitempty,oneof=SMS WhatsApp Email"`
	PhoneNumber  *string    `json:"phone_number,omitempty"`
}

type MessageFeatures struct {
	Content          string   `json:"content"`
	SenderHeader     string   `json:"sender_header"`
	MessageLength    int      `json:"message_length"`
	HasLinks         bool     `json:"has_links"`
	LinkCount        int      `json:"link_count"`
	ExtractedURLs    []string `json:"extracted_urls"`
	HasPhoneNumber   bool     `json:"has_phone_number"`
	PhoneNumberCount int      `json:"phone_number_count"`
	HasUrgentWords   bool     `json:"has_urgent_words"`
	UrgentWordCount  int      `json:"urgent_word_count"`
	SpecialCharRatio float64  `json:"special_char_ratio"`
	CapitalRatio     float64  `json:"capital_ratio"`
	NumberRatio      float64  `json:"number_ratio"`
	HasKYCKeywords   bool     `json:"has_kyc_keywords"`
	HasBankNames     bool     `json:"has_bank_names"`
}

