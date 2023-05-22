package main

import (
	"context"
	"fmt"
	"html"
	"io"
	"os"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/docs/v1"
	"google.golang.org/api/option"
)

type Document interface {
	GetBody() (io.Reader, error)
	GetHTML() (string, error)
}

type GoogleDocument struct {
	client *docs.Service
	docId  string
}

func NewGoogleDocument(client *docs.Service, docId string) *GoogleDocument {
	return &GoogleDocument{
		client: client,
		docId:  docId,
	}
}

func (d *GoogleDocument) GetBody() ([]byte, error) {
	doc, err := d.client.Documents.Get(d.docId).Do()
	if err != nil {
		return nil, err
	}
	return doc.Body.MarshalJSON()
}

func (d *GoogleDocument) GetHTML() (string, error) {
	body, err := d.GetBody()
	if err != nil {
		return "", err
	}
	return html.EscapeString(string(body)), nil
}

func main() {
	docId := "1234567890"

	// クレデンシャル取得
	credentials, err := getCredentials()
	if err != nil {
		fmt.Println(err)
		return
	}

	// GoogleドキュメントAPIクライアントを作成
	client, err := docs.NewService(context.Background(), option.WithCredentials(credentials))
	if err != nil {
		fmt.Println(err)
		return
	}

	html, err := NewGoogleDocument(client, docId).GetHTML()
	if err != nil {
		fmt.Println(err)
		return
	}

	// HTMLを出力します。
	fmt.Println(html)
}

func getCredentials() (*google.Credentials, error) {
	f, err := os.Open("credentials.json")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// io.Readerから読み込んだjsonを[]byteに変換
	bytes, err := io.ReadAll(f)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	credentials, err := google.CredentialsFromJSON(context.Background(), bytes)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return credentials, nil
}
