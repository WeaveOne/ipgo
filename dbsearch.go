package ipgo

import (
	"fmt"
	"math/big"
	"net"
	"os"
)

var HeaderSip []int64
var HeaderPtr []int64
var db *os.File
var headerLength int

type DataBlock struct {
	CityId  int
	Region  string
	dataPtr int
}

func BtreeSearch(ips string) (DataBlock, error) {
	if db == nil {
		return DataBlock{}, fmt.Errorf("Please use the GetFile method to initialize or check the path.")
	}
	ip := InetAtoN(ips)
	if HeaderSip == nil {
		_, e := db.Seek(8, 0)
		if e != nil {
			fmt.Println(e.Error())
		}
		b := make([]byte, 16384)
		_, err := db.Read(b)
		if err != nil {
			fmt.Println(err.Error())
		}
		lenn := len(b) >> 3
		idx := 0
		HeaderSip = make([]int64, lenn)
		HeaderPtr = make([]int64, lenn)
		for i := 0; i < len(b); i += 8 {
			startIp := GetIntLong(b, i)
			dataPtr := GetIntLong(b, i+4)
			if dataPtr == 0 {
				break
			}
			HeaderSip[idx] = startIp
			HeaderPtr[idx] = dataPtr
			idx++
		}
		headerLength = idx
	}
	if ip == HeaderSip[0] {
		return GetByIndexPtr(HeaderPtr[0]), nil
	} else if ip == HeaderSip[headerLength-1] {
		return GetByIndexPtr(HeaderPtr[headerLength-1]), nil
	}

	var l, h = 0, headerLength
	var sptr, eptr int64 = 0, 0
	for l <= h {
		m := (l + h) >> 1
		if ip == HeaderSip[m] {
			if m > 0 {
				sptr = HeaderPtr[m-1]
				eptr = HeaderPtr[m]
			} else {
				sptr = HeaderPtr[m]
				eptr = HeaderPtr[m+1]
			}

			break
		}

		if ip < HeaderSip[m] {
			if m == 0 {
				sptr = HeaderPtr[m]
				eptr = HeaderPtr[m+1]
				break
			} else if ip > HeaderSip[m-1] {
				sptr = HeaderPtr[m-1]
				eptr = HeaderPtr[m]
				break
			}
			h = m - 1
		} else {
			if m == headerLength-1 {
				sptr = HeaderPtr[m-1]
				eptr = HeaderPtr[m]
				break
			} else if ip <= HeaderSip[m+1] {
				sptr = HeaderPtr[m]
				eptr = HeaderPtr[m+1]
				break
			}
			l = m + 1
		}
	}
	if sptr == 0 {
		return DataBlock{}, nil
	}
	blockLen := int(eptr - sptr)
	blen := 12
	iBuffer := make([]byte, blockLen+blen)
	db.Seek(sptr, 0)
	db.Read(iBuffer)

	l = 0
	h = blockLen / blen
	var sip, eip, dataptr int64 = 0, 0, 0
	for l <= h {
		m := (l + h) >> 1
		p := m * blen
		sip = GetIntLong(iBuffer, p)
		if ip < sip {
			h = m - 1
		} else {
			eip = GetIntLong(iBuffer, p+4)
			if ip > eip {
				l = m + 1
			} else {
				dataptr = GetIntLong(iBuffer, p+8)
				break
			}
		}
	}
	if dataptr == 0 {
		return DataBlock{}, nil
	}
	dataLen := (dataptr >> 24) & 0xFF
	dataPtr := dataptr & 0x00FFFFFF
	db.Seek(dataPtr, 0)
	data := make([]byte, dataLen)
	db.Read(data)
	cityId := int(GetIntLong(data, 0))
	region := string(data[4:])
	return DataBlock{cityId, region, int(dataPtr)}, nil
}
func InetAtoN(ip string) int64 {
	ret := big.NewInt(0)
	ret.SetBytes(net.ParseIP(ip).To4())
	return ret.Int64()
}

func GetFile(path string) error {
	file, e := os.Open(path)
	if e != nil {
		return fmt.Errorf("%v please check the path", e)
	}
	db = file
	return e
}
func GetIntLong(b []byte, offset int) int64 {
	v1 := offset
	v2 := v1 + 1
	v3 := v2 + 1
	v4 := v3 + 1
	return (ByteToInt64(b[v1]) & 255) | (ByteToInt64(b[v2]) << 8 & 65280) | (ByteToInt64(b[v3]) << 16 & 16711680) | (ByteToInt64(b[v4]) << 24 & 4278190080)
}
func ByteToInt64(b byte) int64 {
	return int64(int64(b))
}
func GetByIndexPtr(ptr int64) DataBlock {
	db.Seek(ptr, 0)
	buffer := make([]byte, 12)
	db.Read(buffer)
	extra := GetIntLong(buffer, 8)
	dataLen := (extra >> 24) & 0xFF
	dataPtr := extra & 0x00FFFFFF
	db.Seek(dataPtr, 0)
	data := make([]byte, dataLen)
	db.Read(data)
	cityId := int(GetIntLong(data, 0))
	region := string(data[4:])
	return DataBlock{cityId, region, int(dataPtr)}
}
