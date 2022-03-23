package service

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/glog"
	"os"
	"path"
	"zhangyudevops.com/internal/model/model/entity"
	"zhangyudevops.com/internal/utils"
)

var inFile = sFile{}

type sFile struct{}

func File() *sFile {
	return &inFile
}

func (s *sFile) UploadFile(ctx context.Context, inFile *ghttp.UploadFile) (err error) {
	var file = &entity.File{}
	file.Size = inFile.Size
	filePath := ""

	// 查找配置文件所配置的文件上传路径
	configFilePath, err := g.Config().Get(ctx, "file.filePath")
	dirPath := configFilePath.String()

	// 判断上传文件是什么格式，放在不同的目录下
	inFileName := inFile.Filename
	if suffix := path.Ext(inFileName); suffix == ".gz" || suffix == ".tgz" || suffix == ".xz" {
		glog.Debugf(ctx, "%s为应用升级压缩包或者全量压缩包，放在app目录", inFileName)
		filePath = fmt.Sprintf("%s/%s", dirPath, "app")
		err = os.MkdirAll(filePath, os.ModePerm)
		if err != nil {
			return err
		}
	} else if suffix = path.Ext(inFileName); suffix == ".tar" {
		// .tar文件包包放在docker目录下
		glog.Debugf(ctx, "%s为docker image包，放在docker目录下", inFileName)
		filePath = fmt.Sprintf("%s/%s", dirPath, "docker")
		err = os.MkdirAll(filePath, os.ModePerm)
		if err != nil {
			return err
		}
	} else {
		glog.Debugf(ctx, "%s为垃圾收集站，放在trash目录下", inFileName)
		filePath = fmt.Sprintf("%s/%s", dirPath, "trash")
		err = os.MkdirAll(filePath, os.ModePerm)
		if err != nil {
			return err
		}
	}

	file.Path = filePath
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
		id, _ := tx.Ctx(ctx).Model("file").Fields("id").Where("name=", filename).Value()
		_, err = tx.Ctx(ctx).Model("image").Where("file_id", id).Delete()
		if err != nil {
			return err
		}

		_, err = tx.Ctx(ctx).Model("file").Where("name=", filename).Delete()
		if err != nil {
			return err
		}

		return nil
	})
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

	return nil
}
