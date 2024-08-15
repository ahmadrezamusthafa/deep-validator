package validators

import (
	"errors"
	"fmt"
	"github.com/ahmadrezamusthafa/deep-validator/common/utils"
	datetimeformats "github.com/ahmadrezamusthafa/deep-validator/consts/datetime-formats"
	errormessages "github.com/ahmadrezamusthafa/deep-validator/consts/error-messages"
	"github.com/ahmadrezamusthafa/deep-validator/consts/logical-operators"
	"github.com/ahmadrezamusthafa/deep-validator/consts/operators"
	"github.com/ahmadrezamusthafa/deep-validator/enums/value-types"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func (c *Condition) Validate(data interface{}) (isValid bool, err error) {
	if data == nil {
		return false, fmt.Errorf(errormessages.ErrorMessageInvalidData, "nil")
	}
	rType := reflect.TypeOf(data)
	switch rType.Kind() {
	case reflect.Struct, reflect.Map:
		isValid, _, err := c.validateAttribute(rType, data)
		return isValid, err
	default:
		return false, fmt.Errorf(errormessages.ErrorMessageInvalidType, "struct")
	}
}

func (c *Condition) ValidateObjects(data ...interface{}) (isValid bool, err error) {
	if data == nil {
		return false, fmt.Errorf(errormessages.ErrorMessageInvalidData, "nil")
	}
	rType := reflect.TypeOf(data)
	switch rType.Kind() {
	case reflect.Slice:
		dataMap := structsToMap(data)
		return c.Validate(dataMap)
	default:
		return false, fmt.Errorf(errormessages.ErrorMessageInvalidType, "slice")
	}
}

func (c *Condition) FilterSlice(data interface{}) (result interface{}, err error) {
	if data == nil {
		return result, fmt.Errorf(errormessages.ErrorMessageInvalidData, "nil")
	}
	rType := reflect.TypeOf(data)
	switch rType.Kind() {
	case reflect.Slice:
		rValue := reflect.ValueOf(data)
		rSlice := reflect.MakeSlice(rType, 0, 1)
		for i := 0; i < rValue.Len(); i++ {
			obj := rValue.Index(i).Interface()
			isValid, err := c.Validate(obj)
			if err != nil {
				return rSlice, err
			}
			if isValid {
				rSlice = reflect.Append(rSlice, rValue.Index(i))
			}
		}
		result = rSlice.Interface()
		return
	default:
		return result, fmt.Errorf(errormessages.ErrorMessageInvalidType, "slice")
	}
}

func (c *Condition) prepareDataFromSlice(data interface{}) (interface{}, error) {
	var preparedData interface{}
	rValue := reflect.ValueOf(data)
	if rValue.Type().Kind() != reflect.Slice {
		return false, fmt.Errorf(errormessages.ErrorMessageInvalidType, "slice")
	}
	if rValue.Len() == 0 {
		return false, fmt.Errorf(errormessages.ErrorMessageInvalidData, "empty slice")
	}

	firstValue := rValue.Index(0).Interface()
	rFirstValue := reflect.ValueOf(firstValue)
	if firstValue == nil {
		return false, fmt.Errorf(errormessages.ErrorMessageInvalidData, "nil")
	}
	switch rFirstValue.Type().Kind() {
	case reflect.Struct:
		preparedData = firstValue
	default:
		length := rFirstValue.Len()
		switch length {
		case 0:
			return false, fmt.Errorf(errormessages.ErrorMessageInvalidData, "empty slice")
		case 1:
			preparedData = rFirstValue.Index(0).Interface()
		default:
			mapObj := make(map[string]interface{})
			mapValue := reflect.MakeMap(reflect.TypeOf(mapObj))
			for i := 0; i < length; i++ {
				rDetailValue := reflect.ValueOf(rFirstValue.Index(i).Interface())
				mapValue.SetMapIndex(reflect.ValueOf(rDetailValue.Type().Name()), rDetailValue)
			}
			preparedData = mapValue.Interface()
		}
	}
	return preparedData, nil
}

func (c *Condition) validateAttribute(rType reflect.Type, data interface{}) (isValid, isSkip bool, err error) {
	if len(c.Conditions) > 0 {
		for i, subCondition := range c.Conditions {
			con := Condition{Condition: subCondition}
			isSubValid, isSkip, err := con.validateAttribute(rType, data)
			if err != nil {
				return false, false, err
			}
			if isSkip {
				continue
			}
			if i == 0 {
				isValid = isSubValid
			} else {
				if subCondition.Operator == logicaloperators.LogicalOperatorOr {
					isValid = isValid || isSubValid
				} else {
					isValid = isValid && isSubValid
				}
			}
		}
	} else {
		switch rType.Kind() {
		case reflect.Map:
			if value, ok := data.(map[string]interface{}); ok {
				isValid, isSkip, err = c.validateMapValue(value)
			} else {
				return false, false, errors.New(errormessages.ErrorMessageUnableToCastObject)
			}
		default:
			isValid, err = c.validateStructValue("", data)
		}
	}
	return
}

func (c *Condition) validateStructValue(prefix string, data interface{}) (isValid bool, err error) {
	rValue := reflect.ValueOf(data)
	if rValue.Type().Kind() != reflect.Struct {
		return false, fmt.Errorf(errormessages.ErrorMessageInvalidType, "struct")
	}
	for i := 0; i < rValue.NumField(); i++ {
		field := rValue.Field(i)
		typeField := rValue.Type().Field(i)
		tag := typeField.Name
		jsonTag, ok := typeField.Tag.Lookup("json")
		if ok && jsonTag != "" {
			tag = jsonTag
		}
		tag = prefix + tag

		if tag == c.Attribute.Name {
			var conditionValue interface{}
			validationType := valuetypes.Numeric
			value := field.Interface()
			operator := c.Attribute.Operator

			if field.Kind() == reflect.Ptr && field.IsNil() {
				return false, nil
			}

			switch value.(type) {
			case int, int64:
				value = utils.InterfaceToInt64(value)
				conditionValue, err = strconv.ParseInt(c.Attribute.Value, 10, 64)
			case *int, *int64:
				value = utils.InterfacePtrToInt64(value)
				conditionValue, err = strconv.ParseInt(c.Attribute.Value, 10, 64)
			case float32, float64:
				value = utils.InterfaceToFloat64(value)
				conditionValue, err = strconv.ParseFloat(c.Attribute.Value, 64)
			case *float32, *float64:
				value = utils.InterfacePtrToFloat64(value)
				conditionValue, err = strconv.ParseInt(c.Attribute.Value, 10, 64)
			case time.Time:
				validationType = valuetypes.Date
				conditionValue, err = time.Parse(datetimeformats.DateTimeFormat, c.Attribute.Value)
			case *time.Time:
				validationType = valuetypes.Date
				res, ok := value.(*time.Time)
				if ok {
					value = *res
				}
				conditionValue, err = time.Parse(datetimeformats.DateTimeFormat, c.Attribute.Value)
			case bool:
				validationType = valuetypes.Alphanumeric
				conditionValue = utils.StringToBool(c.Attribute.Value)
			default:
				validationType = valuetypes.Alphanumeric
				conditionValue = c.Attribute.Value
			}
			if err != nil {
				return false, err
			}

			switch operator {
			case operators.OperatorEqual:
				isValid = value == conditionValue
			default:
				switch validationType {
				case valuetypes.Date:
					isValid = validateTime(value, operator, conditionValue)
				default:
					isValid = validateNumeric(value, operator, conditionValue)
				}
			}
		}
	}
	return
}

func (c *Condition) validateMapValue(data map[string]interface{}) (isValid, isSkip bool, err error) {
	if data == nil || len(data) == 0 {
		return false, false, fmt.Errorf(errormessages.ErrorMessageInvalidData, "nil")
	}
	isSkip = true
	if v, ok := data[c.Attribute.Name]; ok {
		isSkip = false
		isValid, err = c.validateMap(c.Attribute.Name, v)
		if err != nil {
			return false, false, err
		}
	} else {
		return false, false, err
	}
	return
}

func (c *Condition) validateMap(key string, value interface{}) (isValid bool, err error) {
	if key == c.Attribute.Name {
		var conditionValue interface{}
		validationType := valuetypes.Numeric
		operator := c.Attribute.Operator

		if value == nil {
			return false, nil
		}

		switch value.(type) {
		case int, int64:
			value = utils.InterfaceToInt64(value)
			conditionValue, err = strconv.ParseInt(c.Attribute.Value, 10, 64)
		case *int, *int64:
			value = utils.InterfacePtrToInt64(value)
			conditionValue, err = strconv.ParseInt(c.Attribute.Value, 10, 64)
		case float32, float64:
			value = utils.InterfaceToFloat64(value)
			conditionValue, err = strconv.ParseFloat(c.Attribute.Value, 64)
		case *float32, *float64:
			value = utils.InterfacePtrToFloat64(value)
			conditionValue, err = strconv.ParseInt(c.Attribute.Value, 10, 64)
		case time.Time:
			validationType = valuetypes.Date
			conditionValue, err = time.Parse(datetimeformats.DateTimeFormat, c.Attribute.Value)
		case *time.Time:
			validationType = valuetypes.Date
			res, ok := value.(*time.Time)
			if ok {
				value = *res
			}
			conditionValue, err = time.Parse(datetimeformats.DateTimeFormat, c.Attribute.Value)
		case bool:
			validationType = valuetypes.Alphanumeric
			conditionValue = utils.StringToBool(c.Attribute.Value)
		default:
			validationType = valuetypes.Alphanumeric
			conditionValue = c.Attribute.Value
		}
		if err != nil {
			return false, err
		}

		switch operator {
		case operators.OperatorEqual:
			isValid = value == conditionValue
		default:
			switch validationType {
			case valuetypes.Date:
				isValid = validateTime(value, operator, conditionValue)
			default:
				isValid = validateNumeric(value, operator, conditionValue)
			}
		}
	}
	return
}

func validateTime(firstVal interface{}, operator string, secondVal interface{}) bool {
	firstTime, ok := firstVal.(time.Time)
	if !ok {
		return false
	}
	secondTime, ok := secondVal.(time.Time)
	if !ok {
		return false
	}

	switch operator {
	case operators.OperatorGreaterThan:
		return firstTime.After(secondTime)
	case operators.OperatorLessThan:
		return firstTime.Before(secondTime)
	case operators.OperatorGreaterThanEqual:
		return firstTime.After(secondTime) || firstTime.Equal(secondTime)
	default:
		return firstTime.Before(secondTime) || firstTime.Equal(secondTime)
	}
}

func validateNumeric(firstVal interface{}, operator string, secondVal interface{}) bool {
	firstFloat, ok := firstVal.(float64)
	if !ok {
		firstInt, ok := firstVal.(int64)
		if !ok {
			return false
		}
		firstFloat = float64(firstInt)
	}
	secondFloat, ok := secondVal.(float64)
	if !ok {
		secondInt, ok := secondVal.(int64)
		if !ok {
			return false
		}
		secondFloat = float64(secondInt)
	}

	switch operator {
	case operators.OperatorGreaterThan:
		return firstFloat > secondFloat
	case operators.OperatorLessThan:
		return firstFloat < secondFloat
	case operators.OperatorGreaterThanEqual:
		return firstFloat >= secondFloat
	default:
		return firstFloat <= secondFloat
	}
}

func structsToMap(data interface{}) map[string]interface{} {
	return processStructsToMap(reflect.TypeOf(data).Name(), data)
}

func processStructsToMap(prefix string, data interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	rValue := reflect.ValueOf(data)
	if rValue.Kind() == reflect.Ptr {
		if rValue.IsNil() {
			return result
		}
		rValue = rValue.Elem()
	}

	if rValue.Kind() == reflect.Slice {
		for i := 0; i < rValue.Len(); i++ {
			item := rValue.Index(i).Interface()
			if item == nil {
				return result
			}
			itemPrefix := reflect.TypeOf(item).Name()
			nestedMap := processStructsToMap(itemPrefix, item)
			for k, v := range nestedMap {
				result[k] = v
			}
		}
		return result
	}

	if rValue.Kind() != reflect.Struct {
		return result
	}

	rType := rValue.Type()
	for i := 0; i < rValue.NumField(); i++ {
		field := rValue.Field(i)
		typeField := rType.Field(i)
		key := typeField.Name

		jsonTag := typeField.Tag.Get("json")
		if jsonTag != "" && jsonTag != "-" {
			key = strings.Split(jsonTag, ",")[0]
		}

		if prefix != "" {
			key = prefix + "." + key
		}

		if field.Kind() == reflect.Struct {
			nestedMap := processStructsToMap(key, field.Interface())
			for k, v := range nestedMap {
				result[k] = v
			}
		} else if field.Kind() == reflect.Ptr && !field.IsNil() {
			if field.Elem().Kind() == reflect.Struct {
				nestedMap := processStructsToMap(key, field.Elem().Interface())
				for k, v := range nestedMap {
					result[k] = v
				}
			} else {
				result[key] = field.Interface()
			}
		} else {
			if field.CanInterface() {
				result[key] = field.Interface()
			}
		}
	}

	return result
}
