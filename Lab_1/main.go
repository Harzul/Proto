package main //«Магма» OFB
import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

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
	val1, err := strconv.ParseInt(string(tmp), 16, 64)
	if err != nil {
		fmt.Println("Can't convert hex to int:", err)
		os.Exit(-1)
	}
	val2, err := strconv.ParseInt(string(a1), 16, 64)
	if err != nil {
		fmt.Println("Can't convert hex to int:", err)
		os.Exit(-1)
	}
	return strconv.FormatInt(val1^val2, 16) + a0
}
func encrypt(a, key string) string {
	var (
		a0 = a[8:]
		a1 = a[:8]
	)
	var round_keys = generete_round_keys(key)
	var state = []string{a1, a0}
	for i := range ROUNDS - 1 {
		state = G(round_keys[i], state[0], state[1])
	}
	return G_last(round_keys[31], state[0], state[1])
}

func decrypt(b, key string) string {
	var (
		b0 = b[8:]
		b1 = b[:8]
	)
	var round_keys = generete_round_keys(key)
	var state = []string{b1, b0}
	for i := ROUNDS - 1; i > 0; i-- {
		state = G(round_keys[i], state[0], state[1])
	}
	return G_last(round_keys[0], state[0], state[1])
}

// 11fe7a6d
// f - test:   g[87654321](fedcba98) = fdcbc20c
// T - test:   t(fdb97531) = 2a196f34
func main() {
	//fmt.Println(t("fdb97531"))
	//fmt.Println(g("fdcbc20c", "87654321"))
	//fmt.Println(G("fcfdfeff", "8025c0a5", "b0d66514"))
	fmt.Println("fedcba9876543210")
	cipher := encrypt("fedcba9876543210", "ffeeddccbbaa99887766554433221100f0f1f2f3f4f5f6f7f8f9fafbfcfdfeff")
	fmt.Println(cipher)
	opened := decrypt(cipher, "ffeeddccbbaa99887766554433221100f0f1f2f3f4f5f6f7f8f9fafbfcfdfeff")
	fmt.Println(opened)
}
