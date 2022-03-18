// =================================================================================
// Code generated by GoFrame CLI tool. DO NOT EDIT. Created at 2022-03-18 17:58:03
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// File is the golang structure of table file for DAO operations like Where/Data.
type File struct {
	g.Meta    `orm:"table:file, do:true"`
	Id        interface{} // FILE ID
	CreatedAt *gtime.Time // 文件创建时间
	Name      interface{} // 文件名
	Md5       interface{} // 文件MD5值
	Size      interface{} // 文件大小
	Path      interface{} // 文件存储目录
	Import    interface{} // 文件是否被处理,默认未处理
}
