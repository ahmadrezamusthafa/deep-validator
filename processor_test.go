package deepvalidator

import (
	"encoding/json"
	structs2 "github.com/ahmadrezamusthafa/deep-validator/structs"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestGenerateCondition(t *testing.T) {
	type args struct {
		query string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Normal case",
			args: args{
				query: `
                      (id = 1 
                      && member_id = 2 )
                      || (
                        division = engineering 
                        || division = finance
                      )
`,
			},
			want:    `{"conditions":[{"conditions":[{"attribute":{"name":"id","operator":"=","value":"1","type":"numeric"}},{"operator":"AND","attribute":{"name":"member_id","operator":"=","value":"2","type":"numeric"}}]},{"operator":"OR","conditions":[{"attribute":{"name":"division","operator":"=","value":"engineering","type":"alphanumeric"}},{"operator":"OR","attribute":{"name":"division","operator":"=","value":"finance","type":"alphanumeric"}}]}]}`,
			wantErr: false,
		},
		{
			name: "Normal case",
			args: args{
				query: `
                      (id = 1 
                      && member_id = 2 )
                      || (
                        division = engineering 
                        || division = finance
                      ) && user_id = 43
`,
			},
			want:    `{"conditions":[{"conditions":[{"attribute":{"name":"id","operator":"=","value":"1","type":"numeric"}},{"operator":"AND","attribute":{"name":"member_id","operator":"=","value":"2","type":"numeric"}}]},{"operator":"OR","conditions":[{"attribute":{"name":"division","operator":"=","value":"engineering","type":"alphanumeric"}},{"operator":"OR","attribute":{"name":"division","operator":"=","value":"finance","type":"alphanumeric"}}]},{"operator":"AND","attribute":{"name":"user_id","operator":"=","value":"43","type":"numeric"}}]}`,
			wantErr: false,
		},
		{
			name: "Normal case",
			args: args{
				query: `
                      id = 1 
                      && member_id = 2 
                      && (
                        division = engineering 
                        || division = finance
                      )
`,
			},
			want:    `{"conditions":[{"attribute":{"name":"id","operator":"=","value":"1","type":"numeric"}},{"operator":"AND","attribute":{"name":"member_id","operator":"=","value":"2","type":"numeric"}},{"operator":"AND","conditions":[{"attribute":{"name":"division","operator":"=","value":"engineering","type":"alphanumeric"}},{"operator":"OR","attribute":{"name":"division","operator":"=","value":"finance","type":"alphanumeric"}}]}]}`,
			wantErr: false,
		},
		{
			name: "Normal case",
			args: args{
				query: `id=1 &&  member_id=2   &&   (division=engineering || division=finance)`,
			},
			want:    `{"conditions":[{"attribute":{"name":"id","operator":"=","value":"1","type":"numeric"}},{"operator":"AND","attribute":{"name":"member_id","operator":"=","value":"2","type":"numeric"}},{"operator":"AND","conditions":[{"attribute":{"name":"division","operator":"=","value":"engineering","type":"alphanumeric"}},{"operator":"OR","attribute":{"name":"division","operator":"=","value":"finance","type":"alphanumeric"}}]}]}`,
			wantErr: false,
		},
		{
			name: "Normal case",
			args: args{
				query: `
                  id = 1 
                  && member_id = 2 
                  && user_id = 3 
                  && (
                    province = jatim 
                    || city = mojokerto
                    || (
                      warehouse_id = 1 
                      && warehouse_detail_id = 2
                    )
                  )
`,
			},
			want:    `{"conditions":[{"attribute":{"name":"id","operator":"=","value":"1","type":"numeric"}},{"operator":"AND","attribute":{"name":"member_id","operator":"=","value":"2","type":"numeric"}},{"operator":"AND","attribute":{"name":"user_id","operator":"=","value":"3","type":"numeric"}},{"operator":"AND","conditions":[{"attribute":{"name":"province","operator":"=","value":"jatim","type":"alphanumeric"}},{"operator":"OR","attribute":{"name":"city","operator":"=","value":"mojokerto","type":"alphanumeric"}},{"operator":"OR","conditions":[{"attribute":{"name":"warehouse_id","operator":"=","value":"1","type":"numeric"}},{"operator":"AND","attribute":{"name":"warehouse_detail_id","operator":"=","value":"2","type":"numeric"}}]}]}]}`,
			wantErr: false,
		},
		{
			name: "Normal case",
			args: args{
				query: `
                  id = 1
                  && member_id = 2
                  && user_id = 3
                  && (
                    province = jatim
                    || city = mojokerto
                    || (
                      warehouse_id = 1
                      && warehouse_detail_id = 2
                    )
                  )
				  && data_id = 54
`,
			},
			want:    `{"conditions":[{"attribute":{"name":"id","operator":"=","value":"1","type":"numeric"}},{"operator":"AND","attribute":{"name":"member_id","operator":"=","value":"2","type":"numeric"}},{"operator":"AND","attribute":{"name":"user_id","operator":"=","value":"3","type":"numeric"}},{"operator":"AND","conditions":[{"attribute":{"name":"province","operator":"=","value":"jatim","type":"alphanumeric"}},{"operator":"OR","attribute":{"name":"city","operator":"=","value":"mojokerto","type":"alphanumeric"}},{"operator":"OR","conditions":[{"attribute":{"name":"warehouse_id","operator":"=","value":"1","type":"numeric"}},{"operator":"AND","attribute":{"name":"warehouse_detail_id","operator":"=","value":"2","type":"numeric"}}]}]},{"operator":"AND","attribute":{"name":"data_id","operator":"=","value":"54","type":"numeric"}}]}`,
			wantErr: false,
		},
		{
			name: "Normal case",
			args: args{
				query: "((date<=2019-09-09 && date > 2019-08-08) || (p_date>=2019-01-01 && p_date<2019-02-02)) && (member_type=1||member_type=2)",
			},
			want:    `{"conditions":[{"conditions":[{"conditions":[{"attribute":{"name":"date","operator":"\u003c=","value":"2019-09-09","type":"alphanumeric"}},{"operator":"AND","attribute":{"name":"date","operator":"\u003e","value":"2019-08-08","type":"alphanumeric"}}]},{"operator":"OR","conditions":[{"attribute":{"name":"p_date","operator":"\u003e=","value":"2019-01-01","type":"alphanumeric"}},{"operator":"AND","attribute":{"name":"p_date","operator":"\u003c","value":"2019-02-02","type":"alphanumeric"}}]}]},{"operator":"AND","conditions":[{"attribute":{"name":"member_type","operator":"=","value":"1","type":"numeric"}},{"operator":"OR","attribute":{"name":"member_type","operator":"=","value":"2","type":"numeric"}}]}]}`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateCondition(tt.args.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateCondition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			byteBuf, _ := json.Marshal(got)
			if !strings.EqualFold(string(byteBuf), tt.want) {
				t.Errorf("GenerateCondition() = %v, want %v", string(byteBuf), tt.want)
			}
		})
	}
}

