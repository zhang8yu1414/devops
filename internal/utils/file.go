package utils

import (
	"fmt"
	"path/filepath"
	"regexp"
)

// GetDirectoryDockerFileList 查看指定目录下的文件，返回绝对路径列表
func GetDirectoryDockerFileList(path string) []string {
	files, _ := filepath.Glob(fmt.Sprintf("%s/*.tar", path))
	return files
}

// FindIpAddress 匹配IP
func FindIpAddress(input string) string {
	partIp := "(25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])"
	balance := partIp + "\\." + partIp + "\\." + partIp + "\\." + partIp
	matchMe := regexp.MustCompile(balance)
	return matchMe.FindString(input)
}
