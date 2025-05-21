package main //«Магма» OFB
import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

var users = map[string]string{
	"admin": "61646d696ee3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
	"user":  "75736572e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
}

/*
Готово:
Логирование
Затирка ключа
"Secret" file
диагностический контроль (самотестирование)
самоблокироваться если обнаружено нарушение целостности
Контроль времени жизни ключа +/-
Обновление времени жизни СКЗИ
ограничения числа попыток аутентификации
Остальное либо закрывается астрой либо реализуется гораздо сложнее
*/

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
		t2, t3 := readParams(paramFile) //"secret.key"
		data := readData(filename + ".txt")
		t1 := make([]byte, len(data)*2)
		for i, b := range data {
			t1[i*2] = (b >> 4) & 0x0F
			t1[i*2+1] = b & 0x0F
		}
		fmt.Println("Выберите задачу:\n1-Зашифровать\n2-Расшифровать ")
		task := ""
		if scanner.Scan() {
			task = scanner.Text()
		}
		if task == "1" {
			start := time.Now()
			logger.Println("Шифрование начато")
			cipher := magic(t1, t2, t3)
			logger.Println("Шифрование завершено")
			duration := time.Since(start)
			fmt.Printf("Время выполнения шифрования: %v\n", duration)
			bt := getBytes(cipher)
			writeData(bt, filename+"_C.txt")
		}
		if task == "2" {
			t1 = readData(filename + "_C.txt")
			t2, t3 = readParams(paramFile)
			start := time.Now()
			logger.Println("Расшифрование начато")
			cipher := magic(t1, t2, t3)
			logger.Println("Расшифрование завершено")
			duration := time.Since(start)
			fmt.Printf("Время выполнения расшифрования: %v\n", duration)
			bt := getBytes(cipher)
			writeData([]byte(string(bytes.TrimRight(bt, "\x00"))), filename+"_D.txt")
		}

		
	}
	file.Close()
}
