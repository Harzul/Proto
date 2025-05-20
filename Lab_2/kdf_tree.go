package main

import "fmt"

var label []uint8 = []uint8{
	0x26, 0xbd, 0xb8, 0x78,
}

var seed []uint8 = []uint8{
	0xaf, 0x21, 0x43, 0x41, 0x45, 0x65, 0x63, 0x78,
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
	I := make([]uint8, R)
	var Llen int = R + 1
	L := make([]uint8, Llen)
	bytes(L, Llen, l)
	var Tlen int = R + 4 + 1 + 8 + Llen
	fmt.Printf("KDF_TREE test:\n")
	T := make([]uint8, Tlen)
	zero(T, Tlen)
	copy(label, 0, T, R, 4)
	copy(seed, 0, T, R+4+1, 8)
	copy(L, 0, T, R+1+4+8, Llen)
	res := make([]uint8, 32)
	for i := 1; i <= l/256; i++ {
		bytes(I, R, i)
		copy(I, 0, T, 0, R)
		Hmac256(res, K, klen, T, Tlen)
		fmt.Printf("K%d=\n\t", i)
		print_arr(res, 32)
		copy(res, 0, ret, (i-1)*32, 32)
	}
}

func testKdf_tree() int {
	//test_hmac()
	var K []uint8 = []uint8{
		0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
		0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
	}
	size := 2048 //key_param
	res := make([]uint8, size/8)
	kdf_tree(res, K, 32, 1, size)
	var test []uint8 = []uint8{
		0x22, 0xb6, 0x83, 0x78, 0x45, 0xc6, 0xbe, 0xf6, 0x5e, 0xa7, 0x16, 0x72, 0xb2, 0x65, 0x83, 0x10,
		0x86, 0xd3, 0xc7, 0x6a, 0xeb, 0xe6, 0xda, 0xe9, 0x1c, 0xad, 0x51, 0xd8, 0x3f, 0x79, 0xd1, 0x6b,
		0x07, 0x4c, 0x93, 0x30, 0x59, 0x9d, 0x7f, 0x8d, 0x71, 0x2f, 0xca, 0x54, 0x39, 0x2f, 0x4d, 0xdd,
		0xe9, 0x37, 0x51, 0x20, 0x6b, 0x35, 0x84, 0xc8, 0xf4, 0x3f, 0x9e, 0x6d, 0xc5, 0x15, 0x31, 0xf9,
	}
	if cmp(res, test, 64) == 0 {
		fmt.Printf("test is wrong\n")
		return 0
	}
	fmt.Printf("test valid\n")
	return 1
}

//22b6837845c6bef65ea71672b265831086d3c76aebe6dae91cad51d83f79d16b
//22b6837845c6bef65ea71672b265831086d3c76aebe6dae91cad51d83f79d16b

//074c9330599d7f8d712fca54392f4ddde93751206b3584c8f43f9e6dc51531f9
//074c9330599d7f8d712fca54392f4ddde93751206b3584c8f43f9e6dc51531f9
