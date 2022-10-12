package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
)

type sUpdate struct {
}

var (
	inUpdate = sUpdate{}
)

func Update() *sUpdate {
	return &inUpdate
}

func (s *sUpdate) UpdateApp(ctx context.Context, fileName string) (err error) {
	fileRecord, _ := g.Model("file").Fields("storage_path, uncompressed_path, import").Where("name=", fileName).One()
	record := fileRecord.GMap()
	//g.Log().Print(ctx, record)
	var outImagesPath string

	// @todo: 最好的逻辑是解压覆盖原有目录，需删除解压目录，以及update mysql数据
	// 如果重复点击，文件已经解压过，则不解压
	if record.Get("import") == 0 && record.Get("uncompressed_path") == "" {
		updateFileUncompressedDestPathVar, _ := g.Config().Get(ctx, "update.uncompressedPath")
		updateFileUncompressedDestPath := updateFileUncompressedDestPathVar.String()

		uncompressedOutPath := File().UncompressedFile(fmt.Sprintf("%s/%s", record.Get("storage_path"), fileName), updateFileUncompressedDestPath)
		g.Log().Infof(ctx, "%s文件解压到%s", fileName, uncompressedOutPath)

		imageCompressFile := fmt.Sprintf("%s/%s", uncompressedOutPath, "images.tar.gz")
		outImagesPath = File().UncompressedFile(imageCompressFile, uncompressedOutPath)
		g.Log().Infof(ctx, "%s文件解压到%s", imageCompressFile, outImagesPath)

		_, _ = g.Model("file").Data(g.Map{"import": 1, "uncompressed_path": uncompressedOutPath}).Where("name=", fileName).Update()
	} else if record.Get("import") == 1 && record.Get("uncompressed_path") == "" ||
		record.Get("import") == 0 && record.Get("uncompressed_path") != "" {
		errInfo := errors.New("文件已经不正常解压，出现脏数据")
		g.Log().Error(ctx, errInfo)
		return errInfo
	} else {
		g.Log().Infof(ctx, "%s已经解压，不需重复解压", fileName)
	}

	// 查询解压目录
	imagesFullPath, _ := g.Model("file").Fields("uncompressed_path").Where("name=", fileName).Value()

	err = Docker().LoadImageAndPushToHarbor(ctx, fileName, imagesFullPath.String())
	return err
}
