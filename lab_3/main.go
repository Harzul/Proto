package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
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

		fmt.Println("Выберите задачу:\n1-Изменить ключ\n2-Что-то еще...  ")

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
			fmt.Println("Изменяем ключ")
			file, err := os.Create("secret.key")
			if err != nil {
				logger.Fatalln("ошибка обновления файла")
			}
			defer file.Close()
			file.WriteString(value)
			logger.Println("Ключ изменен")
		}
	} else if (username != "admin") && (users[username] == password) {
		scanner := bufio.NewScanner(file)
		if err := scanner.Err(); err != nil {
			fmt.Println("Ошибка чтения:", err)
		}
		defer file.Close()
		scanner = bufio.NewScanner(os.Stdin)

		fmt.Println("Введите количество значений: ")
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

		fmt.Println("Введите название файле с параметрами: ")
		paramFile := ""
		if scanner.Scan() {
			paramFile = scanner.Text()
		}
		entropy, err := os.ReadFile(paramFile)
		if err != nil {
			logger.Fatalln("ошибка чтение секретного файла")
		}
		nonce := []byte("nonce1234")

		drbg, err := NewHashDRBG(entropy, nonce, nil)
		if err != nil {
			log.Fatal(err)
		}
		for range valsNum {
			output, err := drbg.Generate(32, nil)
			if err != nil {
				log.Fatal(err)
			}
			file_res.WriteString(hex.EncodeToString(output) + "\n")
		}

		src := rand.NewSource(time.Now().UnixNano())
		r := rand.New(src)
		for range 5 {
			drbg.V = []byte(fmt.Sprint(r.Uint64()))

		}

	}
	logger.Println("Программа завершена")
	file.Close()
}
