package main //«Магма» OFB
import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

/*
TODO:
 1. Права доступа
    User - крипта
    Admin - настройка, конфигурирование (ключа, времени жизни ..., время блокировки)
 2. Контроль аутентификации
 3. Контроль времени жизни ключа
 4. Контроль целостности ПО и журнала
 5. Очиста облостей памяти  - сделано
 6. Логирование доступное только админу
*/
func rotateLeft32(val uint32, shift uint) uint32 {
	return (val << shift) | (val >> (32 - shift))
}

func generete_round_keys(key string) []string {
	var keys []string
	for range 3 {
		for i := range 8 {
			keys = append(keys, key[i*BLOCK_SIZE:i*BLOCK_SIZE+8])
		}
	}
	for i := range 8 {
		keys = append(keys, key[64-(i*BLOCK_SIZE+8):64-i*BLOCK_SIZE])
	}

	return keys
}

func t(a string) string {
	var result string
	for i := range 8 {
		value, err := strconv.ParseUint(string(a[i]), 16, 32)
		if err != nil {
			fmt.Println("Can't convert hex to int:", err)
			os.Exit(-1)
		}
		result += strconv.FormatUint(uint64(S_BOX[7-i][value]), 16)
	}
	return result
}

func g(round_key, a string) string {
	var result string
	val1, err := strconv.ParseUint(string(a), 16, 32)
	if err != nil {
		fmt.Println("Can't convert hex to int:", err)
		os.Exit(-1)
	}
	val2, err := strconv.ParseUint(string(round_key), 16, 32)
	if err != nil {
		fmt.Println("Can't convert hex to int:", err)
		os.Exit(-1)
	}

	result = strconv.FormatUint(uint64(val1+val2)%uint64(math.Pow(2, 32)), 16)
	result = t(strings.Repeat("0", 8-len(result)) + result)

	val3, err := strconv.ParseUint(string(result), 16, 32)
	if err != nil {
		fmt.Println("Can't convert hex to int:", err)
		os.Exit(-1)
	}

	result = strconv.FormatUint(uint64(rotateLeft32(uint32(val3), 11)), 16)
	result = strings.Repeat("0", 8-len(result)) + result

	return result[len(result)-8:8] + result[:len(result)-8]
}

func G(round_key, a1, a0 string) []string {
	tmp := g(round_key, a0)
	val1, err := strconv.ParseUint(string(tmp), 16, 32)
	if err != nil {
		fmt.Println("Can't convert hex to int:", err)
		os.Exit(-1)
	}
	val2, err := strconv.ParseUint(string(a1), 16, 32)
	if err != nil {
		fmt.Println("Can't convert hex to int:", err)
		os.Exit(-1)
	}

	return []string{a0, strconv.FormatUint(val1^val2, 16)}
}

func G_last(round_key, a1, a0 string) string {
	tmp := g(round_key, a0)
	val1, err := strconv.ParseUint(string(tmp), 16, 32)
	if err != nil {
		fmt.Println("Can't convert hex to int:", err)
		os.Exit(-1)
	}
	val2, err := strconv.ParseUint(string(a1), 16, 32)
	if err != nil {
		fmt.Println("Can't convert hex to int:", err)
		os.Exit(-1)
	}

	return strconv.FormatUint(val1^val2, 16) + a0
}

func magic(flag, a, IV, key string) string {
	var (
		blocks         = int(math.Ceil(float64(len(a)) / 16))
		result         = ""
		round_keys     = generete_round_keys(key)
		temp_iv        = IV
		a0, a1         = "", ""
		current_iv     = ""
		current_cipher = ""
		tmp            = ""
	)
	for i := 0; i < 5; i++ {
		key = randomHexString(64)
		_ = len(key)
	}

	for i := range blocks {
		if i == blocks-1 {
			tmp = a[i*BLOCK_SIZE*2:]
		} else {
			tmp = a[i*BLOCK_SIZE*2 : (i+1)*BLOCK_SIZE*2]
		}
		current_iv = temp_iv[:16]
		a0 = current_iv[8:]
		a1 = current_iv[:8]

		var state = []string{a1, a0}
		for i := range ROUNDS - 1 {
			state = G(round_keys[i], state[0], state[1])
		}
		current_cipher = G_last(round_keys[31], state[0], state[1])
		temp_iv = temp_iv[len(temp_iv)-16:] + current_cipher

		val1, err := strconv.ParseUint(string(tmp), 16, 64)
		if err != nil {
			fmt.Println("Can't convert hex to int:", err)
			os.Exit(-1)
		}
		val2, err := strconv.ParseUint(string(current_cipher), 16, 64)
		if err != nil {
			fmt.Println("Can't convert hex to int:", err)
			os.Exit(-1)
		}

		data := strconv.FormatUint(val1^val2, 16)
		if flag == "e" {
			data = strings.Repeat("0", 16-len(data)) + data
		}

		result += data
	}
	for j := range ROUNDS - 1 {
		for i := 0; i < 5; i++ {
			round_keys[j] = randomHexString(64)
			_ = len(key)
		}
	}
	return result
}
func randomHexString(n int) string {
	bytes := make([]byte, (n+1)/2)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(err)
	}
	hexStr := hex.EncodeToString(bytes)
	return hexStr[:n]

}
func main() {
	logFile, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Ошибка открытия файла для логов: %v", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)
	log.SetPrefix("[APP] ")
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Println("Приложение запущено")
	/*if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <filename>")
		return
	}*/

	/*if len(os.Args) < 4 {
		fmt.Println("Usage: go run main.go <mode> <input_file> <output_file>")
		return
	}*/
	IV := "1234567890abcdef234567890abcdef1"
	key := "ffeeddccbbaa99887766554433221100f0f1f2f3f4f5f6f7f8f9fafbfcfdfeff"

	start := time.Now()
	data, err := os.ReadFile("test_10MB.txt")
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}
	text := hex.EncodeToString(data)
	log.Println("Шифрование начато")
	cipher := magic("e", text, IV, key)
	log.Println("Шифрование завершено")
	data, err = hex.DecodeString(string(cipher))
	if err != nil {
		log.Fatalf("Error decoding hex data: %v", err)
	}
	err = os.WriteFile("test_10MB_C.txt", data, 0644)
	if err != nil {
		log.Fatalf("Error writing to output file: %v", err)
	}
	duration := time.Since(start)
	fmt.Printf("Время выполнения шифрования: %v\n", duration)

	start = time.Now()
	data, err = os.ReadFile("test_10MB_C.txt")
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}
	text = hex.EncodeToString(data)
	log.Println("Расшифрование начато")
	deciphered := magic("d", text, IV, key)
	log.Println("Расшифрование завершено")
	data, err = hex.DecodeString(string(deciphered))
	if err != nil {
		log.Fatalf("Error decoding hex data: %v", err)
	}
	err = os.WriteFile("test_10MB_D.txt", data, 0644)
	if err != nil {
		log.Fatalf("Error writing to output file: %v", err)
	}
	duration = time.Since(start)
	fmt.Printf("Время выполнения расшифрования: %v\n", duration)

	for i := 0; i < 5; i++ {
		key = randomHexString(64)
		_ = len(key)
	}

	log.Println("Приложение закрыто")
}
