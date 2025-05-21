package main

import (
	"encoding/hex"
	"log"
	"os"
)

func writeData(bt []byte, filename string) {
	err := os.WriteFile(filename, bt, 0644)
	if err != nil {
		log.Fatalf("Error writing to output file: %v", err)
	}
}
func readData(filename string) []byte {
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}
	t1 := make([]byte, len(data)*2)
	for i, b := range data {
		t1[i*2] = (b >> 4) & 0x0F
		t1[i*2+1] = b & 0x0F
	}
	return t1
}
func getBytes(cipher []byte) []byte {
	bt := make([]byte, len(cipher)/2)
	for i := 0; i < len(bt); i++ {
		bt[i] = (cipher[i*2] << 4) | (cipher[i*2+1] & 0xF)
	}
	return bt
}
func readParams() ([]byte, []byte) {
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
	return t2, t3
}
