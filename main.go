package main

import (
	"context"
	"fmt"
	"html"
	"io"

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
	docId := "12345"

	// GoogleドキュメントAPIクライアントを作成
	client, err := docs.NewService(context.Background(), option.WithCredentialsFile("credential.json"))
	if err != nil {
		fmt.Println(err)
		return
	}

	html, err := NewGoogleDocument(client, docId).GetHTML()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(html)
}
