package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const (
	STRING = "string"
	INT    = "int"
	FLOAT  = "float64"
	BOOL   = "bool"
	ARRAY  = "[]interface{}"
	OBJECT = "interface{}"
)

func main() {
	args := os.Args[1:]
	var wg sync.WaitGroup
	c := make(chan error)

	wg.Add(len(args))

	go func() {
		wg.Wait()
		close(c)
	}()

	for _, arg := range args {
		go func(arg string) {
			Convert(arg, c)
			wg.Done()
		}(arg)
	}

	for err := range c {
		if err != nil {
			fmt.Println(err)
		}
	}
}

func Convert(file string, c chan error) {
	if !isJSONFile(file) {
		c <- fmt.Errorf("input file must be json")
		return
	}

	data, err := readJsonFile(file)

	if err != nil {
		c <- fmt.Errorf("error converting json to go: %v", err)
		return
	}

	m, err := getJsonMap(data)

	if err != nil {
		c <- fmt.Errorf("error converting json to go: %v", err)
		return
	}

	filename := getFileNameNoPathExt(file)

	conv, err := buildGoStruct(m, filename)

	if err != nil {
		c <- fmt.Errorf("error converting json to go: %v", err)
		return
	}

	gofile := filename + ".go"
	_, err = writeGoFile(gofile, conv)
	cmd := exec.Command("go", "fmt", gofile)
	if err := cmd.Run(); err != nil {
		c <- fmt.Errorf("Cannot format file %v: %v", gofile, err)
		return
	}
}

func readJsonFile(file string) ([]byte, error) {
	if file == "" {
		return nil, fmt.Errorf("cannot open file with no name")
	}

	f, err := os.ReadFile(file)

	if err != nil {
		return nil, fmt.Errorf("cannot open file %v", file)
	}

	return f, nil
}

func writeGoFile(file, data string) (int, error) {
	f, err := os.Create(file)

	if err != nil {
		return 0, fmt.Errorf("cannot create file %v: %v", file, err)
	}

	defer f.Close()

	n, err := f.Write([]byte("package main\n" + data))

	if err != nil {
		return 0, fmt.Errorf("cannot write to file %v: %v", file, err)
	}

	return n, nil
}

func getJsonMap(input []byte) (*map[string]interface{}, error) {
	m := make(map[string]interface{})

	if err := json.Unmarshal(input, &m); err != nil {
		return nil, fmt.Errorf("could not parse json: %v", err)
	}

	return &m, nil
}

func buildGoStruct(m *map[string]interface{}, name string) (string, error) {
	obj, err := buildUnnamedStruct(m)

	if err != nil {
		return "", fmt.Errorf("cannot build unnamed struct: %v", err)
	}

	return fmt.Sprintf("type %v %v", capFirstLetter(name), obj), nil
}

func buildUnnamedStruct(m *map[string]interface{}) (string, error) {
	res := "struct {"

	for k, v := range *m {
		var r string
		var err error

		res += "\n"
		switch vt := v.(type) {
		case string:
			r, err = buildNVPair(k, STRING)
		case bool:
			r, err = buildNVPair(k, BOOL)
		case []interface{}:
			r, err = buildNVPair(k, ARRAY)
		case float64:
			// golang parser considers all json numbers as floats
			if v == float64(int(vt)) {
				r, err = buildNVPair(k, FLOAT)
			} else {
				r, err = buildNVPair(k, INT)
			}
		default:
			if m, ok := v.(map[string]interface{}); ok {
				r, err = buildNamedStruct(&m, k)
			} else {
				return "", fmt.Errorf("error building struct field: %v", ok)
			}
		}

		if err != nil {
			return "", fmt.Errorf("could not build name/value pair: %v", err)
		}

		res += r
	}

	res += "\n}"

	return res, nil
}

func buildNamedStruct(m *map[string]interface{}, name string) (string, error) {
	r, err := buildUnnamedStruct(m)

	if err != nil {
		return "", fmt.Errorf("error building unnamed struct: %v", err)
	}

	res, err := buildNVPair(name, r)

	if err != nil {
		return "", fmt.Errorf("%v", err)
	}

	return res, nil
}

func buildNVPair(name, value string) (string, error) {
	if len(name) == 0 {
		return "", fmt.Errorf("struct field cannot have empty name")
	}

	if len(value) == 0 {
		return "", fmt.Errorf("struct field type cannot be empty")
	}

	if !isChar(name[0]) {
		return fmt.Sprintf("Go%v %v `json:\"%v\"`", name, value, name), nil
	}

	return fmt.Sprintf("%v %v `json:\"%v\"`", capFirstLetter(name), value, name), nil
}

func capFirstLetter(s string) string {
	return cases.Title(language.English, cases.Compact).String(s)
}

func isChar(s byte) bool {
	return s >= 'a' && s <= 'z' || s >= 'A' && s <= 'Z'
}

func getFileNameNoPathExt(file string) string {
	filename := filepath.Base(file)
	return filename[:len(filename)-len(filepath.Ext(filename))]
}

func isJSONFile(file string) bool {
	return filepath.Ext(file) == ".json"
}