func TestCondition_ValidateStruct(t *testing.T) {
	type args struct {
		query  string
		object interface{}
	}
	tests := []struct {
		name        string
		args        args
		wantIsValid bool
		wantErr     bool
	}{
		{
			name: "Normal case - struct validation",
			args: args{
				query: `(ID=1 && (MemberID=12||MemberID=2))  &&   (Division=engineering || Division=finance)`,
				object: struct {
					ID       string `json:"id"`
					MemberID string `json:"member_id"`
					Division string `json:"division"`
				}{
					ID:       "1",
					MemberID: "2",
					Division: "finance",
				},
			},
			wantIsValid: true,
			wantErr:     false,
		},
		{
			name: "Normal case - struct validation",
			args: args{
				query: `(ID=1 &&  MemberID=2  &&   (Division=engineering || Division=finance))||(MemberID=3)`,
				object: struct {
					ID       int    `json:"id"`
					MemberID int    `json:"member_id"`
					Division string `json:"division"`
				}{
					ID:       1,
					MemberID: 3,
					Division: "finance",
				},
			},
			wantIsValid: true,
			wantErr:     false,
		},
		{
			name: "Normal case - struct validation - brand attribute is not exist",
			args: args{
				query: `(id=1 &&  member_id=2  &&   (division=engineering || division=finance))||(member_id=3&&brand=abc)`,
				object: struct {
					ID       int    `json:"id"`
					MemberID int    `json:"member_id"`
					Division string `json:"division"`
				}{
					ID:       1,
					MemberID: 3,
					Division: "finance",
				},
			},
			wantIsValid: false,
			wantErr:     false,
		},
		{
			name: "Normal case - struct validation",
			args: args{
				query: `ID=1 &&  MemberID=3  && ((Division=engineering || Division=finance || Division=people)&&(MemberID=2||ID=1))`,
				object: struct {
					ID       int    `json:"id"`
					MemberID int    `json:"member_id"`
					Division string `json:"division"`
				}{
					ID:       1,
					MemberID: 3,
					Division: "people",
				},
			},
			wantIsValid: true,
			wantErr:     false,
		},
		{
			name: "Normal case - struct validation",
			args: args{
				query: `ID=1 &&  MemberID=2  &&   (Division=engineering || Division=finance)`,
				object: struct {
					ID       int
					MemberID int
					Division string
				}{
					ID:       1,
					MemberID: 2,
					Division: "engineering",
				},
			},
			wantIsValid: true,
			wantErr:     false,
		},
		{
			name: "Normal case - struct validation - Brand attribute is not exist",
			args: args{
				query: `ID=1 &&  MemberID=2  &&   (Division=engineering || Division=finance) && Brand=Adidas`,
				object: struct {
					ID       int
					MemberID int
					Division string
				}{
					ID:       1,
					MemberID: 2,
					Division: "engineering",
				},
			},
			wantIsValid: false,
			wantErr:     false,
		},
		{
			name: "Normal case - struct validation - skip not exist attribute because using OR condition",
			args: args{
				query: `ID=1 &&  MemberID=2  &&   (Division=engineering || Division=finance) && (Category=Bawahan || ID=1 || Brand=nike)`,
				object: struct {
					ID       int
					MemberID int
					Division string
				}{
					ID:       1,
					MemberID: 2,
					Division: "engineering",
				},
			},
			wantIsValid: true,
			wantErr:     false,
		},
		{
			name: "Error case",
			args: args{
				query:  `ID=1 &&  MemberID=2  &&   (Division=engineering || Division=finance)`,
				object: nil,
			},
			wantIsValid: false,
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			proc := NewProcessor().RegisterCondition(tt.args.query)
			gotIsValid, err := proc.ValidateStruct(tt.args.object)
			if (err != nil) != tt.wantErr {
				t.Errorf("Condition.ValidateStruct() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotIsValid != tt.wantIsValid {
				t.Errorf("Condition.ValidateStruct() = %v, want %v", gotIsValid, tt.wantIsValid)
			}
		})
	}
}

