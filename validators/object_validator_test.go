package validators

import (
	"reflect"
	"testing"
)

func TestProcessStructsToMap(t *testing.T) {
	type PartnerStatementUpdatedEventPayload struct {
		PartnerId string
	}

	type PartnerStatementUpdatedEvent struct {
		EventSource string
		Payload     PartnerStatementUpdatedEventPayload
	}

	type TransactionUpdatedEventPayload struct {
		TransactionId string
	}

	type TransactionUpdatedEvent struct {
		EventSource string
		Payload     TransactionUpdatedEventPayload
	}
	
	tests := []struct {
		name     string
		input    interface{}
		expected map[string]interface{}
	}{
		{
			name: "Single struct",
			input: PartnerStatementUpdatedEvent{
				EventSource: "source",
				Payload: PartnerStatementUpdatedEventPayload{
					PartnerId: "bca",
				},
			},
			expected: map[string]interface{}{
				"EventSource":       "source",
				"Payload.PartnerId": "bca",
			},
		},
		{
			name: "Slice of structs",
			input: []interface{}{
				PartnerStatementUpdatedEvent{
					EventSource: "source1",
					Payload: PartnerStatementUpdatedEventPayload{
						PartnerId: "bca1",
					},
				},
				PartnerStatementUpdatedEvent{
					EventSource: "source2",
					Payload: PartnerStatementUpdatedEventPayload{
						PartnerId: "bca2",
					},
				},
			},
			expected: map[string]interface{}{
				"PartnerStatementUpdatedEvent.EventSource":       "source2",
				"PartnerStatementUpdatedEvent.Payload.PartnerId": "bca2",
			},
		},
		{
			name: "Nested structs",
			input: TransactionUpdatedEvent{
				EventSource: "source",
				Payload: TransactionUpdatedEventPayload{
					TransactionId: "txn123",
				},
			},
			expected: map[string]interface{}{
				"EventSource":           "source",
				"Payload.TransactionId": "txn123",
			},
		},
		{
			name: "Struct with nil pointers",
			input: PartnerStatementUpdatedEvent{
				EventSource: "",
				Payload: PartnerStatementUpdatedEventPayload{
					PartnerId: "",
				},
			},
			expected: map[string]interface{}{
				"EventSource":       "",
				"Payload.PartnerId": "",
			},
		},
		{
			name:     "Empty struct",
			input:    struct{}{},
			expected: map[string]interface{}{},
		},
		{
			name:     "Nil input",
			input:    nil,
			expected: map[string]interface{}{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := processStructsToMap("", tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("processStructsToMap() = %v, want %v", result, tt.expected)
			}
		})
	}
}
