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
	ValidateStruct(data interface{}) (isValid bool, err error)
	ValidateMultipleStructs(data ...interface{}) (isValid bool, err error)
	ValidateCondition(inputCondition structs.Condition) (isValid bool, err error)
	FilterSlice(data interface{}) (result interface{}, err error)
}

type processor struct {
	condition *structs.Condition
}

type validator struct {
	processor *processor
}

func NewProcessor() Processor {
	return &processor{}
}

func newValidator(processor *processor) Validator {
	return &validator{
		processor: processor,
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
		return newValidator(p)
	}
	p.condition = &condition
	return newValidator(p)
}

func (v *validator) ValidateStruct(data interface{}) (isValid bool, err error) {
	if v.processor.condition == nil {
		return false, errors.New("condition is nil")
	}
	con := validators.Condition{Condition: v.processor.condition}
	return con.Validate(data)
}

func (v *validator) ValidateMultipleStructs(data ...interface{}) (isValid bool, err error) {
	if v.processor.condition == nil {
		return false, errors.New("condition is nil")
	}
	con := validators.Condition{Condition: v.processor.condition}
	return con.ValidateObjects(data...)
}

func (v *validator) ValidateCondition(inputCondition structs.Condition) (isValid bool, err error) {
	if v.processor.condition == nil {
		return false, errors.New("condition is nil")
	}
	con := validators.Condition{Condition: v.processor.condition}
	return con.ValidateCondition(inputCondition)
}

func (v *validator) FilterSlice(data interface{}) (result interface{}, err error) {
	if v.processor.condition == nil {
		return false, errors.New("condition is nil")
	}
	con := validators.Condition{Condition: v.processor.condition}
	return con.FilterSlice(data)
}
