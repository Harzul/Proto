package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
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
		fmt.Println("Выберите задачу:\n1-Изменить время жизни СКЗИ\n2-Что-то еще...  ")

		task := ""
		if scanner.Scan() {
			task = scanner.Text()
		}
		if task == "1" {
			var t = conf.TimeLimit
			fmt.Println("Добавляем Год к времени активности")
			conf.TimeLimit = t.Add(8760 * time.Hour)
			file, err := os.Create("config.json")
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
		keysNum := -1
		fmt.Println("Введите количество ключей: ")
		if scanner.Scan() {
			keysNum, err = strconv.Atoi(scanner.Text())
			if err != nil {
				logger.Fatalln("ошибка преобразования str в int")
			}
		}
		res := make([]uint8, keysNum*256/8)
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
		var data string
		if scanner.Scan() {
			data = scanner.Text()
		}
		if err := scanner.Err(); err != nil {
			fmt.Println("Ошибка чтения:", err)
		}
		key, _ := hex.DecodeString(data)
		defer file.Close()
		logger.Println("Выработка начата")
		kdfTree(res, key, 32, 1, keysNum*256)
		logger.Println("Выработка завершена")
		print_arr(res, keysNum)

		src := rand.NewSource(time.Now().UnixNano())
		r := rand.New(src)
		for range 5 {
			for index := range key {
				key[index] = byte(r.Intn(16))
			}
		}
	}

	file.Close()
}
