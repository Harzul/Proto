package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"math"
	"os"
	"reflect"
)

func testAlgo() error {
	var s Splitmix64
	var x Xoshiro256_PP

	var seed = make([]uint64, 4)

	for range int(math.Pow(2, 16)) {
		s.Next()
	}
	for i := range 4 {
		seed[i] = s.Next()
	}
	x.S = [4]uint64(seed)

	predict := []uint64{
		4332022498915112517,
		6777032227142239569,
		17275362995184283823,
		1225089805541233209,
		4554618018609422149,
		14421996961521040327,
		15942364711299469956,
		7510107321175034685,
		1879987344874185030,
		17005661933589271247,
	}

	var temp = make([]uint64, 10)
	for i := range 10 {
		temp[i] = x.Next()
	}

	if !(reflect.DeepEqual(predict, temp)) {
		return errors.New("контроль целостности не пройден")
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
