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

func UploadFiles(ctx context.Context, storageClient *storage.Client, project *common.Project, uid string, board common.Board) error {
	for _, ext := range board.Exts {
		if err := uploadFile(ctx, storageClient, project, uid, ext); err != nil {
			return err
		}
	}
	return nil
}

func uploadFile(ctx context.Context, storageClient *storage.Client, project *common.Project, uid string, ext string) error {
	localFilePath := filepath.Join("/app/build", project.Id, project.Id+".ino."+ext)
	file, err := os.Open(localFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	bucketName := "sketch-bridge.firebasestorage.app"
	remoteFilePath := fmt.Sprintf("build/%s/%s", uid, project.Id+"."+ext)
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
	fmt.Printf("Uploaded %s to %s.\n", localFilePath, remoteFilePath)
	return nil
}