func TestCondition_ValidateMultipleStructs(t *testing.T) {
	type fields struct {
		Operator   string
		Attribute  *structs2.Attribute
		Conditions []*structs2.Condition
	}
	type args struct {
		query string
		data  interface{}
	}
	type firstStruct struct {
		ID       string `json:"id"`
		MemberID string `json:"member_id"`
		Division string `json:"division"`
	}
	type secondStruct struct {
		Name string `json:"name"`
	}
	type thirdStruct struct {
		Type    string `json:"type"`
		Segment string `json:"segment"`
	}

	thirdData := thirdStruct{
		Type:    "ABC",
		Segment: "new-member",
	}

	tests := []struct {
		name        string
		fields      fields
		args        args
		wantIsValid bool
		wantErr     bool
	}{
		{
			name: "Normal case - one struct validation",
			args: args{
				query: `MemberID=345`,
				data: []interface{}{
					firstStruct{
						ID:       "123",
						MemberID: "345",
						Division: "engineering",
					},
				},
			},
			wantIsValid: true,
			wantErr:     false,
		},
		{
			name: "Normal case - one struct validation",
			args: args{
				query: `Division=engineering`,
				data: []interface{}{
					firstStruct{
						ID:       "123",
						MemberID: "345",
						Division: "engineering",
					},
				},
			},
			wantIsValid: true,
			wantErr:     false,
		},
		{
			name: "Normal case - one struct validation - not equal",
			args: args{
				query: `MemberID!=1232323`,
				data: []interface{}{
					firstStruct{
						ID:       "123",
						MemberID: "345",
						Division: "engineering",
					},
				},
			},
			wantIsValid: true,
			wantErr:     false,
		},
		{
			name: "Normal case - one struct validation - not equal",
			args: args{
				query: `Division!=sales`,
				data: []interface{}{
					firstStruct{
						ID:       "123",
						MemberID: "345",
						Division: "engineering",
					},
				},
			},
			wantIsValid: true,
			wantErr:     false,
		},
		{
			name: "Normal case - one struct validation - contains str",
			args: args{
				query: `Division|=eng`,
				data: []interface{}{
					firstStruct{
						ID:       "123",
						MemberID: "345",
						Division: "engineering",
					},
				},
			},
			wantIsValid: true,
			wantErr:     false,
		},
		{
			name: "Normal case - one struct validation - not contains str",
			args: args{
				query: `Division|=eng`,
				data: []interface{}{
					firstStruct{
						ID:       "123",
						MemberID: "345",
						Division: "pasukan katak bersaudara 7896",
					},
				},
			},
			wantIsValid: false,
			wantErr:     false,
		},
		{
			name: "Normal case - one struct validation - contains regex",
			args: args{
				query: `Division |~     "katak[\s][a-z]+[\s][0-9]+"`,
				data: []interface{}{
					firstStruct{
						ID:       "123",
						MemberID: "345",
						Division: "pasukan katak bersaudara 7896",
					},
				},
			},
			wantIsValid: true,
			wantErr:     false,
		},
		{
			name: "Normal case - one struct validation - attribute brand not exist",
			args: args{
				query: `MemberID=345 && Brand=adidas`,
				data: []interface{}{
					firstStruct{
						ID:       "123",
						MemberID: "345",
						Division: "engineering",
					},
				},
			},
			wantIsValid: false,
			wantErr:     false,
		},
		{
			name: "Normal case - multi struct validation - all attributes exist",
			args: args{
				query: `Type=ABC && Name=Test`,
				data: []interface{}{
					thirdData,
					secondStruct{
						Name: "Test",
					},
				},
			},
			wantIsValid: true,
			wantErr:     false,
		},
		{
			name: "Normal case - multi struct validation - attribute memberId in secondStruct not exist",
			args: args{
				query: `Name=Test && MemberID=1010101`,
				data: []interface{}{
					thirdData,
					secondStruct{
						Name: "Test",
					},
				},
			},
			wantIsValid: false,
			wantErr:     false,
		},
		{
			name: "Normal case - multi struct validation",
			args: args{
				query: `ID=123 && Name=Test`,
				data: []interface{}{
					firstStruct{
						ID:       "123",
						MemberID: "345",
						Division: "engineering",
					},
					secondStruct{
						Name: "Test",
					},
				},
			},
			wantIsValid: true,
			wantErr:     false,
		},
		{
			name: "Normal case - struct validation",
			args: args{
				query: `ID=123`,
				data: firstStruct{
					ID:       "123",
					MemberID: "345",
					Division: "engineering",
				},
			},
			wantIsValid: true,
			wantErr:     false,
		},
		{
			name: "Normal case - multi struct validation",
			args: args{
				query: `(ID=1234 || secondStruct=Test || Segment=new-member) && (MemberID=345 && Name=Test) && Type=ABC`,
				data: []interface{}{
					firstStruct{
						ID:       "123",
						MemberID: "345",
						Division: "engineering",
					},
					secondStruct{
						Name: "Test",
					},
					thirdData,
				},
			},
			wantIsValid: true,
			wantErr:     false,
		},
		{
			name: "Error case",
			args: args{
				query: `id=1`,
				data:  []interface{}{},
			},
			wantIsValid: false,
			wantErr:     true,
		},
		{
			name: "Error case",
			args: args{
				query: `id=1`,
				data:  nil,
			},
			wantIsValid: false,
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			proc := NewProcessor().RegisterCondition(tt.args.query)
			gotIsValid, err := proc.ValidateMultipleStructs(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Condition.ValidateMultipleStructs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotIsValid != tt.wantIsValid {
				t.Errorf("Condition.ValidateMultipleStructs() = %v, want %v", gotIsValid, tt.wantIsValid)
			}
		})
	}
}

