package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"os"
	"time"
)

func testAlgo() error {
	var mes = []byte{
		0x92, 0xDE, 0xF0, 0x6B, 0x3C, 0x13, 0x0A, 0x59,
		0xDB, 0x54, 0xC7, 0x04, 0xF8, 0x18, 0x9D, 0x20,
		0x4A, 0x98, 0xFB, 0x2E, 0x67, 0xA8, 0x02, 0x4C,
		0x89, 0x12, 0x40, 0x9B, 0x17, 0xB5, 0x7E, 0x41,
		0x92, 0xDE, 0xF0, 0x6B, 0x3C, 0x13, 0x0A, 0x59,
		0xDB, 0x54, 0xC7, 0x04, 0xF8, 0x18, 0x9D, 0x20,
		0x4A, 0x98, 0xFB, 0x2E, 0x67, 0xA8, 0x02, 0x4C,
		0x89, 0x12, 0x40, 0x9B, 0x17, 0xB5, 0x7E, 0x41,
		0x92, 0xDE, 0xF0, 0x6B, 0x3C, 0x13, 0x0A, 0x59,
		0xDB, 0x54, 0xC7, 0x04, 0xF8, 0x18, 0x9D, 0x20,
	}
	var iv = []byte{
		0x12, 0x34, 0x56, 0x78, 0x90, 0xab, 0xcd, 0xef,
		0x23, 0x45, 0x67, 0x89, 0x0a, 0xbc, 0xde, 0xf1,
	}
	var K []uint8 = []uint8{
		0x00, 0x01, 0x02, 0x03,
		0x04, 0x05, 0x06, 0x07,
		0x08, 0x09, 0x0a, 0x0b,
		0x0c, 0x0d, 0x0e, 0x0f,
		0x10, 0x11, 0x12, 0x13,
		0x14, 0x15, 0x16, 0x17,
		0x18, 0x19, 0x1a, 0x1b,
		0x1c, 0x1d, 0x1e, 0x1f,
	}
	size := 256 * (1 + 1) //key_param
	res := make([]uint8, size/8)
	kdfTree(res, K, 32, 1, size)
	message := createMessage(mes, iv, res[0:32])
	jsonData, err := json.Marshal(message)
	if err != nil {
		return errors.New("ошибка провекрки функционалки")
	}
	good := Message{
		Header: Header{
			ExternalKeyIdFlag: "1",
			Version:           "0",
			CS:                "111000",
			KeyId:             "10000000",
			SeqNum:            "0",
		},
		PayloadData: "de9f19d019f4497edd6da71e0d7d968589e4961717236b2b66946de4b93edccbd2df9e05d1a3e4e4b727c0ddb3d5119b5f02276aed2480fe1c8e75022ff4c0519240012f29fbcca6ce4197ab0d4875dd",
		ICV:         "082a630662158c99220ad778084770d6294fe8cfed2756f375be247b32093abb",
	}
	goodData, err := json.Marshal(good)
	if err != nil {
		return errors.New("ошибка провекрки функционалки")
	}
	if string(goodData[:]) != string(jsonData[:]) {
		return errors.New("контроль функциональности не пройден")
	}
	return nil
}

func checkSumm(conf Config) error {
	expectedHash := conf.Hash
	exePath, err := os.Executable()
	if err != nil {
		return errors.New("не удалось получить путь к исполняемому файлу")
	}
	f, err := os.Open(exePath)
	if err != nil {
		return errors.New("не удалось открыть файл")
	}
	defer f.Close()

	h := sha256.New()

	if _, err := io.Copy(h, f); err != nil {
		return errors.New("ошибка чтения файла")
	}
	actualHash := h.Sum(nil)

	if !(hex.EncodeToString(actualHash) == expectedHash) {
		return errors.New("контроль целостности не пройден")
	}

	return nil
}

func testDate(conf Config) error {
	if conf.TimeLimit.Before(time.Now()) {
		return errors.New("время работы СКЗИ (ключа СКЗИ) истекло, обратитесь к администратору")
	}
	return nil
}
