package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	vsdk "cloud.google.com/go/vision/apiv1"
	"golang.org/x/xerrors"
	"google.golang.org/api/option"
	v "google.golang.org/genproto/googleapis/cloud/vision/v1"
)

func main() {
	ctx := context.Background()
	w := os.Stdout
	gcpKeyJSON, err := ioutil.ReadFile(os.Getenv("GCP_KEY_PATH"))
	if err != nil {
		fmt.Fprintf(w, "%+v", xerrors.Errorf("認証情報のjsonが読み込めませんでした。", err))
		return
	}

	client, err := vsdk.NewImageAnnotatorClient(ctx, option.WithCredentialsJSON(gcpKeyJSON))
	if err != nil {
		fmt.Fprintf(w, "%+v", xerrors.Errorf("クライアントの生成に失敗しました。", err))
		return
	}

	image, err := getImage(getImageFilePathFromArgs(os.Args))
	if err != nil {
		fmt.Fprintf(w, "%+v", xerrors.Errorf("画像が読み込めないす。", err))
		return
	}

	annotation, err := client.DetectDocumentText(ctx, image, &v.ImageContext{
		LanguageHints: []string{"ja"},
	})
	if err != nil {
		fmt.Fprintf(w, "%+v", xerrors.Errorf("画像の解釈に失敗しました。", err))
		return
	}
	outputAnnotation(w, annotation)
}

func getImage(imagePath string) (*v.Image, error) {
	imageReader, err := os.Open(imagePath)
	if err != nil {
		return nil, xerrors.Errorf("画像パスから読めないですわ。", err)
	}

	return vsdk.NewImageFromReader(imageReader)
}

func getImageFilePathFromArgs(args []string) string {
	if len(args) == 1 {
		return ""
	}
	return args[1]
}

func outputAnnotation(w io.Writer, annotation *v.TextAnnotation) {
	if annotation == nil {
		fmt.Fprintln(w, "解釈できるテキストはありませんでした。")
		return
	}
	jsonBytes, err := json.Marshal(annotation)
	if err != nil {

		fmt.Fprintf(w, "%+v", xerrors.Errorf("画像の解釈結果をJSONに変換できませんでした。", err))
	}
	fmt.Fprintln(w, string(jsonBytes))
}
