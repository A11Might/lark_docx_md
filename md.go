package lark_docx_md

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkdocx "github.com/larksuite/oapi-sdk-go/v3/service/docx/v1"
	larkdrive "github.com/larksuite/oapi-sdk-go/v3/service/drive/v1"
)

const (
	Page     = 1
	Text     = 2
	Heading1 = 3
	Heading2 = 4
	Heading3 = 5
	Heading4 = 6
	Heading5 = 7
	Heading6 = 8
	Heading7 = 9
	Heading8 = 10
	Heading9 = 11
	Bullet   = 12
	Ordered  = 13
	Code     = 14
	Image    = 27
)

func DocxMarkdown(ctx context.Context, client *lark.Client, documentId string) (string, error) {
	req := larkdocx.NewListDocumentBlockReqBuilder().DocumentId(documentId).Build()
	iterator, _ := client.Docx.V1.DocumentBlock.ListByIterator(ctx, req)

	var (
		hasMore = true
		block   *larkdocx.Block
		err     error

		buf = new(strings.Builder)
	)
	for ; hasMore && err == nil; hasMore, block, err = iterator.Next() {
		buf.WriteString(DocxBlockMarkdown(ctx, client, block))
		buf.WriteString("\n")
	}

	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func DocxBlockMarkdown(ctx context.Context, client *lark.Client, block *larkdocx.Block) string {
	if block == nil {
		return ""
	}

	switch *block.BlockType {
	case Page:
		return BlockPageMarkdown(ctx, block)
	case Text:
		return BlockTextMarkdown(ctx, block)
	case Heading1, Heading2, Heading3, Heading4, Heading5, Heading6, Heading7, Heading8, Heading9:
		return BlockHeadingMarkdown(ctx, block)
	case Bullet:
		return BlockBulletMarkdown(ctx, block)
	case Ordered:
		return BlockOrderedMarkdown(ctx, block)
	case Code:
		return BlockCodeMarkdown(ctx, block)
	case Image:
		return BlockImageMarkdown(ctx, client, block)
	default:
		return fmt.Sprintf("not support block type:%d", *block.BlockType)
	}
}

func BlockPageMarkdown(ctx context.Context, block *larkdocx.Block) string {
	return fmt.Sprintf("# %s", TextMarkdown(ctx, block.Page))
}

// BlockTextMarkdown 处理文本 Block
func BlockTextMarkdown(ctx context.Context, block *larkdocx.Block) string {
	return TextMarkdown(ctx, block.Text)
}

func BlockHeadingMarkdown(ctx context.Context, block *larkdocx.Block) string {
	var heading *larkdocx.Text
	if block.Heading1 != nil {
		heading = block.Heading1
	} else if block.Heading2 != nil {
		heading = block.Heading2
	} else if block.Heading3 != nil {
		heading = block.Heading3
	} else if block.Heading4 != nil {
		heading = block.Heading4
	} else if block.Heading5 != nil {
		heading = block.Heading5
	} else if block.Heading6 != nil {
		heading = block.Heading6
	} else if block.Heading7 != nil {
		heading = block.Heading7
	} else if block.Heading8 != nil {
		heading = block.Heading8
	} else if block.Heading9 != nil {
		heading = block.Heading9
	}

	// block type: [3, 11] -> heading: [1, 9]
	return fmt.Sprintf("%s %s", strings.Repeat("#", *block.BlockType-2), TextMarkdown(ctx, heading))
}

func BlockBulletMarkdown(ctx context.Context, block *larkdocx.Block) string {
	return fmt.Sprintf("- %s", TextMarkdown(ctx, block.Bullet))
}

func BlockOrderedMarkdown(ctx context.Context, block *larkdocx.Block) string {
	return fmt.Sprintf("1 %s", TextMarkdown(ctx, block.Bullet))
}

const (
	Go   = 22
	JSON = 28
)

var lmap = map[int]string{
	Go:   "go",
	JSON: "json",
}

func BlockCodeMarkdown(ctx context.Context, block *larkdocx.Block) string {
	return fmt.Sprintf("```%s\n%s\n```", lmap[*block.Code.Style.Language], TextMarkdown(ctx, block.Code))
}

func TextMarkdown(ctx context.Context, text *larkdocx.Text) string {
	buf := new(strings.Builder)

	for _, e := range text.Elements {
		if e.TextRun != nil { // 文字
			buf.WriteString(TextAddElementStyle(*e.TextRun.Content, e.TextRun.TextElementStyle))
		} else if e.MentionDoc != nil { // @文档
			buf.WriteString(fmt.Sprintf("[%s](%s)", TextAddElementStyle(*e.MentionDoc.Title, e.TextRun.TextElementStyle), *e.MentionDoc.Url))
		}
	}

	return buf.String()
}

func TextAddElementStyle(text string, style *larkdocx.TextElementStyle) string {
	if *style.Bold {
		text = "**" + text + "**"
	}
	if *style.InlineCode {
		text = "`" + text + "`"
	}
	if *style.Italic {
		text = "*" + text + "*"
	}
	if style.Link != nil {
		text = "[" + text + "]" + "(" + *style.Link.Url + ")"
	}
	if *style.Strikethrough {
		text = "~~" + text + "~~"
	}
	if *style.Underline {
		text = "<u>" + text + "</u>"
	}

	return text
}

func BlockImageMarkdown(ctx context.Context, client *lark.Client, block *larkdocx.Block) string {
	req := larkdrive.NewDownloadMediaReqBuilder().
		FileToken(*block.Image.Token).
		Build()
	resp, err := client.Drive.Media.Download(ctx, req)
	if err != nil {
		log.Printf("lark download drive media %s fail: %s", *block.Image.Token, err)
		return ""
	}
	filename := fmt.Sprintf("%s/%s", "static", *block.Image.Token+".jpg")
	_ = os.MkdirAll(filepath.Dir(filename), 0o755)
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0o666)
	if err != nil {
		log.Printf("open file %s fail: %s", filename, err)
		return ""
	}
	defer f.Close()

	_, _ = io.Copy(f, resp.File)
	return fmt.Sprintf("<img src=%q width=\"%d\" height=\"%d\"/>", filename, *block.Image.Width, *block.Image.Height)
}
