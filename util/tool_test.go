package util

import (
	"fmt"
	"testing"
)

func TestMd5(t *testing.T) {
	str := "abc"
	Md5(str,"")
}

func TestCreateOrder(t *testing.T) {
	for i := 0;i<10;i++ {
		ss := CreateOrder(6)
		fmt.Println("ss:",ss)
	}
}