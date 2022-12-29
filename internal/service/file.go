package service

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"os"
	"path"
	"strings"
	"time"
	"zhangyudevops.com/internal/model"
	"zhangyudevops.com/internal/utils"
)

var inFile = sFile{}

type sFile struct{}

func File() *sFile {
	return &inFile
}

// UploadFile 文件上传
// 总行升级包以center_为前缀开始
// 分行升级包以division_前缀开始
func (s *sFile) UploadFile(ctx context.Context, inFile *ghttp.UploadFile) (err error) {
	var file = &model.File{}
	file.Size = inFile.Size
	filePath := ""

	// 查找配置文件所配置的文件上传路径
	configFilePath, _ := g.Config().Get(ctx, "file.filePath")
	dirPath := configFilePath.String()

	// 判断上传文件是什么格式，放在不同的目录下
	inFileName := inFile.Filename
	if suffix := path.Ext(inFileName); suffix == ".gz" || suffix == ".tgz" || suffix == ".xz" {
		if prefix := strings.HasPrefix(inFileName, "center"); prefix {
			g.Log().Debugf(ctx, "%s为总行应用升级压缩包，放在app/center目录", inFileName)
			filePath = fmt.Sprintf("%s/%s", dirPath, "app/center")
			err = os.MkdirAll(filePath, os.ModePerm)
			if err != nil {
				return err
			}
		} else if prefix = strings.HasPrefix(inFileName, "division"); prefix {
			g.Log().Debugf(ctx, "%s为分行应用全量安装压缩包，放在app/division目录", inFileName)
			filePath = fmt.Sprintf("%s/%s", dirPath, "app/division")
			err = os.MkdirAll(filePath, os.ModePerm)
			if err != nil {
				return err
			}
		} else {
			g.Log().Debugf(ctx, "%s为不知名压缩包，放在app/trash目录", inFileName)
			filePath = fmt.Sprintf("%s/%s", dirPath, "app/trash")
			err = os.MkdirAll(filePath, os.ModePerm)
			if err != nil {
				return err
			}
		}
	} else if suffix = path.Ext(inFileName); suffix == ".tar" {
		// .tar文件包包放在docker目录下
		g.Log().Debugf(ctx, "%s为docker image包，放在docker目录下", inFileName)
		filePath = fmt.Sprintf("%s/%s", dirPath, "docker")
		err = os.MkdirAll(filePath, os.ModePerm)
		if err != nil {
			return err
		}
	} else {
		g.Log().Debugf(ctx, "%s为垃圾收集站，放在trash目录下", inFileName)
		filePath = fmt.Sprintf("%s/%s", dirPath, "trash")
		err = os.MkdirAll(filePath, os.ModePerm)
		if err != nil {
			return err
		}
	}

	file.StoragePath = filePath
	filename, err := inFile.Save(filePath, false)
	if err != nil {
		return err
	}

	file.Name = filename
	md5, err := utils.Md5f(fmt.Sprintf("%s/%s", filePath, filename))
	if err != nil {
		return err
	}
	file.Md5 = md5
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
		result, _ := tx.Ctx(ctx).Model("file").Fields("id,md5").Where("name=", filename).One()
		r := result.GMap()
		if md5Record, _ := tx.Ctx(ctx).Model("file").Where("md5=", r.Get("md5")).One(); md5Record.IsEmpty() != true {
			l, _ := time.LoadLocation("Asia/Shanghai")
			now := time.Now().In(l).Format("2006-01-02 15:04:05")
			_, _ = tx.Ctx(ctx).Model("file").Data("created_at", now).Where("md5=", r.Get("md5")).Update()
		} else {
			_, err = tx.Ctx(ctx).Model("image").Where("file_id", r.Get("id")).Delete()
			if err != nil {
				return err
			}

			_, err = tx.Ctx(ctx).Model("file").Where("name=", filename).Delete()
			if err != nil {
				return err
			}

			err = g.DB().Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
				_, err = tx.Ctx(ctx).Model("file").Data(file).Insert()
				if err != nil {
					return err
				}
				return nil
			})
			if err != nil {
				return err
			}
		}

		return nil
	})

	return nil
}

// UncompressedFile 解压.tar.gz包
// @args: tarFile 压缩包绝对路径
// @args: destPath 解压后的路径
// @return: 解压后的路径
func (s *sFile) UncompressedFile(tarFile string, destPath string) (outPath string) {

	_, outPath = utils.ExtraTarGzip(tarFile, destPath)
	return outPath
}
