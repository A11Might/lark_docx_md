# lark_docx_md

Parse Lark/Feishu Docx into Markdown

## Install

```
go get github.com/A11Might/lark_docx_md
```

## Usage
```go
package main

import (
	"context"
	"fmt"

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
		fmt.Println(err)
		return
	}
	fmt.Println(md)
}

```

## Example

Origin lark docx：[docx](https://r5q4tiv935.feishu.cn/docx/U3hXdQmMAoiNVSxDgPOcu4R8nTd)

Parse into Markdown：[md](./example.md)

## Related repo

[A11Might/lark-docx-readme](https://github.com/A11Might/lark-docx-readme): Use lark docx update github README.md

*Inspired by [chyroc/lark_docs_md](https://github.com/chyroc/lark_docs_md)*