package main

import (
	"fmt"
	"testing"
)

func TestGetJsonMap(t *testing.T) {
	unittest := "Test getJsonMap(input []byte) (*map[string]interface{}, error)"
	tests := []struct {
		input string
		want  *map[string]interface{}
	}{
		{
			"{\"test1\":true,\"test2\":[1,2,3,\"4\",false],\"test3\":{\"test4\":1}}",
			&map[string]interface{}{
				"test1": true,
				"test2": []interface{}{1, 2, 3, "4", false},
				"test3": map[string]interface{}{
					"test4": 1,
				},
			},
		},
		{
			"",
			nil,
		},
		{
			"{\"test1\":true,\"test2\":[1,2,3,\"4\",false],\"test3\":{\"test4\":1}",
			nil,
		},
	}

	for i, tt := range tests {
		testname := fmt.Sprintf("%v, case %v", unittest, i)

		t.Run(testname, func(t *testing.T) {
			got, err := getJsonMap([]byte(tt.input))

			if err != nil && got != nil {
				t.Errorf("%v:\nGot %v, want %v", testname, err, tt.want)
			} else if got != nil && fmt.Sprint(got) != fmt.Sprint(tt.want) {
				t.Errorf("%v:\nGot %v, want %v", testname, got, tt.want)
			}
		})
	}
}

func TestCapFirstLetter(t *testing.T) {
	unittest := "Test capFirstLetter(s string) string"
	tests := []struct {
		input string
		want  string
	}{
		{"test", "Test"},
		{"test", "Test"},
		{"test", "Test"},
	}

	for i, tt := range tests {
		testname := fmt.Sprintf("%v, case %v", unittest, i)

		t.Run(testname, func(t *testing.T) {
			got := capFirstLetter(tt.input)

			if got != tt.want {
				t.Errorf("%v:\nGot %v, want %v", testname, got, tt.want)
			}
		})
	}
}

func TestIsChar(t *testing.T) {
	unittest := "Test isChar(s byte) bool"
	tests := []struct {
		input byte
		want  bool
	}{
		{'a', true},
		{'A', true},
		{'z', true},
		{'Z', true},
		{'0', false},
		{'!', false},
	}

	for i, tt := range tests {
		testname := fmt.Sprintf("%v, case %v", unittest, i)

		t.Run(testname, func(t *testing.T) {
			got := isChar(tt.input)

			if got != tt.want {
				t.Errorf("%v:\nGot %v, want %v", testname, got, tt.want)
			}
		})
	}
}

func TestBuildNVPair(t *testing.T) {
	unittest := "Test buildNVPair(name, ftype string) string"
	tests := []struct {
		input struct {
			name, ftype string
		}
		want string
	}{
		{struct {
			name  string
			ftype string
		}{"test", "string"}, "Test string `json:\"test\"`"},
		{struct {
			name  string
			ftype string
		}{"", "string"}, ""},
		{struct {
			name  string
			ftype string
		}{"test", ""}, ""},
		{struct {
			name  string
			ftype string
		}{"1test", "int"}, "Go1test int `json:\"1test\"`"},
	}

	for i, tt := range tests {
		testname := fmt.Sprintf("%v, case %v", unittest, i)

		t.Run(testname, func(t *testing.T) {
			got, err := buildNVPair(tt.input.name, tt.input.ftype)

			if err != nil && got != "" {
				t.Errorf("%v:\nGot %v, want %v", testname, err, tt.want)
			} else if got != tt.want {
				t.Errorf("%v:\nGot %v, want %v", testname, got, tt.want)
			}
		})
	}
}

func TestGetFileNameNoPathExtGo(t *testing.T) {
	unittest := "func getFileNameNoPathExtGo(file string) string"
	tests := []struct {
		input, want string
	}{
		{"test.json", "test"},
		{`/path/to/test.json`, "test"},
		{`path/to/test.json`, "test"},
	}

	for i, tt := range tests {
		testname := fmt.Sprintf("%v, case %v", unittest, i)

		t.Run(testname, func(t *testing.T) {
			got := getFileNameNoPathExt(tt.input)

			if got != tt.want {
				t.Errorf("%v:\nGot %v, want %v", testname, got, tt.want)
			}
		})
	}
}

func TestIsJSONFile(t *testing.T) {
	unittest := "func isJSONFile(file string) bool"
	tests := []struct {
		input string
		want  bool
	}{
		{"test.json", true},
		{"test", false},
		{"test.go", false},
	}

	for i, tt := range tests {
		testname := fmt.Sprintf("%v, case %v", unittest, i)

		t.Run(testname, func(t *testing.T) {
			got := isJSONFile(tt.input)

			if got != tt.want {
				t.Errorf("%v:\nGot %v, want %v", testname, got, tt.want)
			}
		})
	}
}

// TODO: need ordered map
// func TestBuildUnnamedStruct(t *testing.T) {
// 	unittest := "Test buildNamedStruct(m *map[string]interface{}, name string) (string, error)"
// 	test := []byte("{\"test\":\"1\",\"test1\":true,\"test2\":[1,2,3,\"4\",false],\"test3\":{\"test4\":1}}")
// 	want := "struct {\nTest string `json:\"test\"\nTest1 bool `json:\"test1\"`\nTest2 []interface{} `json:\"test2\"`\nTest3 struct {\nTest4 float64 `json:\"test4\"`\n} `json:\"test3\"`\n}"
// 	res, err := getJsonMap(test)

// 	if err != nil {
// 		t.Errorf("Error creating json map: %v", err)
// 	}

// 	got, err := buildUnnamedStruct(res)

// 	if got != want {
// 		t.Errorf("%v:\nGot %v, want %v", unittest, got, want)
// 	}
// }
