package main

import (
	"bufio"
	"encoding/hex"
	"log"
	"os"
)

func readData(filename string) []byte {
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	return data
}

func readParams(filename string) ([]byte, []byte) {
	file, err := os.Open(filename)
	if err != nil {
		log.Println("Ошибка открытия файла:", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var data string
	if scanner.Scan() {
		data = scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		log.Println("Ошибка чтения:", err)
	}

	IV, _ := hex.DecodeString(data)

	if scanner.Scan() {
		data = scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		log.Println("Ошибка чтения:", err)
	}
	key, _ := hex.DecodeString(data)

	return IV, key
}
