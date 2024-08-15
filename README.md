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

