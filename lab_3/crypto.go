package main

import (
	"encoding/binary"
)

type HashDRBG struct {
	V         []byte
	C         []byte
	ReseedCtr uint64
}

func hashDF(input []byte, length int) []byte {
	counter := byte(1)
	output := []byte{}
	hLen := 32
	iterations := (length + hLen - 1) / hLen

	for i := 0; i < iterations; i++ {
		data := make([]uint8, 32)
		data = append(data, []byte{counter}...)
		data = append(data, input...)
		h := make([]uint8, 32)
		get256(data, h)
		output = append(output, h...)
		counter++
	}
	return output[:length]
}

func NewHashDRBG(entropy, nonce, personalization []byte) (*HashDRBG, error) {
	seedMaterial := append(append(entropy, nonce...), personalization...)
	V := hashDF(seedMaterial, 5)
	C := hashDF(append([]byte{0x00}, V...), 55)
	return &HashDRBG{
		V:         V,
		C:         C,
		ReseedCtr: 1,
	}, nil
}

func (drbg *HashDRBG) Generate(numBytes int, additionalInput []byte) ([]byte, error) {

	if len(additionalInput) > 0 {
		data := append([]byte{0x02}, append(drbg.V, additionalInput...)...)
		drbg.V = hashDF(data, 55)
	}

	requested := []byte{}
	vTemp := make([]byte, len(drbg.V))
	copy(vTemp, drbg.V)

	for len(requested) < numBytes {
		data := make([]uint8, 32)
		data = append(data, vTemp...)
		h := make([]uint8, 32)
		get256(data, h)
		requested = append(requested, h...)

		for i := len(vTemp) - 1; i >= 0; i-- {
			vTemp[i]++
			if vTemp[i] != 0 {
				break
			}
		}
	}

	requested = requested[:numBytes]

	data := make([]uint8, 32)
	data = append(data, []byte{0x03}...)
	data = append(data, drbg.V...)
	data = append(data, drbg.C...)
	h := make([]uint8, 32)
	get256(data, h)
	reseedCtrBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(reseedCtrBytes, drbg.ReseedCtr)
	data = append(data, reseedCtrBytes...)
	get256(data, h)

	for i := range drbg.V {
		drbg.V[i] ^= h[i]
	}

	drbg.ReseedCtr++
	return requested, nil
}

func X(a, b []uint8) {
	for i := range 64 {
		a[i] ^= b[i]
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
