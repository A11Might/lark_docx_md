package main

import (
	"context"
	"fmt"
	"os"

	"github.com/A11Might/lark_docx_md"
	lark "github.com/larksuite/oapi-sdk-go/v3"
)

func main() {
	client := lark.NewClient("appId", "appSecret")
	md, err := lark_docx_md.DocxMarkdown(context.Background(), client, "U3hXdQmMAoiNVSxDgPOcu4R8nTd")
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
