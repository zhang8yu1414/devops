package service

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"os"
)

type sUpdate struct {
}

var (
	inUpdate = sUpdate{}
)

func Update() *sUpdate {
	return &inUpdate
}

func (s *sUpdate) UncompressedTarFileAndPushHarbor(ctx context.Context, fileName string) (err error) {
	fileRecord, _ := g.Model("file").Fields("storage_path, uncompressed_path, import").Where("name=", fileName).One()
	record := fileRecord.GMap()
	//g.Log().Print(ctx, record)
	var outImagesPath string

	// @todo: 最好的逻辑是解压覆盖原有目录，需删除解压目录，以及update mysql数据
	// 如果重复解压，删除之前解压的目录
	uncompressedPathVar := record.Get("uncompressed_path")
	uncompressedPath := gconv.String(uncompressedPathVar)
	if uncompressedPath != "" {
		g.Log().Info(ctx, "压缩文件已经解压过，覆盖原解压目录")
		err = os.RemoveAll(uncompressedPath)
		if err != nil {
			g.Log().Errorf(ctx, "删除原解压目录失败:%s", err)
			return err
		}
	}
	updateFileUncompressedDestPathVar, _ := g.Config().Get(ctx, "update.uncompressedPath")
	updateFileUncompressedDestPath := updateFileUncompressedDestPathVar.String()

	uncompressedOutPath := File().UncompressedFile(fmt.Sprintf("%s/%s", record.Get("storage_path"), fileName), updateFileUncompressedDestPath)
	g.Log().Infof(ctx, "%s文件解压到%s", fileName, uncompressedOutPath)

	imageCompressFile := fmt.Sprintf("%s/%s", uncompressedOutPath, "images.tar.gz")
	outImagesPath = File().UncompressedFile(imageCompressFile, uncompressedOutPath)
	g.Log().Infof(ctx, "%s文件解压到%s", imageCompressFile, outImagesPath)

	_, _ = g.Model("file").Data(g.Map{"import": 1, "uncompressed_path": uncompressedOutPath}).Where("name=", fileName).Update()

	// 查询解压目录
	imagesFullPath, _ := g.Model("file").Fields("uncompressed_path").Where("name=", fileName).Value()

	err = Docker().LoadImageAndPushToHarbor(ctx, fileName, imagesFullPath.String())
	return err
}
