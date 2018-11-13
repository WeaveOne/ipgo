package ipgo

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	file, e := New("d:/ip2region.db")
	defer file.Close()
	if e != nil {
		fmt.Println(e.Error())
	}
	dataBlock, e := BtreeSearch("171.221.151.72")
	if e != nil {
		fmt.Println(e.Error())
	}
	fmt.Println(dataBlock)
}