func TestCondition_ValidateCondition(t *testing.T) {
	tests := []struct {
		name           string
		referenceQuery string
		input          string
		wantIsValid    bool
		wantErr        bool
	}{
		{
			name:           "Normal case",
			referenceQuery: "id=1 && member_id=45",
			input:          "id=2 || member_id=45",
			wantIsValid:    false,
			wantErr:        false,
		},
		{
			name:           "Normal case - ignore case",
			referenceQuery: "name=Budi && brand=Arava && member_id=45",
			input:          "name=budi && brand=arava && member_id=45",
			wantIsValid:    true,
			wantErr:        false,
		},
		{
			name:           "Normal case",
			referenceQuery: "id=1 && member_id=45",
			input:          "(id=2||id=1) && member_id=45",
			wantIsValid:    true,
			wantErr:        false,
		},
		{
			name:           "Normal case",
			referenceQuery: "id=1 && member_id=45",
			input:          "id=1",
			wantIsValid:    false,
			wantErr:        false,
		},
		{
			name:           "Normal case",
			referenceQuery: "id=1 || member_id=45",
			input:          "id=1",
			wantIsValid:    true,
			wantErr:        false,
		},
		{
			name:           "Normal case",
			referenceQuery: "id=1 || member_id=45",
			input:          "member_id=45",
			wantIsValid:    true,
			wantErr:        false,
		},
		{
			name:           "Normal case",
			referenceQuery: "id=1 && member_id=45",
			input:          "id=1 && member_id=45",
			wantIsValid:    true,
			wantErr:        false,
		},
		{
			name:           "Normal case",
			referenceQuery: "id=1 && member_id=45",
			input:          "id=1 && (member_id=23||member_id=35)",
			wantIsValid:    false,
			wantErr:        false,
		},
		{
			name:           "Normal case",
			referenceQuery: "id=1 && member_id=45",
			input:          "id=1 && (member_id=23||member_id=45)",
			wantIsValid:    true,
			wantErr:        false,
		},
		{
			name:           "Normal case",
			referenceQuery: "id=1 && member_id=45",
			input:          "id=1 && member_id=44",
			wantIsValid:    false,
			wantErr:        false,
		},
		{
			name:           "Normal case",
			referenceQuery: "id=1 || member_id=45",
			input:          "id=1 && member_id=22",
			wantIsValid:    true,
			wantErr:        false,
		},
		{
			name:           "Normal case - condition group",
			referenceQuery: "(id=1 || id=2) && member_id=45",
			input:          "id=1 && member_id=45",
			wantIsValid:    true,
			wantErr:        false,
		},
		{
			name:           "Normal case - multi condition group",
			referenceQuery: "id=1 &&  member_id=3  && ((division=engineering || division=finance || division=people)&&(member_id=2||id=1))",
			input:          "id=1 && member_id=3 && division=engineering",
			wantIsValid:    true,
			wantErr:        false,
		},
		{
			name:           "Normal case - multi condition group",
			referenceQuery: "id=1 &&  member_id=3  && ((division=engineering || division=finance || division=people)&&(member_id=2||id=1))",
			input:          "(id=1 && member_id=3) && (division=tech&&division=finance)",
			wantIsValid:    false,
			wantErr:        false,
		},
		{
			name:           "Normal case - multi condition group",
			referenceQuery: "id=1 &&  member_id=3  && ((division=engineering || division=finance || division=people)&&(member_id=2||id=1))",
			input:          "((id=1 && member_id=3) && (division=tech||division=finance))",
			wantIsValid:    true,
			wantErr:        false,
		},
		{
			name:           "Normal case - multi condition group",
			referenceQuery: "(id=1 || id=2) && (member_id=45||member_id=10)",
			input:          "id=1 && (member_id=10)",
			wantIsValid:    true,
			wantErr:        false,
		},
		{
			name:           "Normal case - multi condition group",
			referenceQuery: "(id=1 || id=2) && (member_id=45||member_id=10)",
			input:          "id=3 && member_id=10 || id=14",
			wantIsValid:    false,
			wantErr:        false,
		},
		{
			name:           "Normal case - multi condition group",
			referenceQuery: "(id=1 || id=2) && (member_id=45||member_id=10) && (segment=trial||segment=free)",
			input:          "id=1 && member_id=10 && segment=free",
			wantIsValid:    true,
			wantErr:        false,
		},
		{
			name:           "Normal case - match only one input",
			referenceQuery: "id=1",
			input:          "id=1 && member_id=10 && segment=free",
			wantIsValid:    true,
			wantErr:        false,
		},
		{
			name:           "Normal case - string condition",
			referenceQuery: "deviceType=mobile || memberId=xxx",
			input:          "deviceType=mobile",
			wantIsValid:    true,
			wantErr:        false,
		},
		{
			name:           "Normal case - string condition",
			referenceQuery: "deviceType=mobile && memberId=xxx",
			input:          "deviceType=mobile",
			wantIsValid:    false,
			wantErr:        false,
		},
		{
			name:           "Normal case - rule engine case",
			referenceQuery: "deviceType=mobile && ABTest=xxx ",
			input:          "deviceType=mobile && ABTest=yyy",
			wantIsValid:    false,
			wantErr:        false,
		},
		{
			name:           "Normal case - rule engine case with group condition",
			referenceQuery: "deviceType=mobile && ABTest=xxx ",
			input:          "deviceType=mobile && (ABTest=yyy||ABTest=xxx)",
			wantIsValid:    true,
			wantErr:        false,
		},
		{
			name:           "Normal case - greater than operator - integer",
			referenceQuery: "(id=1 || id=2) && member_id>100",
			input:          "id=1 && member_id=111",
			wantIsValid:    true,
			wantErr:        false,
		},
		{
			name:           "Normal case - greater than operator - integer",
			referenceQuery: "(id=1 || id=2) && member_id>100",
			input:          "id=1 && member_id=111",
			wantIsValid:    true,
			wantErr:        false,
		},
		{
			name:           "Normal case - greater than equal operator - datetime",
			referenceQuery: `(id=1 || id=2) && create_date>="2020-02-02 12:12:12"`,
			input:          `id=1 && create_date="2020-02-02 12:12:12"`,
			wantIsValid:    true,
			wantErr:        false,
		},
		{
			name:           "Normal case - greater than operator - float",
			referenceQuery: "(id=1 || id=2) && price>1200.50 && (segment=hijaber||segment=girl||segment=cantik) && poin>100",
			input:          "id=1 && price=1200.51 && ((segment=cantik&&poin=58)||(segment=girl&&poin=518))",
			wantIsValid:    true,
			wantErr:        false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			proc := NewProcessor().RegisterCondition(tt.referenceQuery)

			inputCondition, err := GenerateCondition(tt.input)
			if err != nil {
				t.Errorf("Condition.ValidateCondition() input error = %v", err)
				return
			}
			gotIsValid, err := proc.ValidateCondition(inputCondition)
			if (err != nil) != tt.wantErr {
				t.Errorf("Condition.ValidateCondition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotIsValid != tt.wantIsValid {
				t.Errorf("Condition.ValidateCondition() = %v, want %v", gotIsValid, tt.wantIsValid)
			}
		})
	}
}

