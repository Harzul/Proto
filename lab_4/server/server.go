package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

var VERSION = "0"
var CS = "111000"
var K []uint8 = []uint8{
	0xff, 0xee, 0xdd, 0xcc,
	0xbb, 0xaa, 0x99, 0x88,
	0x77, 0x66, 0x55, 0x44,
	0x33, 0x22, 0x11, 0x00,
	0xf0, 0xf1, 0xf2, 0xf3,
	0xf4, 0xf5, 0xf6, 0xf7,
	0xf8, 0xf9, 0xfa, 0xfb,
	0xfc, 0xfd, 0xfe, 0xff,
}
var iv = []byte{
	0x12, 0x34, 0x56, 0x78, 0x90, 0xab, 0xcd, 0xef,
	0x23, 0x45, 0x67, 0x89, 0x0a, 0xbc, 0xde, 0xf1,
}
var recievedMessages = make([]bool, SeqNum)

func handler(w http.ResponseWriter, r *http.Request) {
	file, logger := initLogger()
	logger.Println("Программа запущена")
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		logger.Fatalln("Error, Only POST allowed")
		return
	}

	var msg Message
	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		logger.Fatalln("Error, Invalid JSON")
		return
	}
	s, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}
	if len(s) > 2048 {
		http.Error(w, "Invalid JSON len", http.StatusBadRequest)
		logger.Fatalln("Error, Invalid JSON len")
	}
	if msg.Header.Version == VERSION && msg.Header.CS == CS {
		val, _ := strconv.Atoi(msg.Header.SeqNum)
		if val < 0 {
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte("Wrong SeqNum"))
			logger.Fatalln("Error, wrong SeqNum")
			return
		}

		Snum, err := strconv.Atoi(msg.Header.SeqNum)
		if err != nil {
			fmt.Println("Conversion error:", err)
			return
		}
		if (uint64(Snum) < SeqNum) && (recievedMessages[Snum]) {
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte("Repeated SeqNum"))
			logger.Fatalln("Error, Repeated SeqNum")
			return
		}
		if uint64(Snum) > SeqNum {
			ttt := make([]bool, Snum)
			copy(ttt, recievedMessages)
			ttt[Snum-1] = true
			recievedMessages = ttt
		}

		size := 256 * (val) //key_param
		res := make([]uint8, size/8)
		kdf_tree(res, K, 32, 1, size)
		tmp := res[(val-1)*32 : (val-1)*32+32]
		hash, err := hex.DecodeString(msg.ICV)
		if err != nil {
			logger.Fatalln("Error while decoding")
			return
		}
		msg.ICV = ""
		data, err := json.Marshal(msg)
		if err != nil {
			logger.Fatalln("Error while decoding")
			return
		}
		h := make([]uint8, 32)
		get256(data, h)

		if !bytes.Equal(hash[:], h[:]) {
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte("Wrong ICV"))
			logger.Fatalln("Error, wrong ICV")
			return
		}
		data, err = hex.DecodeString(msg.PayloadData)
		if err != nil {
			logger.Fatalln("Error while decoding")
			return
		}

		t1 := make([]byte, len(data)*2)
		for i, b := range data {
			t1[i*2] = (b >> 4) & 0x0F
			t1[i*2+1] = b & 0x0F
		}
		t2 := make([]byte, len(iv)*2)
		for i, b := range iv {
			t2[i*2] = (b >> 4) & 0x0F
			t2[i*2+1] = b & 0x0F
		}
		t3 := make([]byte, len(tmp)*2)
		for i, b := range tmp {
			t3[i*2] = (b >> 4) & 0x0F
			t3[i*2+1] = b & 0x0F
		}
		openText := magic(t1, t2, t3)
		logger.Printf("Получено сообщение:\n%+v\n", string(getBt(openText)))
	} else {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte("Wrong Version of CS"))
		logger.Fatalln("Error, wrong CS")
	}
	file.Close()
}

func main() {
	http.HandleFunc("/submit", handler)
	fmt.Println("Сервер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
