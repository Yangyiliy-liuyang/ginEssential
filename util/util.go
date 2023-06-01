package util

import (
	"math/rand"
	"time"
)

func RandomString(n int) string {
	var letters = []byte("assadsasadfsadsaGdsfsFFSWLKDKSKDAXDF")
	result := make([]byte, n)
	//使用时间作为种子值，然后生成不同系列的随机数
	rand.Seed(time.Now().Unix())
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}