func TestCondition_FilterSlice(t *testing.T) {
	type Account struct {
		ID        int        `json:"id"`
		MemberID  int        `json:"member_id"`
		Division  string     `json:"division"`
		Score     *int       `json:"score"`
		Point     *int64     `json:"point"`
		Wallet    *float32   `json:"wallet"`
		Money     *float64   `json:"money"`
		JoinDate  time.Time  `json:"join_date"`
		LeaveDate *time.Time `json:"leave_date"`
	}

	fInt := func(i int) *int {
		return &i
	}
	fInt64 := func(i int64) *int64 {
		return &i
	}
	fFloat := func(f float32) *float32 {
		return &f
	}
	fFloat64 := func(f float64) *float64 {
		return &f
	}
	fTime := func(t time.Time) *time.Time {
		return &t
	}

	testData := []Account{
		{
			ID:        1,
			MemberID:  21,
			Division:  "people",
			Score:     fInt(90),
			Point:     fInt64(12000),
			Wallet:    fFloat(100000),
			Money:     fFloat64(10000),
			JoinDate:  time.Date(2020, 3, 9, 0, 0, 0, 0, time.UTC),
			LeaveDate: fTime(time.Date(2020, 12, 9, 0, 0, 0, 0, time.UTC)),
		},
		{
			ID:        2,
			MemberID:  22,
			Division:  "finance",
			Score:     fInt(40),
			Point:     fInt64(1000),
			Wallet:    fFloat(1000),
			Money:     fFloat64(50000),
			JoinDate:  time.Date(2014, 1, 9, 0, 0, 0, 0, time.UTC),
			LeaveDate: fTime(time.Date(2015, 12, 9, 0, 0, 0, 0, time.UTC)),
		},
		{
			ID:        3,
			MemberID:  23,
			Division:  "business",
			Score:     fInt(60),
			Point:     fInt64(5000),
			Wallet:    fFloat(5000),
			Money:     fFloat64(80000),
			JoinDate:  time.Date(2016, 12, 9, 0, 0, 0, 0, time.UTC),
			LeaveDate: fTime(time.Date(2017, 12, 9, 0, 0, 0, 0, time.UTC)),
		},
		{
			ID:        4,
			MemberID:  24,
			Division:  "managerial",
			Score:     fInt(70),
			Point:     fInt64(20000),
			Wallet:    fFloat(4000),
			Money:     fFloat64(900000),
			JoinDate:  time.Date(2018, 4, 9, 0, 0, 0, 0, time.UTC),
			LeaveDate: fTime(time.Date(2019, 12, 9, 0, 0, 0, 0, time.UTC)),
		},
		{
			ID:        5,
			MemberID:  25,
			Division:  "engineering",
			Score:     fInt(100),
			Point:     fInt64(3000),
			Wallet:    fFloat(100),
			Money:     fFloat64(1500000),
			JoinDate:  time.Date(2015, 10, 9, 0, 0, 0, 0, time.UTC),
			LeaveDate: nil,
		},
		{
			ID:        5,
			MemberID:  25,
			Division:  "engineering",
			Score:     nil,
			Point:     nil,
			Wallet:    nil,
			Money:     nil,
			JoinDate:  time.Date(2015, 7, 9, 0, 0, 0, 0, time.UTC),
			LeaveDate: fTime(time.Date(2016, 12, 9, 0, 0, 0, 0, time.UTC)),
		},
	}

	type args struct {
		query   string
		objects interface{}
	}
	tests := []struct {
		name        string
		args        args
		wantResults interface{}
		wantErr     bool
	}{
		{
			name: "Normal case",
			args: args{
				query:   "ID=1||ID=2",
				objects: testData,
			},
			wantResults: []Account{
				{
					ID:        1,
					MemberID:  21,
					Division:  "people",
					Score:     fInt(90),
					Point:     fInt64(12000),
					Wallet:    fFloat(100000),
					Money:     fFloat64(10000),
					JoinDate:  time.Date(2020, 3, 9, 0, 0, 0, 0, time.UTC),
					LeaveDate: fTime(time.Date(2020, 12, 9, 0, 0, 0, 0, time.UTC)),
				},
				{
					ID:        2,
					MemberID:  22,
					Division:  "finance",
					Score:     fInt(40),
					Point:     fInt64(1000),
					Wallet:    fFloat(1000),
					Money:     fFloat64(50000),
					JoinDate:  time.Date(2014, 1, 9, 0, 0, 0, 0, time.UTC),
					LeaveDate: fTime(time.Date(2015, 12, 9, 0, 0, 0, 0, time.UTC)),
				},
			},
			wantErr: false,
		},
		{
			name: "Normal case",
			args: args{
				query:   "(MemberID=23)||((ID=1 && MemberID=21)&&(Division=people||Division=managerial))",
				objects: testData,
			},
			wantResults: []Account{
				{
					ID:        1,
					MemberID:  21,
					Division:  "people",
					Score:     fInt(90),
					Point:     fInt64(12000),
					Wallet:    fFloat(100000),
					Money:     fFloat64(10000),
					JoinDate:  time.Date(2020, 3, 9, 0, 0, 0, 0, time.UTC),
					LeaveDate: fTime(time.Date(2020, 12, 9, 0, 0, 0, 0, time.UTC)),
				},
				{
					ID:        3,
					MemberID:  23,
					Division:  "business",
					Score:     fInt(60),
					Point:     fInt64(5000),
					Wallet:    fFloat(5000),
					Money:     fFloat64(80000),
					JoinDate:  time.Date(2016, 12, 9, 0, 0, 0, 0, time.UTC),
					LeaveDate: fTime(time.Date(2017, 12, 9, 0, 0, 0, 0, time.UTC)),
				},
			},
			wantErr: false,
		},
		{
			name: "Normal case",
			args: args{
				query:   "MemberID=21 && Score=90",
				objects: testData,
			},
			wantResults: []Account{
				{
					ID:        1,
					MemberID:  21,
					Division:  "people",
					Score:     fInt(90),
					Point:     fInt64(12000),
					Wallet:    fFloat(100000),
					Money:     fFloat64(10000),
					JoinDate:  time.Date(2020, 3, 9, 0, 0, 0, 0, time.UTC),
					LeaveDate: fTime(time.Date(2020, 12, 9, 0, 0, 0, 0, time.UTC)),
				},
			},
			wantErr: false,
		},
		{
			name: "Normal case",
			args: args{
				query:   "MemberID=25 && Point>=3000",
				objects: testData,
			},
			wantResults: []Account{
				{
					ID:        5,
					MemberID:  25,
					Division:  "engineering",
					Score:     fInt(100),
					Point:     fInt64(3000),
					Wallet:    fFloat(100),
					Money:     fFloat64(1500000),
					JoinDate:  time.Date(2015, 10, 9, 0, 0, 0, 0, time.UTC),
					LeaveDate: nil,
				},
			},
			wantErr: false,
		},
		{
			name: "Normal case",
			args: args{
				query:   `JoinDate>"2015-01-01T00:00:00+07:00" && JoinDate<="2016-01-01T00:00:00+07:00" && Score>80 && Point<4000 && Wallet>90 && Money>=1500000`,
				objects: testData,
			},
			wantResults: []Account{
				{
					ID:        5,
					MemberID:  25,
					Division:  "engineering",
					Score:     fInt(100),
					Point:     fInt64(3000),
					Wallet:    fFloat(100),
					Money:     fFloat64(1500000),
					JoinDate:  time.Date(2015, 10, 9, 0, 0, 0, 0, time.UTC),
					LeaveDate: nil,
				},
			},
			wantErr: false,
		},
		{
			name: "Normal case - empty",
			args: args{
				query:   "(member_id=25)||((id=1 && member_id=21)&&(division=people||division=managerial))",
				objects: []Account{},
			},
			wantResults: []Account{},
			wantErr:     false,
		},
		{
			name: "Error case - nil object",
			args: args{
				query:   "(member_id=25)||((id=1 && member_id=21)&&(division=people||division=managerial))",
				objects: nil,
			},
			wantResults: nil,
			wantErr:     true,
		},
		{
			name: "Error case - invalid type",
			args: args{
				query:   "(member_id=25)||((id=1 && member_id=21)&&(division=people||division=managerial))",
				objects: Account{},
			},
			wantResults: nil,
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			proc := NewProcessor().RegisterCondition(tt.args.query)
			gotResults, err := proc.FilterSlice(tt.args.objects)
			if (err != nil) != tt.wantErr {
				t.Errorf("Condition.FilterSlice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResults, tt.wantResults) {
				t.Errorf("Condition.FilterSlice() = %v, want %v", gotResults, tt.wantResults)
			}
		})
	}
}
