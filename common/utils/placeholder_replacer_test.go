package utils

import (
	"testing"
)

func TestReplacePlaceholders(t *testing.T) {
	tests := []struct {
		name        string
		inputStr    string
		attributes  map[string]interface{}
		expectedStr string
	}{
		{
			name:     "Replace all placeholders",
			inputStr: "(transactionId = {{transaction_id}} && (confirmedAt >= {{start_confirmed_at}} && confirmedAt <= {{end_confirmed_at}}) && beneficiaryProviderAccountId = {{beneficiary_account_id}})",
			attributes: map[string]interface{}{
				"transaction_id":         "616994753",
				"start_confirmed_at":     12312312312,
				"end_confirmed_at":       122131231231,
				"beneficiary_account_id": "1770019821328",
			},
			expectedStr: "(transactionId = 616994753 && (confirmedAt >= 12312312312 && confirmedAt <= 122131231231) && beneficiaryProviderAccountId = 1770019821328)",
		},
		{
			name:     "Missing attribute in map",
			inputStr: "(transactionId = {{transaction_id}} && beneficiaryProviderAccountId = {{beneficiary_account_id}} && unknown = {{unknown}})",
			attributes: map[string]interface{}{
				"transaction_id":         "616994753",
				"beneficiary_account_id": "1770019821328",
			},
			expectedStr: "(transactionId = 616994753 && beneficiaryProviderAccountId = 1770019821328 && unknown = {{unknown}})",
		},
		{
			name:        "Empty string",
			inputStr:    "",
			attributes:  map[string]interface{}{},
			expectedStr: "",
		},
		{
			name:     "No placeholders",
			inputStr: "(transactionId = 616994753 && beneficiaryProviderAccountId = 1770019821328)",
			attributes: map[string]interface{}{
				"transaction_id":         "616994753",
				"beneficiary_account_id": "1770019821328",
			},
			expectedStr: "(transactionId = 616994753 && beneficiaryProviderAccountId = 1770019821328)",
		},
		{
			name:     "Partial placeholders",
			inputStr: "(transactionId = {{transaction_id}} && confirmedAt >= {{start_confirmed_at}} && confirmedAt <= {{end_confirmed_at}})",
			attributes: map[string]interface{}{
				"transaction_id": "616994753",
			},
			expectedStr: "(transactionId = 616994753 && confirmedAt >= {{start_confirmed_at}} && confirmedAt <= {{end_confirmed_at}})",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ReplacePlaceholders(tt.inputStr, tt.attributes)
			if result != tt.expectedStr {
				t.Errorf("got %s, want %s", result, tt.expectedStr)
			}
		})
	}
}
