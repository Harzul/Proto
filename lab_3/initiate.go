package main

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"time"
)

type Config struct {
	Hash      string    `json:"hash"`
	TimeLimit time.Time `json:"timeLimit"`
}

func initiate() error {
	file, err := os.Open("config.json")
	if err != nil {
		return errors.New("ошибка открытия файла")
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return errors.New("ошибка чтения файла")
	}

	var conf Config
	err = json.Unmarshal(data, &conf)
	if err != nil {
		return errors.New("ошибка декодирования JSON")
	}
	err = testAlgo()
	if err != nil {
		return err
	}

	err = checkSumm(conf)
	if err != nil {
		return err
	}

	return nil
}
