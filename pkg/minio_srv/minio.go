package minio_srv

import (
	"aapi/config"
	"context"
	"io"
	"log"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func NewMinIOClient(c *config.Config) (*minio.Client, error) {
	// Initialize minio client object.
	log.Printf("Init MiniIO cline witn id = %s pass = %s\n", c.MinIO.AccessKeyID, c.MinIO.SecretAccessKey)

	minioClient, err := minio.New(c.MinIO.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(c.MinIO.AccessKeyID, c.MinIO.SecretAccessKey, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	log.Printf("%#v\n", minioClient)

	return minioClient, nil
}

func NewAndroidMinIOClient(c *config.Config) (*minio.Client, error) {
	// Initialize minio client object.
	log.Printf("Init NewAndroidMinIOClient cline witn id = %s pass = %s\n", c.AndroidMinIO.AccessKeyID, c.AndroidMinIO.SecretAccessKey)

	minioClient, err := minio.New(c.AndroidMinIO.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(c.AndroidMinIO.AccessKeyID, c.AndroidMinIO.SecretAccessKey, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	log.Printf("%#v\n", minioClient)

	return minioClient, nil
}

func CreateMinIOBucket(minioClient *minio.Client, bucket string) bool {
	ctx := context.Background()

	if minioClient == nil {
		log.Println("minioClient is nil.")
		return false
	}

	err := minioClient.MakeBucket(ctx, bucket, minio.MakeBucketOptions{Region: "us-east-1"})
	if err != nil {
		exists, errBucketExists := minioClient.BucketExists(ctx, bucket)
		if errBucketExists == nil && exists {
			log.Printf("We already own %s\n", bucket)
			return true
		} else {
			log.Fatalln(err)
			return false
		}
	} else {
		log.Printf("Successfully created %s\n", bucket)
		return true
	}

}

func UploadFile(minioClient *minio.Client, bucket string, target_file_name string, file_reader io.Reader, objectSize int64) string {
	if minioClient == nil {
		log.Println("minioClient is nil.")
		return target_file_name
	}

	uploadInfo, err := minioClient.PutObject(context.Background(), bucket, target_file_name, file_reader, objectSize, minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		log.Println(err)
		return target_file_name
	}
	log.Println("Successfully uploaded bytes: ", uploadInfo)

	return GetTmpURL(minioClient, bucket, target_file_name)
}

func DownloadFile(minioClient *minio.Client, bucket string, target_file_name string) []byte {
	if minioClient == nil {
		log.Println("minioClient is nil.")
		return nil
	}

	file_object, err := minioClient.GetObject(context.Background(), bucket, target_file_name, minio.GetObjectOptions{})
	if err != nil {
		log.Println(err)
		return nil
	}
	log.Println("Successfully downloaded file: ", target_file_name)
	fileInfo, _ := file_object.Stat()
	file_bytes := make([]byte, fileInfo.Size)
	file_object.Read(file_bytes)

	return file_bytes
}

func GetTmpURL(minioClient *minio.Client, bucket string, file_name string) string {
	if minioClient == nil {
		log.Println("minioClient is nil")
		return file_name
	}

	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition", "attachment; filename=\""+file_name+"\"")
	// Generates a presigned url which expires in a day.
	presignedURL, err := minioClient.PresignedGetObject(context.Background(), bucket, file_name, time.Hour*24*7, reqParams)
	if err != nil {
		log.Println(err)
		return file_name
	}
	return presignedURL.String()
}

func GetSignedURL(minioClient *minio.Client, bucket string, file_name string, expireHour time.Duration) string {
	if minioClient == nil {
		log.Println("minioClient is nil")
		return file_name
	}

	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition", "attachment; filename=\""+file_name+"\"")

	// Generates a presigned url which expires in a day.
	presignedURL, err := minioClient.PresignedGetObject(context.Background(), bucket, file_name, expireHour, reqParams)
	if err != nil {
		log.Println(err)
		return file_name
	}
	return presignedURL.String()
}

func UploadFileOnly(minioClient *minio.Client, bucket string, target_file_name string, file_reader io.Reader, objectSize int64) string {
	if minioClient == nil {
		log.Println("minioClient is nil.")
		return target_file_name
	}

	uploadInfo, err := minioClient.PutObject(context.Background(), bucket, target_file_name, file_reader, objectSize, minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		log.Println(err)
		return target_file_name
	}
	log.Println("Successfully uploaded bytes: ", uploadInfo)

	return uploadInfo.Key
}
