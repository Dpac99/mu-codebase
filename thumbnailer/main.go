package main

import (
	"bytes"
	"context"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/disintegration/imaging"
)

type Request struct {
	Input_bucket  string `json:"input_bucket"`
	Output_bucket string `json:"output_bucket"`
	Input_key     string `json:"input_key"`
	Output_key    string `json:"output_key"`
	Width         int    `json:"width"`
	Height        int    `json:"height"`
}

func HandleRequest(ctx context.Context, r Request) (int, error) {
	var region = os.Getenv("AWS_REGION")
	log.Println(region)

	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("eu-west-3")}))
	downloader := s3manager.NewDownloader(sess)

	// var r Request
	// log.Println(req.Body)
	// json.Unmarshal([]byte(req.Body), &r)
	// log.Println(r)

	buff := &aws.WriteAtBuffer{}
	_, err := downloader.Download(buff, &s3.GetObjectInput{
		Bucket: aws.String(r.Input_bucket),
		Key:    aws.String(r.Input_key),
	})
	if err != nil {
		return -1, err
	}

	thumb_reader := resizeImage(buff.Bytes(), r.Width, r.Height)

	uploader := s3manager.NewUploader(sess)

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(r.Output_bucket),
		Key:    aws.String(r.Output_key),
		Body:   thumb_reader,
	})
	if err != nil {
		return -1, err
	}

	return 1, nil

}

func resizeImage(image_bytes []byte, w int, h int) *bytes.Reader {
	src, err := imaging.Decode(bytes.NewReader(image_bytes))
	if err != nil {
		log.Fatalf(err.Error())
	}

	thumbnail := imaging.Thumbnail(src, w, h, imaging.CatmullRom)
	dst := imaging.New(w, h, color.NRGBA{0, 0, 0, 0})
	dst = imaging.Paste(dst, thumbnail, image.Pt(0, 0))
	buf := new(bytes.Buffer)
	err = png.Encode(buf, dst)
	if err != nil {
		log.Fatalf(err.Error())
	}

	return bytes.NewReader(buf.Bytes())
}

func main() {
	lambda.Start(HandleRequest)
}
