package main

import (
	"encoding/binary"
	"fmt"
	"math"
	"math/rand"
	"slices"
	"time"
)

func getBt(cipher []byte) []byte {
	bt := make([]byte, len(cipher)/2)
	for i := 0; i < len(bt); i++ {
		bt[i] = (cipher[i*2] << 4) | (cipher[i*2+1] & 0xF)
	}
	return bt
}

func X(a, b []uint8) {
	for i := range 64 {
		a[i] ^= b[i]
	}
}

func zero(a []uint8) {
	for i := range a {
		a[i] = 0
	}
}

func get64(a []uint8) uint64 {
	return binary.BigEndian.Uint64(a)
}

func getBytes(a []uint8, dig uint64) {
	for i := 7; i >= 0; i-- {
		a[i] = uint8(dig & 0xFF)
		dig >>= 8
	}
}

func l(a uint64) uint64 {
	var result uint64 = 0
	for i := 0; i < 64; i++ {
		if (a>>i)&1 == 1 {
			result ^= A[63-i]
		}
	}
	return result
}

func SPL(a []uint8) {
	for i := range 64 {
		a[i] = PI[a[i]]
	}
	temp := make([]uint8, 64)
	copy(temp, a)
	for i := range 64 {
		a[i] = temp[T[i]]
	}
	for i := range 8 {
		getBytes(a[i*8:], l(get64(a[i*8:])))
	}
}

func keySchedule(k []uint8, i int) {
	X(k, C[i])
	SPL(k)
}

func E(m, k []uint8) {
	for i := range 12 {
		X(m, k)
		SPL(m)
		keySchedule(k, i)
	}
	X(m, k)
}

func g_N(N, h, m []uint8) {
	LPS := make([]uint8, 64)
	copy(LPS, h)
	tempm := make([]uint8, 64)
	copy(tempm, m)
	X(LPS, N)
	SPL(LPS)
	E(tempm, LPS)
	X(tempm, h)
	X(tempm, m)
	copy(h, tempm)
}

func initiateHash(out int) {
	for i := range 64 {
		if out == 256 {
			IV[i] = 0x01
		} else if out == 512 {
			IV[i] = 0x00
		}
		N[i] = 0x00
		Sigma[i] = 0x00
	}
	N_512 = make([]uint8, 64)
	N_512[62] = 0x02
	N_0 = make([]uint8, 64)
}

func add(a, b []uint8) {
	var t int = 0
	for i := 63; i >= 0; i-- {
		t = int(a[i]) + int(b[i]) + (t >> 8)
		a[i] = uint8(t & 0xFF)
	}
}

func alg(message, h []uint8) {
	var inc int = 0
	var tlen int = len(message)
	for tlen >= 64 {
		inc++
		g_N(N, IV, message[len(message)-inc*64:])
		add(N, N_512)
		add(Sigma, message[len(message)-inc*64:])
		tlen -= 64
	}
	var len1 = len(message) - inc*64
	paddedMes := make([]uint8, 64)
	if len(message)-inc*64 < 64 {
		for i := 0; i < (64 - len1 - 1); i++ {
			paddedMes[i] = 0
		}
		paddedMes[64-len1-1] = 0x01
		copy(paddedMes[64-len1:], message)
	}
	g_N(N, IV, paddedMes)
	a := make([]uint8, 64)
	a[63] = uint8(((len1) * 8) % 256)
	a[62] = uint8((((len1) * 8) / 256) % 256)
	add(N, a)
	add(Sigma, paddedMes)
	g_N(N_0, IV, N)
	g_N(N_0, IV, Sigma)
	copy(h, IV)
}

func get256(message, h []uint8) {
	initiateHash(256)
	temph := make([]uint8, 64)
	alg(message, temph)
	copy(h, temph)
}

func Hmac256(ret, K, T []uint8, Klen, Tlen int) {
	ipad := make([]uint8, 64)
	opad := make([]uint8, 64)
	for i := range 64 {
		ipad[i] = 0x36
		opad[i] = 0x5C
	}
	var ilen uint = uint(64 + Tlen)
	inner := make([]uint8, ilen)
	tK := make([]uint8, 64)
	copy(tK, K)
	X(tK, ipad)
	copy(inner, tK)
	copy(inner[64:], T)
	X(tK, ipad)
	innerh := make([]uint8, 32)
	slices.Reverse(inner)
	get256(inner, innerh)
	slices.Reverse(innerh)
	outer := make([]uint8, 96)
	X(tK, opad)
	copy(outer, tK)
	copy(outer[64:], innerh)
	slices.Reverse(outer)
	zero(ret)
	get256(outer, ret)
	slices.Reverse(ret)
}
func bytess(arr []uint8, step int) {
	for j := len(arr) - 1; j >= 0; j-- {
		arr[j] = uint8(step % 256)
		step = step / 256
	}
}
func kdf_tree(ret, K []uint8, klen, R, l int) {
	if !(R > 0 && R < 5) {
		fmt.Printf("Wrong R\n")
		return
	}
	if !(l >= 256) {
		fmt.Printf("Wrong l\n")
	}
	var Llen int = R + 1
	var Tlen int = R + 4 + 1 + 8 + Llen
	I := make([]uint8, R)
	T := make([]uint8, Tlen)
	L := make([]uint8, Llen)
	bytess(L, l)
	copy(T[R:], label)
	copy(T[R+4+1:], seed)
	copy(T[R+1+4+8:], L)
	res := make([]uint8, 32)
	for i := 1; i <= l/256; i++ {
		bytess(I, i)
		copy(T, I)
		Hmac256(res, K, T, klen, Tlen)
		copy(ret[(i-1)*32:], res)
	}
}
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
