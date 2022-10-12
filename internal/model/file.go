package model

import "github.com/gogf/gf/v2/os/gtime"

// File is the golang structure for table file.
type File struct {
	Id               uint        `json:"id"        `              // FILE ID
	CreatedAt        *gtime.Time `json:"createdAt" `              // 文件创建时间
	Name             string      `json:"name"      `              // 文件名
	Md5              string      `json:"md5"       `              // 文件MD5值
	Size             int64       `json:"size"      `              // 文件大小
	StoragePath      string      `json:"storage_path"      `      // 上传文件存储目录
	UncompressedPath string      `json:"uncompressed_path"      ` // 上传文件存储目录
	Import           int         `json:"import"    `              // 文件是否被解压,默认未处理
}
