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
