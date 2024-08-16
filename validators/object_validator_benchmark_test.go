package validators

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func BenchmarkStructsToMap(b *testing.B) {
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

	var jsTrx = `
		{
		  "payload": {
			"transactionId": "616994753",
			"transactionType": "WithdrawalTransfer.v1",
			"totalAmount": 50000,
			"senderProviderId": "mandiri",
			"senderType": "FlipBankAccount",
			"senderProviderAccountId": "1570009951832",
			"transactionReference": "20240624141310607",
			"status": "Completed",
			"createdAt": "2024-06-24T14:12:10+07:00",
			"confirmedAt": "2024-06-24T14:13:08+07:00",
			"completedAt": "2024-06-24T14:13:11+07:00",
			"currency": "IDR",
			"beneficiaryType": "BankAccount",
			"beneficiaryProviderId": "mandiri",
			"beneficiaryProviderAccountId": "1770019821328",
			"version": 1
		  }
		}`

	var trx TransactionUpdatedEvent
	err := json.Unmarshal([]byte(jsTrx), &trx)
	if err != nil {
		fmt.Println(err)
	}

	attributeNames := map[string]interface{}{
		"transactionId":                nil,
		"totalAmount":                  nil,
		"beneficiaryProviderAccountId": nil,
		"senderProviderId":             nil,
		"createdAt":                    nil,
	}

	for n := 0; n < b.N; n++ {
		structsToMap(attributeNames, trx)
	}
}
