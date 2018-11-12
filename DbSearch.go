package main

import (
	"fmt"
	"ipgo/util"
	"math/big"
	"net"
	"os"
)

var HeaderSip []int64
var HeaderPtr []int
var db *os.File
var headerLength int

func main() {
	GetFile()
	BtreeSearch(InetAtoN("183.226.61.202"))
}
func BtreeSearch(ip int64) {
	var h int
	var sptr int
	if HeaderSip == nil {
		_, e := db.Seek(8, 1)
		if e != nil {
			fmt.Println(e.Error())
		}
		b := make([]byte, 16384)
		_, err := db.Read(b)
		if err != nil {
			fmt.Println(err.Error())
		}
		h = len(b) >> 3
		sptr = 0
		HeaderSip = make([]int64, h)
		HeaderPtr = make([]int, h)
		for i := 0; i < len(b); i += 8 {
			startIp := GetIntLong(b, i)
			dataPtr := GetIntLong(b, i+4)
			if dataPtr == 0 {
				break
			}
			HeaderSip[sptr] = startIp
			HeaderPtr[sptr] = int(dataPtr)
			sptr++
		}
		headerLength = sptr

		fmt.Println(ip)
		fmt.Println(db.Name())
		fmt.Println(b)
		fmt.Println(HeaderSip)
		fmt.Println(HeaderPtr)
	}
	if ip == HeaderSip[0] {

	}
}
func InetAtoN(ip string) int64 {
	ret := big.NewInt(0)
	ret.SetBytes(net.ParseIP(ip).To4())
	return ret.Int64()
}

func GetFile() {
	file, e := os.Open("d:/ip2region.db")
	if e != nil {

	}
	db = file
}
func GetIntLong(b []byte, offset int) int64 {
	v1 := offset
	v2 := v1 + 1
	v3 := v2 + 1
	v4 := v3 + 1
	fmt.Println(v1,v2,v3,v4)
	return (util.ByteToInt64(b[v1]) & 255) | (util.ByteToInt64(b[v2]) << 8 & 65280) | (util.ByteToInt64(b[v3]) << 16 & 16711680) | (util.ByteToInt64(b[v4]) << 24 & 4278190080)
}
