package validators

import (
	"reflect"
	"testing"
)

func TestStructsToMap(t *testing.T) {
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
		name         string
		input        interface{}
		removePrefix bool
		expected     map[string]interface{}
	}{
		{
			name: "Single struct",
			input: PartnerStatementUpdatedEvent{
				EventSource: "source",
				Payload: PartnerStatementUpdatedEventPayload{
					PartnerId: "bca",
				},
			},
			removePrefix: false,
			expected: map[string]interface{}{
				"event_source":       "source",
				"payload.partner_id": "bca",
			},
		},
		{
			name: "Single struct - no prefix",
			input: PartnerStatementUpdatedEvent{
				EventSource: "source",
				Payload: PartnerStatementUpdatedEventPayload{
					PartnerId: "bca",
				},
			},
			removePrefix: true,
			expected: map[string]interface{}{
				"event_source": "source",
				"partner_id":   "bca",
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
				TransactionUpdatedEvent{
					EventSource: "source2",
					Payload: TransactionUpdatedEventPayload{
						TransactionId: "123455",
					},
				},
			},
			removePrefix: false,
			expected: map[string]interface{}{
				"PartnerStatementUpdatedEvent.event_source":       "source1",
				"TransactionUpdatedEvent.event_source":            "source2",
				"PartnerStatementUpdatedEvent.payload.partner_id": "bca1",
				"TransactionUpdatedEvent.payload.transaction_id":  "123455",
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
				TransactionUpdatedEvent{
					EventSource: "source2",
					Payload: TransactionUpdatedEventPayload{
						TransactionId: "123455",
					},
				},
			},
			removePrefix: true,
			expected: map[string]interface{}{
				"event_source":   "source2",
				"partner_id":     "bca1",
				"transaction_id": "123455",
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
			removePrefix: false,
			expected: map[string]interface{}{
				"event_source":           "source",
				"payload.transaction_id": "txn123",
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
			removePrefix: false,
			expected: map[string]interface{}{
				"event_source":       "",
				"payload.partner_id": "",
			},
		},
		{
			name:         "Empty struct",
			input:        struct{}{},
			removePrefix: false,
			expected:     map[string]interface{}{},
		},
		{
			name:         "Nil input",
			input:        nil,
			removePrefix: false,
			expected:     map[string]interface{}{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := structsToMap(tt.input, tt.removePrefix)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("processStructsToMap() = %v, want %v", result, tt.expected)
			}
		})
	}
}
