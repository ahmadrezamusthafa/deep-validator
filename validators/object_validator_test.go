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
				"EventSource":       "source",
				"Payload.PartnerId": "bca",
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
				"EventSource": "source",
				"PartnerId":   "bca",
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
				"PartnerStatementUpdatedEvent.EventSource":       "source1",
				"TransactionUpdatedEvent.EventSource":            "source2",
				"PartnerStatementUpdatedEvent.Payload.PartnerId": "bca1",
				"TransactionUpdatedEvent.Payload.TransactionId":  "123455",
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
				"EventSource":   "source2",
				"PartnerId":     "bca1",
				"TransactionId": "123455",
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
			removePrefix: false,
			expected: map[string]interface{}{
				"EventSource":       "",
				"Payload.PartnerId": "",
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
