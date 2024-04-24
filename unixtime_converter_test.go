package freak_conventer

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type TestStruct struct {
	IntField    int64     `json:"intField"`
	StringField string    `json:"stringField"`
	TimeField   time.Time `json:"timeField"`
}

type TestStructWithOmitempty struct {
	IntField    int64     `json:"intField"`
	StringField string    `json:"stringField,omitempty"`
	TimeField   time.Time `json:"timeField"`
}

type TestStructWithNestedStruct struct {
	IntField            int64                   `json:"intField"`
	StringField         string                  `json:"stringField"`
	NestedField         TestStruct              `json:"nestedField"`
	NestedFieldWithOmni TestStructWithOmitempty `json:"nestedOmni"`
}

type TestStructWithArray struct {
	IntField      int64                        `json:"intField"`
	Array         []TestStruct                 `json:"arrayField"`
	ArrayWithOmni []TestStructWithNestedStruct `json:"arrayOmni"`
}

type AnonStructField struct {
	IntField  int64     `json:"intField"`
	TimeField time.Time `json:"timeField"`
}
type TestStructWithAnonStruct struct {
	IntField int64 `json:"intField"`
	AnonStructField
}

func TestConvertToUnixTime(t *testing.T) {
	testCases := []struct {
		name        string
		input       interface{}
		expected    map[string]interface{}
		expectedErr error
	}{
		{
			name: "valid struct",
			input: TestStruct{
				IntField:    1,
				StringField: "TestStruct",
				TimeField:   time.Unix(1633024900, 0),
			},
			expected: map[string]interface{}{
				"intField":    int64(1),
				"stringField": "TestStruct",
				"timeField":   int64(1633024900),
			},
			expectedErr: nil,
		},
		{
			name: "valid struct with omit empty tag",
			input: TestStructWithOmitempty{
				IntField:  1,
				TimeField: time.Unix(1633024862, 0),
			},
			expected: map[string]interface{}{
				"intField":  int64(1),
				"timeField": int64(1633024862),
			},
			expectedErr: nil,
		},
		{
			name: "valid struct with nested struct",
			input: TestStructWithNestedStruct{
				IntField:    1,
				StringField: "TestStructWithNestedStruct",
				NestedField: TestStruct{
					IntField:    2,
					StringField: "TestStruct",
					TimeField:   time.Unix(1633024900, 0),
				},
				NestedFieldWithOmni: TestStructWithOmitempty{
					IntField:  3,
					TimeField: time.Unix(1633025000, 0),
				},
			},
			expected: map[string]interface{}{
				"intField":    int64(1),
				"stringField": "TestStructWithNestedStruct",
				"nestedField": map[string]interface{}{
					"intField":    int64(2),
					"stringField": "TestStruct",
					"timeField":   int64(1633024900),
				},
				"nestedOmni": map[string]interface{}{
					"intField":  int64(3),
					"timeField": int64(1633025000),
				},
			},
			expectedErr: nil,
		},
		{
			name: "valid struct with array",
			input: TestStructWithArray{
				IntField: 1,
				Array: []TestStruct{
					{
						IntField:    2,
						StringField: "TestStruct1",
						TimeField:   time.Unix(1633024862, 0),
					},
					{
						IntField:    3,
						StringField: "TestStruct2",
						TimeField:   time.Unix(1633024864, 0),
					},
				},
				ArrayWithOmni: []TestStructWithNestedStruct{
					{
						IntField:    4,
						StringField: "TestStructWithNestedStruct1",
						NestedField: TestStruct{
							IntField:    5,
							StringField: "TestStruct1",
							TimeField:   time.Unix(1633024866, 0),
						},
						NestedFieldWithOmni: TestStructWithOmitempty{
							IntField:  6,
							TimeField: time.Unix(1633024867, 0),
						},
					},
					{
						IntField:    7,
						StringField: "TestStructWithNestedStruct2",
						NestedField: TestStruct{
							IntField:    8,
							StringField: "TestStruct2",
							TimeField:   time.Unix(1633024869, 0),
						},
						NestedFieldWithOmni: TestStructWithOmitempty{
							IntField:  9,
							TimeField: time.Unix(1633024870, 0),
						},
					},
				},
			},
			expected: map[string]interface{}{
				"intField": int64(1),
				"arrayField": []interface{}{
					map[string]interface{}{
						"intField":    int64(2),
						"stringField": "TestStruct1",
						"timeField":   int64(1633024862),
					},
					map[string]interface{}{
						"intField":    int64(3),
						"stringField": "TestStruct2",
						"timeField":   int64(1633024864),
					},
				},
				"arrayOmni": []interface{}{
					map[string]interface{}{
						"intField":    int64(4),
						"stringField": "TestStructWithNestedStruct1",
						"nestedField": map[string]interface{}{
							"intField":    int64(5),
							"stringField": "TestStruct1",
							"timeField":   int64(1633024866),
						},
						"nestedOmni": map[string]interface{}{
							"intField":  int64(6),
							"timeField": int64(1633024867),
						},
					},
					map[string]interface{}{
						"intField":    int64(7),
						"stringField": "TestStructWithNestedStruct2",
						"nestedField": map[string]interface{}{
							"intField":    int64(8),
							"stringField": "TestStruct2",
							"timeField":   int64(1633024869),
						},
						"nestedOmni": map[string]interface{}{
							"intField":  int64(9),
							"timeField": int64(1633024870),
						},
					},
				},
			},
			expectedErr: nil,
		},
		{
			name: "valid struct with anonymous struct",
			input: TestStructWithAnonStruct{
				IntField: 1,
				AnonStructField: AnonStructField{
					IntField:  2,
					TimeField: time.Unix(1633024862, 0),
				},
			},
			expected: map[string]interface{}{
				"intField":  int64(2),
				"timeField": int64(1633024862),
			},
			expectedErr: nil,
		},
		{
			name:        "nil input",
			input:       nil,
			expected:    nil,
			expectedErr: fmt.Errorf("input data is nil"),
		},
		{
			name:        "zero input",
			input:       TestStruct{},
			expected:    nil,
			expectedErr: fmt.Errorf("input data is zero"),
		},
		{
			name:        "non struct input",
			input:       123,
			expected:    nil,
			expectedErr: fmt.Errorf("input data should be a struct"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := ConvertToUnixTime(tc.input)
			assert.Equal(t, tc.expected, result)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}
