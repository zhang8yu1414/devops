package service

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"os"
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
	path, err := g.Config().Get(ctx, "docker.filePath")
	dirPath := path.String()
	filePath := fmt.Sprintf("%s/%s", dirPath, "docker")
	err = os.MkdirAll(filePath, os.ModePerm)
	if err != nil {
		return err
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
	_, err = g.Model("file").Where("name=", filename).Delete()
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
