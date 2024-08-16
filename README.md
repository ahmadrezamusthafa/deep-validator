# Deep Validator

`deep-validator` is a Golang library designed to validate objects based on complex conditions defined by query strings.
This library supports validation across multiple struct types, allowing developers to enforce complex business rules
with ease.

## Features

- **Single Struct Validation**: Validate a single struct based on specified conditions.
- **Multi-Struct Validation**: Validate multiple structs simultaneously with conditions spanning across them.
- **Attribute Presence Checking**: Ensure required attributes are present in the data being validated.
- **Error Handling**: Handles scenarios where the input data is nil or empty.

## Installation

```bash
go get github.com/ahmadrezamusthafa/deep-validators
```

## Usage

### Supported Operators

The `deep-validator` library supports a variety of operators for use in query strings. Below is a table of the supported operators:

| Operator                    | Symbol | Description                                                  |
|-----------------------------|--------|--------------------------------------------------------------|
| **Equal**                   | `=`    | Checks if a field is equal to a specified value.             |
| **Not Equal**               | `!=`   | Checks if a field is not equal to a specified value.         |
| **Less Than**               | `<`    | Checks if a field is less than a specified value.            |
| **Less Than or Equal**      | `<=`   | Checks if a field is less than or equal to a specified value. |
| **Greater Than**            | `>`    | Checks if a field is greater than a specified value.         |
| **Greater Than or Equal**   | `>=`   | Checks if a field is greater than or equal to a specified value.|
| **Contains**                | `\|=`  | Checks if a field contains a specified substring.              |
| **Contains Regex Match**    | `\|~`  | Checks if a field matches a specified regex pattern.           |

### Basic Validation

To validate a single struct:

```go
package main

import (
	"fmt"
	"github.com/ahmadrezamusthafa/deep-validator"
)

type FirstStruct struct {
	ID       string `json:"id"`
	MemberID string `json:"member_id"`
	Division string `json:"division"`
}

func main() {
	query := `id=123`
	data := FirstStruct{
		ID:       "123",
		MemberID: "345",
		Division: "engineering",
	}

	isValid, err := deepvalidator.NewProcessor().
		RegisterCondition(query).
		ValidateStruct(data)

	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Is valid:", isValid)
	}
}
```

### Multi-Struct Validation

You can validate multiple structs together:

```go
package main

import (
	"fmt"
	"github.com/ahmadrezamusthafa/deep-validator"
)

type FirstStruct struct {
	ID       string `json:"id"`
	MemberID string `json:"member_id"`
	Division string `json:"division"`
}

type SecondStruct struct {
	Name string `json:"name"`
}

type ThirdStruct struct {
	Type    string `json:"type"`
	Segment string `json:"segment"`
}

func main() {
	query := `(FirstStruct.id=123 || SecondStruct.name=Test || ThirdStruct.segment=new-member) && (FirstStruct.member_id=345 && SecondStruct.name=Test) && ThirdStruct.type=ABC`
	data := []interface{}{
		FirstStruct{
			ID:       "123",
			MemberID: "345",
			Division: "engineering",
		},
		SecondStruct{
			Name: "Test",
		},
		ThirdStruct{
			Type:    "ABC",
			Segment: "new-member",
		},
	}

	isValid, err := deepvalidator.NewProcessor().
		RegisterCondition(query).
		ValidateMultipleStructs(data)

	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Is valid:", isValid)
	}
}
```

## Unit Tests

The library comes with comprehensive unit tests to validate its functionality. You can run the tests using:

```bash
go test ./...
```

## Contributing

Feel free to open issues or submit pull requests for any bugs or enhancements.

