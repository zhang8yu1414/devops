package utils

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/os/glog"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// GetDirectoryDockerFileList 查看docker目录下的文件，返回绝对路径列表
func GetDirectoryDockerFileList(path string) []string {
	files, _ := filepath.Glob(fmt.Sprintf("%s/*.tar", path))
	return files
}

// GetDirectoryUpdateAndFullFileList 查看app目录下不同的文件
func GetDirectoryUpdateAndFullFileList(path string) []string {
	//fileList := make([]string, 10)
	files, _ := filepath.Glob(path)
	for _, file := range files {
		fileStr := filepath.Base(file)
		if prefix := strings.HasPrefix(fileStr, "update"); prefix {
			fileInfo, err := os.Stat(file)
			if err != nil {
				glog.Debugf(context.Background(), "%s文件获取状态失败", file)
			}
			modTime := fileInfo.ModTime()
			fmt.Println("modTime", modTime)
		}
	}
	return files
}

// FindIpAddress 匹配IP
func FindIpAddress(input string) string {
	partIp := "(25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])"
	balance := partIp + "\\." + partIp + "\\." + partIp + "\\." + partIp
	matchMe := regexp.MustCompile(balance)
	return matchMe.FindString(input)
}
