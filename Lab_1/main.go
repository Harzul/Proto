package main //«Магма» OFB
import (
	"bytes"
	"encoding/hex"
	"fmt"
	"log"
	"math"
	"os"
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

func generete_round_keys(key []byte) [][]byte {
	var keys [][]byte
	for range 3 {
		for i := range 8 {
			keys = append(keys, key[i*8:i*8+8])
		}
	}
	for i := range 8 {
		keys = append(keys, key[64-(i*8+8):64-i*8])
	}
	return keys
}

func t(a []byte) []byte {
	var result = make([]byte, 8)
	for i := range 8 {
		result[i] = byte(S_BOX[7-i][a[i]])
	}
	return result
}

func g(round_key, a []byte) []byte {
	var (
		result = make([]byte, 8)
		val1   uint32
		val2   uint32
		val3   uint32
	)
	for _, nibble := range a {
		val1 = (val1 << 4) | uint32(nibble)
	}
	for _, nibble := range round_key {
		val2 = (val2 << 4) | uint32(nibble)
	}
	tmp := uint32(val1 + val2)
	for i := 7; i >= 0; i-- {
		result[i] = byte(tmp & 0xF)
		tmp >>= 4
	}

	result = t(result)

	for _, nibble := range result {
		val3 = (val3 << 4) | uint32(nibble)
	}

	tmp = rotateLeft32(val3, 11)
	for i := 7; i >= 0; i-- {
		result[i] = byte(tmp & 0xF)
		tmp >>= 4
	}
	return result
}

func G(round_key, a1, a0 []byte) [][]byte {
	tmp := g(round_key, a0)

	var (
		result = make([][]byte, 2)
		val1   uint32
		val2   uint32
	)
	for _, nibble := range tmp {
		val1 = (val1 << 4) | uint32(nibble)
	}
	for _, nibble := range a1 {
		val2 = (val2 << 4) | uint32(nibble)
	}

	n := uint32(val1 ^ val2)

	tmp = make([]byte, 8)
	for i := 7; i >= 0; i-- {
		tmp[i] = byte(n & 0xF)
		n >>= 4
	}
	result[0] = a0
	result[1] = tmp
	return result
}

func G_last(round_key, a1, a0 []byte) []byte {
	tmp := g(round_key, a0)

	var (
		result []byte
		val1   uint32
		val2   uint32
	)
	for _, nibble := range tmp {
		val1 = (val1 << 4) | uint32(nibble)
	}
	for _, nibble := range a1 {
		val2 = (val2 << 4) | uint32(nibble)
	}

	n := uint32(val1 ^ val2)
	tmp = make([]byte, 8)
	for i := 7; i >= 0; i-- {
		tmp[i] = byte(n & 0xF)
		n >>= 4
	}
	result = append(result, tmp...)
	result = append(result, a0...)

	return result
}

func magic(a, IV, key []byte) []byte {
	var (
		blocks         = int(math.Ceil(float64(len(a)) / 16))
		result         = []byte{}
		round_keys     = generete_round_keys(key)
		temp_iv        = IV
		a0, a1         = make([]byte, 16), make([]byte, 16)
		n              = []byte{}
		current_cipher = []byte{}
		tmp            = make([]byte, blocks)
		val1           uint64
		val2           uint64
	)
	for i := range blocks {
		if i == blocks-1 {
			tmp = a[i*8*2:]
		} else {
			tmp = a[i*8*2 : (i+1)*8*2]
		}
		n = temp_iv[:16]
		a0 = n[8:]
		a1 = n[:8]
		var state = [][]byte{a1, a0}
		for i := range ROUNDS - 1 {
			state = G(round_keys[i], state[0], state[1])
		}
		current_cipher = G_last(round_keys[31], state[0], state[1])

		for _, nibble := range tmp {
			val1 = (val1 << 4) | uint64(nibble)
		}
		for _, nibble := range current_cipher {
			val2 = (val2 << 4) | uint64(nibble)
		}
		for index, val := range tmp {
			current_cipher[index] = current_cipher[index] ^ tmp[index]
			_ = val
		}
		temp_iv = append(temp_iv[16:], current_cipher...)
		result = append(result, current_cipher...)
	}
	return result
}

func main() {
	data, err := os.ReadFile("test_10MB.txt")
	if err != nil {
		fmt.Println("Ошибка чтения файла:", err)
		return
	}
	t1 := make([]byte, len(data)*2)
	for i, b := range data {
		t1[i*2] = (b >> 4) & 0x0F
		t1[i*2+1] = b & 0x0F
	}
	IV, _ := hex.DecodeString("1234567890abcdef234567890abcdef1")
	t2 := make([]byte, len(IV)*2)
	for i, b := range IV {
		t2[i*2] = (b >> 4) & 0x0F
		t2[i*2+1] = b & 0x0F
	}
	key, _ := hex.DecodeString("ffeeddccbbaa99887766554433221100f0f1f2f3f4f5f6f7f8f9fafbfcfdfeff")
	t3 := make([]byte, len(key)*2)
	for i, b := range key {
		t3[i*2] = (b >> 4) & 0x0F
		t3[i*2+1] = b & 0x0F
	}

	start := time.Now()
	cipher := magic(t1, t2, t3)

	duration := time.Since(start)
	fmt.Printf("Время выполнения шифрования: %v\n", duration)
	bt := make([]byte, len(cipher)/2)
	for i := 0; i < len(bt); i++ {
		bt[i] = (cipher[i*2] << 4) | (cipher[i*2+1] & 0xF)
	}
	data, err = hex.DecodeString(hex.EncodeToString(bt))
	if err != nil {
		log.Fatalf("Error decoding hex data: %v", err)
	}
	err = os.WriteFile("test_10MB_C.txt", data, 0644)
	if err != nil {
		log.Fatalf("Error writing to output file: %v", err)
	}
	data, err = os.ReadFile("test_10MB_C.txt")
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}
	t1 = make([]byte, len(data)*2)
	for i, b := range data {
		t1[i*2] = (b >> 4) & 0x0F
		t1[i*2+1] = b & 0x0F
	}
	start = time.Now()
	cipher = magic(t1, t2, t3)
	duration = time.Since(start)

	fmt.Printf("Время выполнения расшифрования: %v\n", duration)
	bt = make([]byte, len(cipher)/2)
	for i := 0; i < len(bt); i++ {
		bt[i] = (cipher[i*2] << 4) | (cipher[i*2+1] & 0xF)
	}

	err = os.WriteFile("test_10MB_D.txt", []byte(string(bytes.TrimRight(bt, "\x00"))), 0644)
	if err != nil {
		log.Fatalf("Error writing to output file: %v", err)
	}

}
