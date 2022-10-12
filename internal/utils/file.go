package utils

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"fmt"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"io"
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

// FindIpAddress 匹配IP
func FindIpAddress(input string) string {
	partIp := "(25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])"
	balance := partIp + "\\." + partIp + "\\." + partIp + "\\." + partIp
	matchMe := regexp.MustCompile(balance)
	return matchMe.FindString(input)
}

// ExtraTarGzip 解压.tar.gz 文件
func ExtraTarGzip(tarFile, destPath string) (err error, uncompressedFullPath string) {
	// 打开压缩文件流
	srcFile, err := os.Open(tarFile)
	if err != nil {
		return err, ""
	}
	defer srcFile.Close()

	gr, err := gzip.NewReader(srcFile)
	if err != nil {
		return err, ""
	}
	defer gr.Close()

	// tar read
	tr := tar.NewReader(gr)

	var outFullPath string
	for true {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err, ""
		}

		outFullPath = fmt.Sprintf("%s/%s", destPath, header.Name)

		g.Log().Debugf(context.Background(), "%s文件被解压", outFullPath)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.Mkdir(outFullPath, 0755); err != nil {
				return err, ""
			}
		case tar.TypeReg:
			outFile, err := os.Create(outFullPath)
			if err != nil {
				return err, ""
			}
			if _, err := io.Copy(outFile, tr); err != nil {
				return err, ""
			}
			outFile.Close()
		default:
			err = gerror.New("未知文件类型")
			return err, ""
		}
	}
	suffixPath := strings.TrimPrefix(outFullPath, destPath)
	uncompressedPath := strings.Split(suffixPath, "/")[1]
	uncompressedFullPath = fmt.Sprintf("%s/%s", destPath, uncompressedPath)
	return nil, uncompressedFullPath
}
