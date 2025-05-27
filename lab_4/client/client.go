package main

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
		file, err := os.Open("config.json")
		if err != nil {
			logger.Fatalln("ошибка открытия файла")
		}
		defer file.Close()

		data, err := io.ReadAll(file)
		if err != nil {
			logger.Fatalln("ошибка чтения файла")
		}
		var conf Config
		err = json.Unmarshal(data, &conf)
		if err != nil {
			logger.Fatalln("ошибка декодирования JSON")
		}
		fmt.Println("Выберите задачу:\n1-Изменить ключ\n2-Изменить время жизни СКЗИ\n3-Что-то еще...  ")

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
		} else if task == "2" {
			var t = conf.TimeLimit
			fmt.Println("Добавляем Год к времени активности")
			conf.TimeLimit = t.Add(8760 * time.Hour)
			file, err = os.Create("config.json")
			if err != nil {
				logger.Fatalln("ошибка обновления JSON")
			}
			defer file.Close()

			encoder := json.NewEncoder(file)
			encoder.SetIndent("", "  ")
			if err := encoder.Encode(conf); err != nil {
				logger.Fatalln("ошибка обновления JSON")
			}
			logger.Println("Время жизни изменено")
		}
	} else if (username != "admin") && (users[username] == password) {
		fmt.Println("Введите название файла для обработки: ")

		filename := ""
		if scanner.Scan() {
			filename = scanner.Text()
		}
		fmt.Println("Введите название файле с параметрами: ")
		paramFile := ""
		if scanner.Scan() {
			paramFile = scanner.Text()
		}
		iv, key := readParams(paramFile) //"secret.key"
		data := readData(filename)
		fmt.Println("Сколько сообщений отправить?")
		tasks := ""
		if scanner.Scan() {
			tasks = scanner.Text()
		}
		num, err := strconv.Atoi(tasks)
		if err != nil {
			fmt.Println("Conversion error:", err)
			return
		}
		for i := range num {
			logger.Println("Начало генерации сообщения", err)
			size := 256 * (i + 1) //key_param
			res := make([]uint8, size/8)
			kdfTree(res, key, 32, 1, size)
			message := createMessage(data, iv, res[i*32:i*32+32])
			jsonData, err := json.Marshal(message)
			if err != nil {
				logger.Fatalln("Ошибка сериализации:", err)
				return
			}
			resp, err := http.Post("http://localhost:8080/submit", "application/json", bytes.NewBuffer(jsonData))
			if err != nil {
				logger.Fatalln("Ошибка запроса:", err)
				return
			}
			defer resp.Body.Close()

			logger.Println("Ответ сервера:", resp.Status)
		}

	}
	file.Close()
}
