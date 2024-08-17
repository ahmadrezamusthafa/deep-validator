package utils

import (
	"testing"
)

func BenchmarkReplacePlaceholders(b *testing.B) {
	query := "(transactionId = {{transaction_id}} && (confirmedAt >= {{start_confirmed_at}} && confirmedAt <= {{end_confirmed_at}}) && beneficiaryProviderAccountId = {{beneficiary_account_id}})"
	attributes := map[string]interface{}{
		"transaction_id":         "616994753",
		"start_confirmed_at":     12312312312,
		"end_confirmed_at":       122131231231,
		"beneficiary_account_id": "1770019821328",
	}
	for n := 0; n < b.N; n++ {
		ReplacePlaceholders(query, attributes)
	}
}
