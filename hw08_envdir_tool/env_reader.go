package main

import (
	"bufio"
	"os"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	env := make(Environment)
	for _, envFile := range files {
		envVarName := envFile.Name()
		info, err := envFile.Info()
		if err != nil {
			return nil, err
		}
		if info.Size() == 0 {
			env[envVarName] = EnvValue{NeedRemove: true}
		} else {
			f, err := os.Open(dir + "/" + envFile.Name())
			if err != nil {
				_ = f.Close()
				return nil, err
			}
			r := bufio.NewReader(f)
			envVarValue, err := readValue(r)
			if err != nil {
				_ = f.Close()
				return nil, err
			}
			_ = f.Close()
			env[envVarName] = EnvValue{Value: cleanValue(envVarValue), NeedRemove: false}
		}
	}
	return env, nil
}

func readValue(reader *bufio.Reader) (string, error) {

	isPrefix := true
	var err error
	var rawLine, outLine []byte

	for isPrefix && err == nil {
		rawLine, isPrefix, err = reader.ReadLine()
		outLine = append(outLine, rawLine...)
	}
	return string(outLine), err
}

func cleanValue(inValue string) string {
	if inValue == "" {
		return ""
	}
	v := strings.TrimRight(inValue, ` 	`)
	return strings.ReplaceAll(v, `\0`, "\n")
}
