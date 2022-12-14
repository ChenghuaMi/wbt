package util

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
)
const (
	prefixCode = "mch"
	OrderLength = 6
)
func Md5(str string,pre string) string {
	if len(pre) == 0 {
		pre = prefixCode
	}
	h := md5.New()
	newStr := str + prefixCode
	fmt.Println(newStr)
	h.Write([]byte(newStr))
	byt := h.Sum(nil)
	ss := hex.EncodeToString(byt)
	fmt.Println("ss:",ss)
	return ss
}

func CreateRand(length int) string {
	str := "0123456789"
	bytes := []byte(str)
	result := []byte{}
	rand.Seed(time.Now().UnixNano()+ int64(rand.Intn(100)))
	for i := 0; i < length; i++ {
		result = append(result, bytes[rand.Intn(len(bytes))])
	}
	return string(result)
}
func CreateOrder(length int) string {
	dt := fmt.Sprintf("%s",time.Now().Format("20060102150405"))
	rd := CreateRand(length)
	return dt + rd
}

func JsonData(data interface{}) ([]byte,error) {
	bt,err := json.Marshal(data)
	return bt,err
}

func UnJsonData(data []byte,ret interface{}) error {
	return json.Unmarshal(data,&ret)
}
func ByteToStr(bt byte) string {
	bbt := []byte{bt}
	return string(bbt)
}