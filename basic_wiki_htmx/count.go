package main

import (
	"os"
	"strconv"
)

const countFile = "data/count.txt"

func ReadCount() (int, error) {
	data, err := os.ReadFile(countFile)

	if err != nil {
		return 0, err
	}

	content := string(data)
	val, err := strconv.Atoi(content)

	if err != nil {
		return 0, err
	}

	return val, nil
}

func WriteCount(val int) error {
	s:= []byte(strconv.Itoa(val))
	return os.WriteFile(countFile, s, 0600)
}