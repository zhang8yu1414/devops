package utils

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
)

func Md5f(filename string) (md string, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()
	md5h := md5.New()
	_, err = io.Copy(md5h, f)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(md5h.Sum(nil)), nil
}
