package utilities

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func ArrayToString(a interface{}, delim string) string {
	return strings.Trim(strings.Replace(fmt.Sprint(a), " ", delim, -1), "[]")
}

func StringToStringArray(array string, delim string) []string {
	s := strings.Split(array, delim)
	return s
}

func StringToIntArray(array string, delim string) []int {
	s := strings.Split(array, delim)
	arrayInt := []int{}
	for _, v := range s {
		value, err := strconv.Atoi(v)
		if err == nil {
			arrayInt = append(arrayInt, value)
		}
	}
	return arrayInt
}

func StringToInt32Array(array string, delim string) []int32 {
	s := strings.Split(array, delim)
	arrayInt := []int32{}
	for _, v := range s {
		value, err := strconv.ParseInt(v, 10, 32)
		if err == nil {
			arrayInt = append(arrayInt, int32(value))
		}
	}
	return arrayInt
}

func StringToInt64Array(array string, delim string) []int64 {
	s := strings.Split(array, delim)
	arrayInt := []int64{}
	for _, v := range s {
		value, err := strconv.Atoi(v)
		if err == nil {
			arrayInt = append(arrayInt, int64(value))
		}
	}
	return arrayInt
}

const charset = "abcdefghijklmnopqrstuvwxyz" + "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func RandStringBytes(length int) string {
	return StringWithCharset(length, charset)
}
