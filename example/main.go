package main

import (
	"context"
	"fmt"
	"os"

	"github.com/A11Might/lark_docx_md"
	lark "github.com/larksuite/oapi-sdk-go/v3"
)

func main() {
	processor := lark_docx_md.NewDocxMarkdownProcessor(
		lark.NewClient("appId", "appSecret"),
		"documentType", "documentToken",
		lark_docx_md.DownloadStatic("static", "static"),
		lark_docx_md.UseGhCalloutStyle(),
	)
	md, err := processor.DocxMarkdown(context.Background())
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	f, err := os.Create("example.md")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer f.Close()
	if _, err := f.WriteString(md); err != nil {
		fmt.Println(err.Error())
		return
	}
}
