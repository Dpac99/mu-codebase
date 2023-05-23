package tasks

import (
	"bytes"
	"encoding/json"
	"errors"
	"image"
	"image/color"
	"image/png"
	"log"
	"serverless/types"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/disintegration/imaging"
)

type CPURequest struct {
	N       int         `json:"n"`
	Vectors [][]float64 `json:"vectors"`
}

type SleepRequest struct {
	Duration int `json:"duration"`
}

type MatrixRequest struct {
	A [][]int64 `json:"a"`
	B [][]int64 `json:"b"`
}

type ThumbnailRequest struct {
	Input_bucket  string `json:"input_bucket"`
	Output_bucket string `json:"output_bucket"`
	Input_key     string `json:"input_key"`
	Output_key    string `json:"output_key"`
	Width         int    `json:"width"`
	Height        int    `json:"height"`
}

func ExecuteTask(req types.TaskRequest) (interface{}, error) {
	switch req.Type {
	case "sleep":
		duration := req.Args["duration"].(float64)
		log.Println(duration)
		ret := placeHolderSleep(float64(duration))

		return ret, nil

	case "cpu":
		var r CPURequest
		jsonData, _ := json.Marshal(req.Args)
		json.Unmarshal(jsonData, &r)
		// r.N = int(req.Args["n"].(float64))
		// r.Vectors = req.Args["vectors"].([][]float64)
		ret := CPUintensive(&r)

		return ret, nil

	case "matrix":
		var r MatrixRequest
		jsonData, _ := json.Marshal(req.Args)
		json.Unmarshal(jsonData, &r)
		ret := multiplyMatrix(r.A, r.B)
		return ret, nil

	case "test":
		name := req.Args["name"].(string)
		ret := testfunction(name)
		log.Println(ret)

		return ret, nil
	case "thumbnail":
		var r ThumbnailRequest
		jsonData, _ := json.Marshal(req.Args)
		json.Unmarshal(jsonData, &r)
		ret, err := HandleThumbnailRequest(r)
		return ret, err
	default:
		return nil, errors.New("unrecognized request")
	}
}

func multiplyMatrix(a [][]int64, b [][]int64) [][]int64 {
	n_rows := len(a)
	n_cols := len(b[0])
	n_elems := len(b)
	c := make([][]int64, n_rows)
	for i := range c {
		c[i] = make([]int64, n_cols)
	}

	for i := 0; i < n_rows; i++ {
		for j := 0; j < n_cols; j++ {
			c[i][j] = 0
			for k := 0; k < n_elems; k++ {
				c[i][j] += a[i][k] * b[k][j]
			}
		}
	}

	return c
}

func placeHolderSleep(duration float64) int {

	time.Sleep(time.Duration(duration * float64(time.Second)))

	return 1
}

func CPUintensive(req *CPURequest) (result float64) {
	result = 0
	for i := 0; i < req.N; i += 2 {
		result += dotProduct(req.Vectors[i], req.Vectors[i+1])
	}
	return
}

func dotProduct(a []float64, b []float64) (result float64) {
	k := len(a)
	for i := 0; i < k; i++ {
		result += a[i] * b[i]
	}
	return
}

func testfunction(name string) string {
	return "hello " + name + "!"
}

func HandleThumbnailRequest(r ThumbnailRequest) (int, error) {

	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("eu-west-3")}))
	downloader := s3manager.NewDownloader(sess)

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
