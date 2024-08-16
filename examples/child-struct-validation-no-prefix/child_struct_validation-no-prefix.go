package main

import (
	"fmt"
	"github.com/ahmadrezamusthafa/deep-validator"
	"time"
)

type TransactionUpdatedEventPayload struct {
	TransactionId                string                 `json:"transactionId"`
	TransactionType              string                 `json:"transactionType"`
	TotalAmount                  float64                `json:"totalAmount"`
	SenderType                   string                 `json:"senderType"`
	SenderProviderId             *string                `json:"senderProviderId"`
	SenderProviderAccountId      string                 `json:"senderProviderAccountId"`
	TransactionReference         string                 `json:"transactionReference"`
	Status                       string                 `json:"status"`
	CreatedAt                    *time.Time             `json:"createdAt"`
	ConfirmedAt                  *time.Time             `json:"confirmedAt"`
	CompletedAt                  *time.Time             `json:"completedAt"`
	ExtraInfo                    map[string]interface{} `json:"extraInfo"`
	Currency                     string                 `json:"currency"`
	BeneficiaryType              string                 `json:"beneficiaryType"`
	BeneficiaryProviderId        string                 `json:"beneficiaryProviderId"`
	BeneficiaryProviderAccountId string                 `json:"beneficiaryProviderAccountId"`
}

type TransactionUpdatedEvent struct {
	EventSource          string                         `json:"eventSource"`
	EventId              string                         `json:"eventId"`
	EventPublishedAt     time.Time                      `json:"eventPublishedAt"`
	MessageFormatVersion string                         `json:"messageFormatVersion"`
	Payload              TransactionUpdatedEventPayload `json:"payload"`
}

func strToStrPtr(s string) *string {
	return &s
}

func main() {
	query := `(SenderProviderId=bca && TotalAmount=50000) && 
				(SenderProviderAccountId=121 || SenderProviderAccountId=1212121212)`
	val := "2022-11-03T16:50:16+07:00"
	createdAt, _ := time.Parse(time.RFC3339, val)
	transaction := TransactionUpdatedEvent{
		EventSource:          "",
		EventId:              "",
		EventPublishedAt:     time.Time{},
		MessageFormatVersion: "",
		Payload: TransactionUpdatedEventPayload{
			TransactionId:                "",
			TransactionType:              "",
			TotalAmount:                  50000,
			SenderType:                   "",
			SenderProviderId:             strToStrPtr("bca"),
			SenderProviderAccountId:      "1212121212",
			TransactionReference:         "",
			Status:                       "",
			CreatedAt:                    &createdAt,
			ConfirmedAt:                  nil,
			CompletedAt:                  nil,
			ExtraInfo:                    nil,
			Currency:                     "",
			BeneficiaryType:              "",
			BeneficiaryProviderId:        "",
			BeneficiaryProviderAccountId: "",
		},
	}

	isValid, err := deepvalidator.NewProcessor().
		RegisterCondition(query).
		SetRemovePrefix(true).
		ValidateMultipleStructs(transaction.Payload)

	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Is valid:", isValid)
	}
}
