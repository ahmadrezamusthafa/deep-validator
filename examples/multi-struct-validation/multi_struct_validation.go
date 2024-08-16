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
	query := `(ID=123 || Name=Test || Segment=new-member) && (MemberID=345 && Name=Test) && Type=ABC`
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
