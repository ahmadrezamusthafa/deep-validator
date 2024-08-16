package main

import (
	"fmt"
	"github.com/ahmadrezamusthafa/deep-validator"
	"time"
)

type PartnerStatementUpdatedEventPayload struct {
	EventType                 string                 `json:"eventType"`
	PartnerStatementId        string                 `json:"partnerStatementId"`
	PartnerStatementType      string                 `json:"partnerStatementType"`
	PartnerId                 string                 `json:"partnerId"`
	PartnerAccountId          string                 `json:"partnerAccountId"`
	TransactionReference      string                 `json:"transactionReference"`
	Description               string                 `json:"description"`
	Currency                  string                 `json:"currency"`
	HashCode                  string                 `json:"hashCode"`
	Status                    string                 `json:"status"`
	TransactionAt             time.Time              `json:"transactionAt"`
	CreatedAt                 time.Time              `json:"createdAt"`
	Version                   int64                  `json:"version"`
	AssignmentTransactionId   *string                `json:"assignmentTransactionId"`
	AssignmentTransactionType *string                `json:"assignmentTransactionType"`
	ExtraInfo                 map[string]interface{} `json:"extraInfo"`
	Credit                    float64                `json:"credit"`
	Debit                     float64                `json:"debit"`
}

type PartnerStatementUpdatedEvent struct {
	EventSource          string                              `json:"eventSource"`
	EventId              string                              `json:"eventId"`
	EventPublishedAt     time.Time                           `json:"eventPublishedAt"`
	MessageFormatVersion string                              `json:"messageFormatVersion"`
	Payload              PartnerStatementUpdatedEventPayload `json:"payload"`
	RetryAt              *time.Time                          `json:"retryAt"`
}

type TransactionUpdatedEventPayload struct {
	TransactionId                string                 `json:"transactionId"`
	TransactionType              string                 `json:"transactionType"`
	TotalAmount                  float64                `json:"totalAmount"`
	SenderType                   *string                `json:"senderType"`
	SenderProviderId             *string                `json:"senderProviderId"`
	SenderProviderAccountId      *string                `json:"senderProviderAccountId"`
	TransactionReference         string                 `json:"transactionReference"`
	Status                       string                 `json:"status"`
	CreatedAt                    *time.Time             `json:"createdAt"`
	ConfirmedAt                  *time.Time             `json:"confirmedAt"`
	CompletedAt                  *time.Time             `json:"completedAt"`
	ExtraInfo                    map[string]interface{} `json:"extraInfo"`
	Currency                     string                 `json:"currency"`
	BeneficiaryType              *string                `json:"beneficiaryType"`
	BeneficiaryProviderId        *string                `json:"beneficiaryProviderId"`
	BeneficiaryProviderAccountId *string                `json:"beneficiaryProviderAccountId"`
	Version                      *int64                 `json:"version"`
}

type TransactionUpdatedEvent struct {
	EventSource          string                         `json:"eventSource"`
	EventId              string                         `json:"eventId"`
	EventPublishedAt     time.Time                      `json:"eventPublishedAt"`
	MessageFormatVersion string                         `json:"messageFormatVersion"`
	Payload              TransactionUpdatedEventPayload `json:"payload"`
}

func main() {
	query := `(PartnerId=bca && Debit=50000 && TotalAmount=50000)`
	data := []interface{}{
		PartnerStatementUpdatedEvent{
			EventSource:          "",
			EventId:              "123",
			EventPublishedAt:     time.Time{},
			MessageFormatVersion: "",
			Payload: PartnerStatementUpdatedEventPayload{
				EventType:                 "",
				PartnerStatementId:        "",
				PartnerStatementType:      "",
				PartnerId:                 "bca",
				PartnerAccountId:          "",
				TransactionReference:      "",
				Description:               "",
				Currency:                  "",
				HashCode:                  "",
				Status:                    "",
				TransactionAt:             time.Time{},
				CreatedAt:                 time.Time{},
				Version:                   0,
				AssignmentTransactionId:   nil,
				AssignmentTransactionType: nil,
				ExtraInfo:                 nil,
				Credit:                    0,
				Debit:                     50000,
			},
			RetryAt: &time.Time{},
		},
		TransactionUpdatedEvent{
			EventSource:          "",
			EventId:              "",
			EventPublishedAt:     time.Time{},
			MessageFormatVersion: "",
			Payload: TransactionUpdatedEventPayload{
				TransactionId:                "",
				TransactionType:              "",
				TotalAmount:                  50000,
				SenderType:                   nil,
				SenderProviderId:             nil,
				SenderProviderAccountId:      nil,
				TransactionReference:         "",
				Status:                       "",
				CreatedAt:                    &time.Time{},
				ConfirmedAt:                  &time.Time{},
				CompletedAt:                  &time.Time{},
				ExtraInfo:                    nil,
				Currency:                     "",
				BeneficiaryType:              nil,
				BeneficiaryProviderId:        nil,
				BeneficiaryProviderAccountId: nil,
				Version:                      nil,
			},
		},
	}

	isValid, err := deepvalidator.NewProcessor().
		RegisterCondition(query).
		ValidateMultipleStructs(data)

	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Is valid:", isValid)
	}
}
