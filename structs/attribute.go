package structs

import (
	"github.com/ahmadrezamusthafa/deep-validator/enums/value-types"
)

type Attribute struct {
	Name     string               `json:"name"`
	Operator string               `json:"operator"`
	Value    string               `json:"value"`
	Type     valuetypes.ValueType `json:"type,omitempty"`
}

type TokenAttribute struct {
	Value          string
	IsAlphanumeric bool
	HasCalled      bool
}
