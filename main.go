package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/storage"
)

var (
	bucket string
	file   string
	key    string

	accessKey string //ACCESS_KEY
	secretKey string //SECRET_KEY
	zone      string //zone

	upToken      string
	formUploader *storage.FormUploader
	ret          storage.PutRet
	putExtra     storage.PutExtra
)

func init() {
	flag.StringVar(&bucket, "bucket", "", "upload qiniu bucket")
	flag.StringVar(&file, "file", "", "upload file name")
	flag.StringVar(&key, "key", "", "bucket key")
}

func initQiniu() {
	accessKey = os.Getenv("QINIU_ACCESS_KEY") //ACCESS_KEY
	secretKey = os.Getenv("QINIU_SECRET_KEY") //SECRET_KEY
	zone = os.Getenv("QINIU_ZONE")            //zone

	putPolicy := storage.PutPolicy{
		Scope: bucket,
	}

	mac := qbox.NewMac(accessKey, secretKey)
	upToken = putPolicy.UploadToken(mac)

	cfg := storage.Config{}
	switch zone {
	case "Huadong":
		cfg.Zone = &storage.ZoneHuadong
	case "Huabei":
		cfg.Zone = &storage.ZoneHuabei
	case "Huanan":
		cfg.Zone = &storage.ZoneHuanan
	case "Beimei":
		cfg.Zone = &storage.ZoneBeimei
	}
	cfg.UseHTTPS = false

	formUploader = storage.NewFormUploader(&cfg)
	ret = storage.PutRet{}
	putExtra = storage.PutExtra{}
}

func main() {
	flag.Parse()
	initQiniu()
	if bucket == "" || key == "" || file == "" {
		flag.PrintDefaults()
	}

	f, err := os.Open(file)
	if err != nil {
		log.Printf("Open file error:%s", err.Error())
		return
	}
	defer f.Close()

	err = formUploader.Put(context.Background(), &ret, upToken, key+"/"+file, f, fileSize(file), &putExtra)
	if err != nil {
		log.Printf("upload file err: %s, fileName: %s", err, file)
		return
	}
	log.Printf("has upload file %s sucess\n", file)
}

func fileSize(path string) int64 {
	fileInfo, err := os.Stat(path)
	if err != nil {
		panic(err)
	}
	fileSize := fileInfo.Size()
	return fileSize
}
