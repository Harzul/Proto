package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"time"
)

var users = map[string]string{
	"admin": "61646d696ee3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
	"user":  "75736572e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
}

func main() {
	file, logger := initLogger()
	logger.Println("Программа запущена")
	err := initiate()
	if err != nil {
		logger.Fatalln("критичесая ошибка при инициализации: ", err)
	}
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

	scanner := bufio.NewScanner(os.Stdin)
	username := ""
	fmt.Print("Введите имя пользователя: ")
	if scanner.Scan() {
		username = scanner.Text()
	}
	password := ""
	var counter = 2
	for {
		fmt.Printf("Введите пароль пользователя %s: ", username)
		if scanner.Scan() {
			password = hex.EncodeToString(sha256.New().Sum([]byte(scanner.Text())))
		}
		if users[username] == password {
			break
		} else {
			fmt.Printf("Пароль/имя пользователя неверено, осталось попыток %d\n", counter)
			counter -= 1
		}
		if counter == -1 {
			logger.Fatalln("исчерпано число попыток входа: ", err)
		}
	}

	if username == "admin" && users[username] == password {
		file, err := os.Open("secret.key")
		if err != nil {
			logger.Fatalln("ошибка открытия файла")
		}
		defer file.Close()

		fmt.Println("Выберите задачу:\n1-Изменить позицию генерации ключа\n2-Что-то еще...  ")

		task := ""
		if scanner.Scan() {
			task = scanner.Text()
		}
		if task == "1" {
			value := ""
			fmt.Println("Введите новое значение: ")
			if scanner.Scan() {
				value = scanner.Text()
			}
			fmt.Println("Изменяем отступ ключа")
			file, err := os.Create("secret.key")
			if err != nil {
				logger.Fatalln("ошибка обновления файла")
			}
			defer file.Close()
			file.WriteString(value)
			logger.Println("Отступ ключа изменен")
		}
	} else if (username != "admin") && (users[username] == password) {
		fmt.Println("Введите название файле с параметрами: ")
		paramFile := ""
		if scanner.Scan() {
			paramFile = scanner.Text()
		}
		file, err := os.Open(paramFile)
		if err != nil {
			fmt.Println("Ошибка открытия файла:", err)
		}

		scanner := bufio.NewScanner(file)
		var key int
		if scanner.Scan() {
			key, err = strconv.Atoi(scanner.Text())
			if err != nil {
				logger.Fatalln("ошибка преобразования str в int")
			}
		}
		if err := scanner.Err(); err != nil {
			fmt.Println("Ошибка чтения:", err)
		}
		defer file.Close()
		scanner = bufio.NewScanner(os.Stdin)
		var s Splitmix64
		var x Xoshiro256_PP

		var seed = make([]uint64, 4)
		for range int(key * 4) {
			s.Next()
		}
		for i := range 4 {
			seed[i] = s.Next()
		}
		x.S = [4]uint64(seed)

		fmt.Println("Введите количество uint64 значений: ")
		valsNum := -1
		if scanner.Scan() {
			valsNum, err = strconv.Atoi(scanner.Text())
			if err != nil {
				logger.Fatalln("ошибка преобразования str в int")
			}
		}

		file_res, err := os.Create("output.txt")
		if err != nil {
			logger.Fatalln("ошибка создания файла вывода")
		}
		defer file_res.Close()

		for range valsNum {
			file_res.WriteString(strconv.FormatUint(x.Next(), 10) + "\n")
		}

		src := rand.NewSource(time.Now().UnixNano())
		r := rand.New(src)
		for range 5 {
			for index := range x.S {
				x.S[index] = r.Uint64()
			}
		}

	}
	logger.Println("Программа завершена")
	file.Close()
}
