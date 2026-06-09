package oss

import (
	"bytes"
	"fmt"
	"io"
	"super-order-web/internal/config"
	"strings"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

var client *oss.Client

// Initialize 初始化OSS客户端
func Initialize(cfg *config.OSSConfig) error {
	if cfg.Endpoint == "" || cfg.AccessKeyID == "" || cfg.AccessKeySecret == "" {
		return fmt.Errorf("OSS配置不完整")
	}

	var err error
	client, err = oss.New(cfg.Endpoint, cfg.AccessKeyID, cfg.AccessKeySecret)
	if err != nil {
		return fmt.Errorf("创建OSS客户端失败: %w", err)
	}

	return nil
}

// Upload 上传文件到OSS
func Upload(bucketName, objectName string, reader io.Reader) (string, error) {
	if client == nil {
		return "", fmt.Errorf("OSS客户端未初始化")
	}

	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return "", fmt.Errorf("获取Bucket失败: %w", err)
	}

	if err := bucket.PutObject(objectName, reader); err != nil {
		return "", fmt.Errorf("上传文件失败: %w", err)
	}

	return fmt.Sprintf("https://%s.%s/%s", bucketName, client.Config.Endpoint, objectName), nil
}

// UploadBytes 上传字节数组到OSS
func UploadBytes(bucketName, objectName string, data []byte) (string, error) {
	return Upload(bucketName, objectName, bytes.NewReader(data))
}

// Delete 删除OSS文件
func Delete(bucketName, objectName string) error {
	if client == nil {
		return fmt.Errorf("OSS客户端未初始化")
	}

	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return fmt.Errorf("获取Bucket失败: %w", err)
	}

	if err := bucket.DeleteObject(objectName); err != nil {
		return fmt.Errorf("删除文件失败: %w", err)
	}

	return nil
}

// GetObjectName 生成OSS对象名（按日期分目录）
func GetObjectName(prefix, filename string) string {
	ext := ""
	if idx := strings.LastIndex(filename, "."); idx > 0 {
		ext = filename[idx:]
		filename = filename[:idx]
	}
	return fmt.Sprintf("%s/%s%s", prefix, filename, ext)
}
