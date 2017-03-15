package utils

import (
	"strings"
	"fmt"
	"github.com/Unknwon/com"
	"os"
	"path"
	"github.com/sluu99/uuid"
	"time"
)

func GetUserPicSavePath(filename string, uname string) string {

	arr := strings.Split(filename, ".")
	savefile := fmt.Sprint(uname, ".", arr[len(arr) - 1])
	savePath := fmt.Sprint("static/upload/pic/", savefile)
	if !com.IsExist(savePath) {
		os.MkdirAll(path.Dir(savePath), os.ModePerm)
	}
	return savePath;
}

func GetSavePathBySize(filename string, size string) string {
	arr := strings.Split(filename, ".")
	Extension := arr[len(arr) - 1]
	return fmt.Sprint(arr[0], ".", Extension, "_", size, ".", Extension)
}
func GetSavePathArr(filename string, sizes []string, types string) (string, []string) {

	gs := make([]string, len(sizes) + 1)
	var id uuid.UUID = uuid.Rand()
	arr := strings.Split(filename, ".")

	Extension := arr[len(arr) - 1]
	hex := id.Hex()

	//a120161221_0000000001.jpg
	savefile := fmt.Sprint(hex, ".", Extension)
	t := Substr(time.Now().Format(time.RFC3339), 0, 10)
	savePath := fmt.Sprint(types, t, "/", savefile)

	//a120161221_0000000001.jpg_l.jpg
	for k, v := range sizes {
		s := fmt.Sprint(hex, ".", Extension, "_", v, ".", Extension)
		gs[k] = fmt.Sprint(types, t, "/", s)
	}
	if !com.IsExist(savePath) {
		os.MkdirAll(path.Dir(savePath), os.ModePerm)
	}
	return savePath, gs;
}

func GetSavePath(filename string, types string) string {

	var id uuid.UUID = uuid.Rand()
	arr := strings.Split(filename, ".")
	savefile := fmt.Sprint(id.Hex(), ".", arr[len(arr) - 1])
	t := Substr(time.Now().Format(time.RFC3339), 0, 10)
	savePath := fmt.Sprint(types, t, "/", savefile)
	if !com.IsExist(savePath) {
		os.MkdirAll(path.Dir(savePath), os.ModePerm)
	}
	return savePath;
}

func StringIsNotEmpty(val string) bool {
	if len(val) != 0 || val != "" {
		return true
	}
	return false
}

func Substr(str string, start, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0

	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length

	if start > end {
		start, end = end, start
	}

	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}

	return string(rs[start:end])
}

