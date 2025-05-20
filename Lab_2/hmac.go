package main

import "fmt"

func Hmac256(ret []uint8, K []uint8, Klen int, T []uint8, Tlen int) {
	tK := make([]uint8, 64)
	zero(tK, 64)
	copy(K, 0, tK, 0, Klen)
	//fmt.Printf("padded key=\n\t");
	//print_arr(tK,64);
	ipad := make([]uint8, 64)
	opad := make([]uint8, 64)
	for i := range 64 {
		ipad[i] = 0x36
		opad[i] = 0x5C
	}
	var ilen uint = uint(64 + Tlen)
	//fmt.Printf("ilen=%d\n", ilen)
	inner := make([]uint8, ilen)
	X(tK, ipad)
	copy(tK, 0, inner, 0, 64)
	copy(T, 0, inner, 64, Tlen)
	//fmt.Printf("inner=\n\t")
	//print_arr(inner, int(ilen))
	X(tK, ipad)
	innerh := make([]uint8, 32)
	zero(innerh, 32)
	reverse(inner, int(ilen))
	get256(inner, int(ilen), innerh)
	//fmt.Printf("innerh=\n\t")
	//print_arr(innerh,32);
	reverse(innerh, 32)
	outer := make([]uint8, 96)
	X(tK, opad)
	copy(tK, 0, outer, 0, 64)
	copy(innerh, 0, outer, 64, 32)
	//fmt.Printf("outer=\n\t")
	//print_arr(outer,96);
	reverse(outer, 96)
	zero(ret, 32)
	get256(outer, 96, ret)
	reverse(ret, 32)
	//print_arr(ret,32);
}

func test_hmac() int {
	var K []uint8 = []uint8{
		0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
		0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
		0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17,
		0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
	}
	var T []uint8 = []uint8{
		0x01, 0x26, 0xbd, 0xb8, 0x78, 0x00, 0xaf, 0x21,
		0x43, 0x41, 0x45, 0x65, 0x63, 0x78, 0x01, 0x00,
	}
	ret := make([]uint8, 32)
	var t []uint8 = []uint8{
		0xa1, 0xaa, 0x5f, 0x7d, 0xe4, 0x02, 0xd7, 0xb3,
		0xd3, 0x23, 0xf2, 0x99, 0x1c, 0x8d, 0x45, 0x34,
		0x01, 0x31, 0x37, 0x01, 0x0a, 0x83, 0x75, 0x4f,
		0xd0, 0xaf, 0x6d, 0x7c, 0xd4, 0x92, 0x2e, 0xd9,
	}
	zero(ret, 32)
	test_stribog()
	Hmac256(ret, K, 32, T, 16)

	fmt.Printf("HMAC256 test:\n\t")
	print_arr(ret, 32)
	if cmp(ret, t, 32) == 0 {
		fmt.Printf("test failed\n")
		return 0
	}
	fmt.Printf("test valid\n")
	return 1
}
