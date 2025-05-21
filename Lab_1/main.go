package main //«Магма» OFB
import (
	"bytes"
	"encoding/hex"
	"fmt"
	"time"
)

/*
Готово:
Логирование
Затирка ключа

TODO:
 1. Права доступа
    User - крипта
    Admin - настройка, конфигурирование (ключа, времени жизни ..., время блокировки)
 2. Контроль аутентификации
 3. Контроль времени жизни ключа
 4. Контроль целостности ПО и журнала
 5. Очиста облостей памяти  - сделано
 6. Логирование доступное только админу
 7. диагностический контроль
 8. самоблокироваться если обнаружено нарушение целостности 8
 9. ограничения числа попыток аутентификации
*/

func main() {
	file, logger := initLogger()
	logger.Println("Программа запущена")
	t2, t3 := readParams()
	filename := "./tests/test_10MB"

	//t1 := readData(filename + ".txt")
	data, _ := hex.DecodeString("92def06b3c130a59db54c704f8189d204a98fb2e67a8024c8912409b17b57e41")
	t1 := make([]byte, len(data)*2)
	for i, b := range data {
		t1[i*2] = (b >> 4) & 0x0F
		t1[i*2+1] = b & 0x0F
	}
	start := time.Now()
	logger.Println("Шифрование начато")
	cipher := magic(t1, t2, t3)
	logger.Println("Шифрование завершено")
	duration := time.Since(start)
	fmt.Printf("Время выполнения шифрования: %v\n", duration)
	bt := getBytes(cipher)
	writeData(bt, filename+"_C.txt")
	fmt.Printf("%#v\n\n", cipher)

	t1 = readData(filename + "_C.txt")
	t2, t3 = readParams()
	start = time.Now()
	logger.Println("Расшифрование начато")
	cipher = magic(t1, t2, t3)
	logger.Println("Расшифрование завершено")
	duration = time.Since(start)
	fmt.Printf("Время выполнения расшифрования: %v\n", duration)
	bt = getBytes(cipher)
	writeData([]byte(string(bytes.TrimRight(bt, "\x00"))), filename+"_D.txt")
	fmt.Printf("%#v\n", cipher)

	file.Close()
}
