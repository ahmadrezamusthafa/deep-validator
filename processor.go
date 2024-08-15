package deepvalidator

import (
	structgen "github.com/ahmadrezamusthafa/deep-validator/struct-gen"
	"github.com/ahmadrezamusthafa/deep-validator/structs"
	"github.com/ahmadrezamusthafa/deep-validator/validator"
)

type Processor struct{}

func NewProcessor() *Processor {
	return &Processor{}
}

/*
GenerateCondition
-----------------------------------------------------------------------
is a function to generate condition object as validator that used by
  - Validate
  - ValidateObjects
  - ValidateCondition

Param:
@astQuery is abstract syntax tree query
*/
func GenerateCondition(astQuery string) (structs.Condition, error) {
	var gen structgen.StructGen
	return gen.GenerateCondition(astQuery)
}

func Validate(referenceCondition structs.Condition, data interface{}) (isValid bool, err error) {
	con := validator.Condition{Condition: &referenceCondition}
	return con.Validate(data)
}

func ValidateObjects(referenceCondition structs.Condition, data ...interface{}) (isValid bool, err error) {
	con := validator.Condition{Condition: &referenceCondition}
	return con.ValidateObjects(data...)
}

func ValidateCondition(referenceCondition structs.Condition, inputCondition structs.Condition) (isValid bool, err error) {
	con := validator.Condition{Condition: &referenceCondition}
	return con.ValidateCondition(inputCondition)
}

func FilterSlice(referenceCondition structs.Condition, data interface{}) (result interface{}, err error) {
	con := validator.Condition{Condition: &referenceCondition}
	return con.FilterSlice(data)
}
