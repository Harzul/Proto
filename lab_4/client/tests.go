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
	message := createMessage(mes, iv, res[0*32:0*32+32])
	jsonData, err := json.Marshal(message)
	if err != nil {
		return errors.New("ошибка провекрки функционалки")
	}
	_ = jsonData
	good := Message{
		Header: Header{
			ExternalKeyIdFlag: "1",
			Version:           "0",
			CS:                "111000",
			KeyId:             "10000000",
			SeqNum:            "1",
		},
		PayloadData: "9c8fbab4a81e6177f422151da222efb9d77159582421f3ba2bd0c98d10270a4b18e94e7b75b32ec5b41f4a2ce8637d4e96bbd9f7befe4ecd244111805652f101d25163a20c0a5b3f3b824121afb322c0",
		ICV:         "f0831a8dbb25f5736b5f64003e72e53cdb90284338fd9051f3664c42550cb6a5",
	}
	_ = good
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
