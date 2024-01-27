package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/A11Might/lark_docx_md"
	lark "github.com/larksuite/oapi-sdk-go/v3"
)

func main() {
	client := lark.NewClient(os.Getenv("APP_ID"), os.Getenv("APP_SECRET"))
	md, err := lark_docx_md.DocxMarkdown(context.Background(), client, os.Getenv("DOCUMENT_ID"))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	filename := "dist/README.md"
	_ = os.MkdirAll(filepath.Dir(filename), 0o755)
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0o666)
	if err != nil {
		fmt.Printf("open file %s fail: %s", filename, err)
		return
	}
	defer f.Close()
	if _, err := f.WriteString(md); err != nil {
		fmt.Println(err.Error())
		return
	}

	if err = os.Rename("static", "dist/static"); err != nil {
		fmt.Println(err.Error())
		return
	}
}
