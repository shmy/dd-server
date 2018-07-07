package util

import (
	"io"
	"math/rand"
	"mime/multipart"
	"os"
	"path"
	"regexp"
	"strconv"
	"time"
)

const UPLOAD_PATH = "public/upload"

// 判断路径是否存在
func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// 创建目录
func mkDir(path string, perm os.FileMode) error {
	b, err := pathExists(path)
	if err != nil {
		return err
	}
	if !b {
		err := os.MkdirAll(path, perm)
		if err != nil {
			return err
		}
	}
	return nil
}

// 获取当前日期的路径 upload/2018/06/17
func getUploadSubPath() string {
	now := time.Now()
	year := now.Year()
	month := int(now.Month())
	day := now.Day()
	_year := strconv.Itoa(year)
	_month := strconv.Itoa(month)
	_day := strconv.Itoa(day)
	if month < 10 {
		_month = "0" + _month
	}
	if day < 10 {
		_day = "0" + _day
	}
	return path.Join(UPLOAD_PATH, _year, _month, _day)
}

// 生成随机字符串
func getRandomString(l int) string {
	str := "_-0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}
func AutoGenerateUploadDir() error {

	return mkDir(UPLOAD_PATH, 0777)
}
func SaveUploadFile(src multipart.File, name string) (string, error) {
	ext := path.Ext(name)
	basePath := getUploadSubPath()
	err := mkDir(basePath, 0777)
	if err != nil {
		return "", err
	}
	filename := path.Join(basePath, getRandomString(32)+ext)
	dst, err := os.Create(filename)
	defer dst.Close()
	if err != nil {
		return "", err
	}

	if _, err = io.Copy(dst, src); err != nil {
		return "", err
	}
	// 删掉 UPLOAD_PATH/
	// TODO windows目录换行符
	reg := regexp.MustCompile("^" + UPLOAD_PATH + "/")
	return reg.ReplaceAllString(filename, ""), nil
}
