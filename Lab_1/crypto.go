package main

import (
	"math"
	"math/rand"
	"time"
)

func rotateLeft32(val uint32, shift uint) uint32 {
	return (val << shift) | (val >> (32 - shift))
}

func generete_round_keys(key []byte) [][]byte {
	keys := make([][]byte, 32)
	for i := range keys {
		keys[i] = make([]byte, 8)
	}

	for j := range 3 {
		for i := range 8 {
			copy(keys[i+j*8], key[i*8:i*8+8])
		}
	}
	for i := range 8 {
		copy(keys[i+24], key[64-(i*8+8):64-i*8])
	}
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)
	for range 5 {
		for index := range key {
			key[index] = byte(r.Intn(16))
		}
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
	)
	for index, val := range tmp {
		tmp[index] = a1[index] ^ tmp[index]
		_ = val
	}
	result[0] = a0
	result[1] = tmp
	return result
}

func G_last(round_key, a1, a0 []byte) []byte {
	tmp := g(round_key, a0)

	var (
		result []byte
	)
	for index, val := range tmp {
		tmp[index] = a1[index] ^ tmp[index]
		_ = val
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
		temp_iv = append(temp_iv[16:], current_cipher...)
		for index, val := range tmp {
			current_cipher[index] = current_cipher[index] ^ tmp[index]
			_ = val
		}
		result = append(result, current_cipher...)
	}
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)
	for range 5 {
		for index := range round_keys {
			for j := range 8 {
				round_keys[index][j] = byte(r.Intn(16))
			}
		}
	}
	return result
}
