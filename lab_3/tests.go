package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"os"
	"reflect"
)

func testAlgo() error {

	predict := []string{
		"ce84432f8890d84d0f7e93e6ac140dc99406174bafa6b56904817be1df51246a",
		"f58dbcd590c89cea62d608691534680cbf00cd600bb2f425caf6092eab272b71",
		"4a61ee21996a8eebf2178b2d4f38f7bea4a6d259f1985b04b5158db89ee3a351",
		"0d01a136ec0cebd055515079821023ccb825a27e42abdb2ec77a896e41367eae",
		"e3e5c1c13532b007ee3c8c69933069da0f04769a2b0e75b8a1d5cc516c007c60",
		"1a5a933601e7074db42ecdde06cbd4aa8f1582388c1d011286f24583b4925e27",
		"e859646ad9fd193534015b49385a09da1b42c01572d264ece8b65f9a19ad00e7",
		"b76e69eb3ec129fe337250a74395dbe439e9042e27451469962500c00c2315eb",
		"b581baeb0b78995a945094efdb8806ca458a38f79e217b71e830fcc445ecf78a",
		"d4acc9d149d201a7b2c20173f6720d4830f9e47003193f875312576e09e78e79",
	}
	nonce := []byte("nonce1234")
	entropy := []byte("Chumanov Nikita KKSO-03-20")

	drbg, err := NewHashDRBG(entropy, nonce, nil)
	if err != nil {
		return errors.New("контроль функциональности не пройден")
	}
	temp := make([]string, 10)
	for i := range 10 {
		output, err := drbg.Generate(32, nil)
		if err != nil {
			return errors.New("контроль функциональности не пройден")
		}
		temp[i] = hex.EncodeToString(output)
	}
	if !(reflect.DeepEqual(predict, temp)) {
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
