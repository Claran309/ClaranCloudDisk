package minIO

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime"
	"os"
	"path"
	"path/filepath"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.uber.org/zap"
)

type MinIOClient struct {
	Client     *minio.Client
	BucketName string
}

func NewMinIOClient(endpoint, accessKeyID, secretAccessKey, bucketName, DefaultAvatarPath string) (*MinIOClient, error) {
	//初始化minIOClient
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: false, // http
	})
	if err != nil {
		return nil, fmt.Errorf("初始化minIO客户端失败: %v", err)
	}

	//验证bucket是否存在
	exist, err := minioClient.BucketExists(context.Background(), bucketName)
	if !exist {
		err = minioClient.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return nil, fmt.Errorf("新建bucket失败: %v", err)
		}
	}
	if err != nil {
		return nil, fmt.Errorf("验证minIO_bucket失败: %v", err)
	}

	//上传默认头像
	Path := filepath.Join(".", DefaultAvatarPath)
	DefaultAvatar, err := os.Open(Path)
	if err != nil {
		fmt.Printf("获取默认头像失败:%v\n", err)
	}
	DefaultAvatarData, _ := io.ReadAll(DefaultAvatar)
	ext := path.Ext(DefaultAvatarPath)
	reader := bytes.NewReader(DefaultAvatarData)
	size := int64(len(DefaultAvatarData))
	mimeType := mime.TypeByExtension(ext)
	opts := minio.PutObjectOptions{ContentType: mimeType}
	_, err = minioClient.PutObject(context.Background(), bucketName, DefaultAvatarPath, reader, size, opts)
	if err != nil {
		zap.S().Errorf("上传默认头像失败:%v\n", err)
		fmt.Printf("上传默认头像失败:%v\n", err)
	}

	return &MinIOClient{
		Client:     minioClient,
		BucketName: bucketName,
	}, nil
}

func (m *MinIOClient) Save(ctx context.Context, objectName string, data []byte, ext string) error {
	//捕获数据
	//mimeType := mime.TypeByExtension(ext)
	mimeType := mime.TypeByExtension(ext)
	reader := bytes.NewReader(data)
	size := int64(len(data))
	//Options
	opts := minio.PutObjectOptions{ContentType: mimeType}

	//保存文件
	//PutObject(ctx context.Context, bucketName, objectName string, reader io.Reader, objectSize int64,opts PutObjectOptions) (info UploadInfo, err error)
	_, err := m.Client.PutObject(ctx, m.BucketName, objectName, reader, size, opts)
	if err != nil {
		zap.S().Errorf("保存到minIO失败: %v", err)
		return fmt.Errorf("保存到minIO失败: %v", err)
	}

	return nil
}

func (m *MinIOClient) Delete(ctx context.Context, objectName string) error {
	//RemoveObject(ctx context.Context, bucketName, objectName string, opts minio.RemoveObjectOptions) error
	err := m.Client.RemoveObject(ctx, m.BucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		zap.S().Errorf("从minIO删除文件失败: %v", err)
		return fmt.Errorf("从minIO删除文件失败: %v", err)
	}

	return nil
}

func (m *MinIOClient) Update(ctx context.Context, objectName string, data []byte, ext string) error {
	return m.Save(ctx, objectName, data, ext)
}

func (m *MinIOClient) Exists(ctx context.Context, objectName string) (bool, error) {
	//StatObject(ctx context.Context, bucketName, objectName string, opts StatObjectOptions) (ObjectInfo, error)
	//NoSuchKey
	_, err := m.Client.StatObject(ctx, m.BucketName, objectName, minio.StatObjectOptions{})
	if err != nil {
		if errType, ok := err.(minio.ErrorResponse); ok && errType.Code == "NoSuchKey" {
			return false, nil
		}
		return false, fmt.Errorf("检查文件失败: %v", err)
	}
	return true, nil
}

func (m *MinIOClient) GetBytes(ctx context.Context, objectName string) ([]byte, error) {
	//GetObject(ctx context.Context, bucketName, objectName string, opts GetObjectOptions) (*Object, error)
	obj, err := m.Client.GetObject(ctx, m.BucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("获取文件失败: %v", err)
	}
	defer obj.Close()

	//检查文件是否存在
	if _, err := obj.Stat(); err != nil {
		return nil, fmt.Errorf("文件不存在: %v", err)
	}

	//读取数据
	data, err := io.ReadAll(obj)
	if err != nil {
		return nil, fmt.Errorf("读取文件失败: %v", err)
	}

	return data, nil
}

func (m *MinIOClient) GetStream(ctx context.Context, objectName string) (io.ReadCloser, error) { //GetObject(ctx context.Context, bucketName, objectName string, opts GetObjectOptions) (*Object, error)
	obj, err := m.Client.GetObject(ctx, m.BucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("获取文件失败: %v", err)
	}
	//流式传输不关闭obj，直接返回

	return obj, nil
}
