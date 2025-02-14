package database

import (
	"context"
	"firebase.google.com/go/v4/storage"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sketch-bridge/arduino-compile-server/common"
)

func UploadHexFile(ctx context.Context, storageClient *storage.Client, project *common.Project, uid string) error {
	localFilePath := filepath.Join("/app/build", project.Id, project.Id+".ino.hex")
	file, err := os.Open(localFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	bucketName := "sketch-bridge.firebasestorage.app"
	remoteFilePath := fmt.Sprintf("build/%s/%s", uid, project.Id+".hex")
	bucket, err := storageClient.Bucket(bucketName)
	if err != nil {
		return err
	}
	writer := bucket.Object(remoteFilePath).NewWriter(ctx)
	if _, err := io.Copy(writer, file); err != nil {
		return err
	}
	if err := writer.Close(); err != nil {
		return err
	}
	return nil
}
