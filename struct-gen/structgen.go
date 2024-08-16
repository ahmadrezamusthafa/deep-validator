package structgen

import (
	"bytes"
	bytescodes "github.com/ahmadrezamusthafa/deep-validator/consts/bytes-codes"
	datetimeformats "github.com/ahmadrezamusthafa/deep-validator/consts/datetime-formats"
	"github.com/ahmadrezamusthafa/deep-validator/consts/logical-operators"
	"github.com/ahmadrezamusthafa/deep-validator/consts/operators"
	"github.com/ahmadrezamusthafa/deep-validator/enums/value-types"
	"github.com/ahmadrezamusthafa/deep-validator/structs"
	"time"
)

type StructGen struct {
}

var (
	operatorMap = map[string]interface{}{
		operators.OperatorEqual:              nil,
		operators.OperatorLessThan:           nil,
		operators.OperatorGreaterThan:        nil,
		operators.OperatorLessThanEqual:      nil,
		operators.OperatorGreaterThanEqual:   nil,
		operators.OperatorContains:           nil,
		operators.OperatorContainsRegexMatch: nil,
	}

	logicalOperatorMap = map[string]string{
		logicaloperators.LogicalOperatorAndSyntax: logicaloperators.LogicalOperatorAnd,
		logicaloperators.LogicalOperatorOrSyntax:  logicaloperators.LogicalOperatorOr,
	}
)

func (s *StructGen) GenerateCondition(query string) (structs.Condition, error) {
	tokenAttributes := getTokenAttributes(query)
	if len(tokenAttributes) == 0 {
		return structs.Condition{Attribute: &structs.Attribute{}}, nil
	}
	_, condition := buildCondition(structs.Condition{}, tokenAttributes)
	return condition, nil
}

func buildCondition(condition structs.Condition, attrs []*structs.TokenAttribute) (int, structs.Condition) {
	var (
		conditionItem *structs.Condition
		lastPos       int
		operator      string
	)
	for i := 0; i < len(attrs); i++ {
		lastPos = i
		attr := attrs[i]
		if attr.HasCalled {
			continue
		}
		attr.HasCalled = true
		if attr.Value == ")" {
			break
		}
		if attr.Value == "(" {
			newCondition := structs.Condition{
				Operator: operator,
			}
			lastPos, resp := buildCondition(newCondition, attrs[i+1:])
			condition.Conditions = append(condition.Conditions, &resp)
			i = i + lastPos + 1
			continue
		}

		if val, ok := logicalOperatorMap[attr.Value]; ok {
			operator = val
			conditionItem = nil
		} else if _, ok := operatorMap[attr.Value]; ok {
			if conditionItem != nil {
				conditionItem.Attribute.Operator = attr.Value
			}
		} else {
			if conditionItem == nil {
				conditionItem = &structs.Condition{
					Attribute: &structs.Attribute{
						Name: attr.Value,
					},
				}
				conditionItem.Attribute = &structs.Attribute{
					Name: attr.Value,
				}
			} else {
				conditionItem.Attribute.Value = attr.Value
				if !attr.IsAlphanumeric {
					conditionItem.Attribute.Type = getValueType(attr.Value)
				}
				if condition.Conditions == nil {
					condition.Conditions = []*structs.Condition{}
				}
				conditionItem.Operator = operator
				condition.Conditions = append(condition.Conditions, conditionItem)
			}
		}
	}
	return lastPos, condition
}

func getTokenAttributes(query string) []*structs.TokenAttribute {
	var tokenAttributes []*structs.TokenAttribute
	buffer := &bytes.Buffer{}
	isOpenQuote := false
	isAlphanumeric := false
	for _, char := range query {
		switch char {
		case ' ', '\n', '\'':
			if !isOpenQuote {
				continue
			} else {
				buffer.WriteRune(char)
			}
		case '|', '&', '<', '>':
			if buffer.Len() > 0 {
				bufBytes := buffer.Bytes()
				switch bufBytes[0] {
				case bytescodes.ByteVerticalBar:
					tokenAttributes = appendAttribute(tokenAttributes, buffer, logicaloperators.LogicalOperatorOrSyntax, isAlphanumeric)
					isAlphanumeric = false
				case bytescodes.ByteAmpersand:
					tokenAttributes = appendAttribute(tokenAttributes, buffer, logicaloperators.LogicalOperatorAndSyntax, isAlphanumeric)
					isAlphanumeric = false
				default:
					tokenAttributes = appendAttribute(tokenAttributes, buffer, string(bufBytes), isAlphanumeric)
					isAlphanumeric = false
					buffer.WriteRune(char)
				}
			} else {
				buffer.WriteRune(char)
			}
		case '=', '(', ')', '~':
			if buffer.Len() > 0 {
				bufBytes := buffer.Bytes()
				switch bufBytes[0] {
				case bytescodes.ByteLessThan, bytescodes.ByteGreaterThan:
					tokenAttributes = appendAttribute(tokenAttributes, buffer, string(bufBytes)+string(char), isAlphanumeric)
					isAlphanumeric = false
					continue
				case bytescodes.ByteVerticalBar:
					tokenAttributes = appendAttribute(tokenAttributes, buffer, string(bufBytes)+string(char), isAlphanumeric)
					isAlphanumeric = false
					continue
				default:
					tokenAttributes = appendAttribute(tokenAttributes, buffer, string(bufBytes), isAlphanumeric)
					isAlphanumeric = false
				}
			}
			tokenAttributes = append(tokenAttributes, &structs.TokenAttribute{
				Value: string(char),
			})
		case '"':
			isOpenQuote = !isOpenQuote
			if !isOpenQuote {
				isAlphanumeric = true
			}
		default:
			if buffer.Len() > 0 {
				bufByte := buffer.Bytes()[0]
				if bufByte == bytescodes.ByteLessThan || bufByte == bytescodes.ByteGreaterThan {
					tokenAttributes = appendAttribute(tokenAttributes, buffer, string(bufByte), isAlphanumeric)
					isAlphanumeric = false
				}
			}
			buffer.WriteRune(char)
		}
	}
	if buffer.Len() > 0 {
		tokenAttributes = appendAttribute(tokenAttributes, buffer, buffer.String(), isAlphanumeric)
		isAlphanumeric = false
	}
	return tokenAttributes
}

func appendAttribute(tokenAttributes []*structs.TokenAttribute, buffer *bytes.Buffer, value string, isAlphanumeric bool) []*structs.TokenAttribute {
	tokenAttributes = append(tokenAttributes, &structs.TokenAttribute{
		Value:          value,
		IsAlphanumeric: isAlphanumeric,
	})
	buffer.Reset()
	return tokenAttributes
}

func getValueType(value string) valuetypes.ValueType {
	varType, indexVal, dotCount := valuetypes.Alphanumeric, 0, 0
	for _, char := range value {
		if char == ',' {
			continue
		}
		if '0' <= char && char <= '9' {
			if indexVal == 0 || (indexVal > 0 && dotCount == 1) {
				varType = valuetypes.Numeric
			}
		} else if char == '.' {
			if indexVal > 0 && varType == valuetypes.Numeric {
				dotCount++
				varType = valuetypes.Alphanumeric
			}
			if dotCount > 1 {
				varType = valuetypes.Alphanumeric
				break
			}
		} else {
			varType = valuetypes.Alphanumeric
			break
		}
		indexVal++
	}
	if varType == valuetypes.Alphanumeric {
		if _, err := time.Parse(datetimeformats.DateTimeFormat, value); err == nil {
			varType = valuetypes.Date
		}
	}
	return varType
}
