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
	client := lark.NewClient("appId", "appSecret")
	md, err := lark_docx_md.DocxMarkdown(context.Background(), client, "documentId")
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

## Related Repo

[lark_docs_md](https://github.com/chyroc/lark_docs_md)