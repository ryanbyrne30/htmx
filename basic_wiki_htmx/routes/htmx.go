package main

import (
	"net/http"
	"os"
	"strconv"
)

const countFile = "data/count.txt"

func ButtonClickHandler(w http.ResponseWriter, r *http.Request) {
	count, err := ReadCount()
	if err != nil {
		handleError(w, err)
	}
	err = templates.ExecuteTemplate(w, "click.html", count)
	if err != nil {
		handleError(w, err)
	}
	err = WriteCount(count + 1)
	if err != nil {
		handleError(w, err)
	}
}


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