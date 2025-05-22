package main

import (
	"errors"
	"fmt"
	"os"
)

func print_arr(a []uint8, nums int) error {
	file, err := os.Create("resultKeys.txt")
	if err != nil {
		return errors.New("ошибка создания файла")
	}
	defer file.Close()

	for j := range nums {
		for i := range 32 {
			if a[i+j*32] < 0x10 {
				fmt.Fprintf(file, "0%x", a[i+j*32])
			} else {
				fmt.Fprintf(file, "%x", a[i+j*32])
			}
		}
		fmt.Fprintf(file, "\n")
	}
	return nil
}

func X(a, b []uint8) {
	for i := range 64 {
		a[i] ^= b[i]
	}
}

func copy(a []uint8, n1 int, b []uint8, n2 int, len int) {
	for i := range len {
		b[n2+i] = a[n1+i]
	}
}

func zero(a []uint8, len int) {
	for i := range len {
		a[i] = 0
	}
}

func reverse(a []uint8, len int) {
	temp := make([]uint8, len)
	copy(a, 0, temp, 0, len)
	for i := range len {
		a[i] = temp[len-1-i]
	}
}

func get64(a []uint8) uint64 {
	var result uint64 = 0
	for i := range 7 {
		result ^= uint64(a[i])
		result <<= 8
	}
	result ^= uint64(a[7])
	return result
}

func getBytes(a []uint8, dig uint64) {
	var temp = dig
	for i := 7; i >= 0; i-- {
		a[i] = uint8(temp % 256)
		temp /= 256
	}
}

func l(a uint64) uint64 {
	var result uint64 = 0
	var temp = a
	for i := range 64 {
		if temp%2 == 1 {
			result ^= A[63-i]
		}
		temp /= 2
	}
	return result
}
func SPL(a []uint8) {
	for i := range 64 {
		a[i] = PI[a[i]]
	}
	temp := make([]uint8, 64)
	copy(a, 0, temp, 0, 64)
	for i := range 64 {
		a[i] = temp[T[i]]
	}
	t := make([]uint8, 8)
	temp = make([]uint8, 8)

	for i := range 8 {
		copy(a, i*8, t, 0, 8)
		copy(t, 0, temp, 0, 8)
		getBytes(t, l(get64(temp)))
		copy(t, 0, a, i*8, 8)
	}
}

func keySchedule(k []uint8, i int) {
	temp := make([]uint8, 64)
	copy(k, 0, temp, 0, 64)
	X(temp, C[i])
	SPL(temp)
	copy(temp, 0, k, 0, 64)
}

func E(m, k []uint8) {
	temp := make([]uint8, 64)
	copy(k, 0, temp, 0, 64)
	for i := range 12 {
		X(m, temp)
		SPL(m)
		keySchedule(temp, i)
	}
	X(m, temp)
}

func g_N(N, h, m []uint8) {
	LPS := make([]uint8, 64)
	copy(h, 0, LPS, 0, 64)
	tempm := make([]uint8, 64)
	copy(m, 0, tempm, 0, 64)
	X(LPS, N)
	SPL(LPS)
	E(tempm, LPS)
	X(tempm, h)
	X(tempm, m)
	copy(tempm, 0, h, 0, 64)
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
	zero(N_512, 64)
	N_512[62] = 0x02
	zero(N_0, 64)
}

func add(a, b []uint8) {
	tempa := make([]uint8, 64)
	tempb := make([]uint8, 64)
	copy(a, 0, tempa, 0, 64)
	copy(b, 0, tempb, 0, 64)
	var t int = 0
	for i := 63; i >= 0; i-- {
		t = int(tempa[i]) + int(tempb[i]) + (t >> 8)
		a[i] = uint8(t & 0xFF)
	}
}

func alg(message []uint8, len int, h []uint8) {
	temph := make([]uint8, 64)
	copy(IV, 0, temph, 0, 64)
	var inc int = 0
	var tlen int = len
	for tlen >= 64 {
		inc++
		tempmes := make([]uint8, 64)
		copy(message, len-inc*64, tempmes, 0, 64)
		g_N(N, temph, tempmes)
		add(N, N_512)
		add(Sigma, tempmes)
		tlen -= 64
	}
	var len1 = len - inc*64
	mes1 := make([]uint8, len1)
	copy(message, 0, mes1, 0, len1)
	paddedMes := make([]uint8, 64)
	if len-inc*64 < 64 {
		for i := 0; i < (64 - len1 - 1); i++ {
			paddedMes[i] = 0
		}
		paddedMes[64-len1-1] = 0x01
		copy(mes1, 0, paddedMes, 64-len1, len1)
	}
	g_N(N, temph, paddedMes)
	arr := make([]uint8, 64)
	arr[63] = uint8(((len1) * 8) % 256)
	arr[62] = uint8((((len1) * 8) / 256) % 256)
	add(N, arr)
	add(Sigma, paddedMes)
	g_N(N_0, temph, N)
	g_N(N_0, temph, Sigma)
	copy(temph, 0, h, 0, 64)
}

func get256(message []uint8, len int, h []uint8) {
	initiateHash(256)
	temph := make([]uint8, 64)
	alg(message, len, temph)
	copy(temph, 0, h, 0, 32)
}

func get512(message []uint8, len int, h []uint8) {
	initiateHash(512)
	alg(message, len, h)
}

func cmp(a []uint8, b []uint8, len int) int {
	for i := 0; i < len; i++ {
		if a[i] != b[i] {
			return 0
		}
	}
	return 1
}
func Hmac256(ret []uint8, K []uint8, Klen int, T []uint8, Tlen int) {
	tK := make([]uint8, 64)
	copy(K, 0, tK, 0, Klen)
	ipad := make([]uint8, 64)
	opad := make([]uint8, 64)
	for i := range 64 {
		ipad[i] = 0x36
		opad[i] = 0x5C
	}
	var ilen uint = uint(64 + Tlen)
	inner := make([]uint8, ilen)
	X(tK, ipad)
	copy(tK, 0, inner, 0, 64)
	copy(T, 0, inner, 64, Tlen)
	X(tK, ipad)
	innerh := make([]uint8, 32)
	reverse(inner, int(ilen))
	get256(inner, int(ilen), innerh)
	reverse(innerh, 32)
	outer := make([]uint8, 96)
	X(tK, opad)
	copy(tK, 0, outer, 0, 64)
	copy(innerh, 0, outer, 64, 32)
	reverse(outer, 96)
	zero(ret, 32)
	get256(outer, 96, ret)
	reverse(ret, 32)
}
func bytes(arr []uint8, len int, i int) {
	for j := len - 1; j >= 0; j-- {
		arr[j] = uint8(i % 256)
		i = i / 256
	}
}

func kdf_tree(ret []uint8, K []uint8, klen int, R int, l int) {
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
	bytes(L, Llen, l)
	copy(label, 0, T, R, 4)
	copy(seed, 0, T, R+4+1, 8)
	copy(L, 0, T, R+1+4+8, Llen)
	res := make([]uint8, 32)
	for i := 1; i <= l/256; i++ {
		bytes(I, R, i)
		copy(I, 0, T, 0, R)
		Hmac256(res, K, klen, T, Tlen)
		copy(res, 0, ret, (i-1)*32, 32)
	}
}
