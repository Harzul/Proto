package main

type Splitmix64 struct {
	S uint64
}

func (x *Splitmix64) Next() uint64 {
	(*x).S += 0x9e3779b97f4a7c15
	var z uint64 = x.S
	z = (z ^ (z >> 30)) * 0xbf58476d1ce4e5b9
	z = (z ^ (z >> 27)) * 0x94d049bb133111eb
	return z ^ (z >> 31)
}

type Xoshiro256_PP struct {
	S [4]uint64
}

func (x Xoshiro256_PP) rotl(a uint64, b int) uint64 {
	return (a << b) | (a >> (64 - b))
}

func (x *Xoshiro256_PP) Next() uint64 {
	var result uint64 = x.rotl(x.S[0]+x.S[3], 23) + x.S[0]

	var t uint64 = x.S[1] << 17

	(*x).S[2] ^= x.S[0]
	(*x).S[3] ^= x.S[1]
	(*x).S[1] ^= x.S[2]
	(*x).S[0] ^= x.S[3]

	(*x).S[2] ^= t

	(*x).S[3] = x.rotl(x.S[3], 45)

	return result
}
