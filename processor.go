package deepvalidator

import (
	"errors"
	structgen "github.com/ahmadrezamusthafa/deep-validator/struct-gen"
	"github.com/ahmadrezamusthafa/deep-validator/structs"
	"github.com/ahmadrezamusthafa/deep-validator/validators"
)

type Processor interface {
	RegisterCondition(astQuery string) Validator
}

type Validator interface {
	SetRemovePrefix(value bool) Validator
	ValidateStruct(data interface{}) (isValid bool, err error)
	ValidateMultipleStructs(data ...interface{}) (isValid bool, err error)
	ValidateCondition(inputCondition structs.Condition) (isValid bool, err error)
	FilterSlice(data interface{}) (result interface{}, err error)
	GetCondition() *structs.Condition
}

type processor struct{}

type validator struct {
	attributeNames     map[string]interface{}
	conditionValidator validators.ConditionValidator
}

func NewProcessor() Processor {
	return &processor{}
}

func newValidator(attributeNames map[string]interface{}, conditionValidator validators.ConditionValidator) Validator {
	return &validator{
		attributeNames:     attributeNames,
		conditionValidator: conditionValidator,
	}
}

/*
GenerateCondition
-----------------------------------------------------------------------
is a function to generate condition object as validators that used by
  - ValidateStruct
  - ValidateMultipleStructs
  - ValidateCondition

Param:
@astQuery is abstract syntax tree query
*/
func GenerateCondition(astQuery string) (structs.Condition, error) {
	var gen structgen.StructGen
	return gen.GenerateCondition(astQuery)
}

func (p *processor) RegisterCondition(astQuery string) Validator {
	var gen structgen.StructGen
	condition, err := gen.GenerateCondition(astQuery)
	if err != nil {
		return newValidator(nil, nil)
	}
	return newValidator(gen.AttributeNames, validators.NewConditionValidator(&condition))
}

func (v *validator) SetRemovePrefix(value bool) Validator {
	if v.conditionValidator.GetCondition() == nil {
		return v
	}
	v.conditionValidator.SetRemovePrefix(value)
	return v
}

func (v *validator) ValidateStruct(data interface{}) (isValid bool, err error) {
	if v.conditionValidator.GetCondition() == nil {
		return false, errors.New("condition is nil")
	}
	return v.conditionValidator.Validate(data)
}

func (v *validator) ValidateMultipleStructs(data ...interface{}) (isValid bool, err error) {
	if v.conditionValidator.GetCondition() == nil {
		return false, errors.New("condition is nil")
	}
	return v.conditionValidator.ValidateObjects(v.attributeNames, data...)
}

func (v *validator) ValidateCondition(inputCondition structs.Condition) (isValid bool, err error) {
	if v.conditionValidator.GetCondition() == nil {
		return false, errors.New("condition is nil")
	}
	return v.conditionValidator.ValidateCondition(inputCondition)
}

func (v *validator) FilterSlice(data interface{}) (result interface{}, err error) {
	if v.conditionValidator.GetCondition() == nil {
		return false, errors.New("condition is nil")
	}
	return v.conditionValidator.FilterSlice(data)
}

func (v *validator) GetCondition() *structs.Condition {
	if v.conditionValidator.GetCondition() == nil {
		return nil
	}
	return v.conditionValidator.GetCondition()
}
