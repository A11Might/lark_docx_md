package lark_docx_md

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	larkdocx "github.com/larksuite/oapi-sdk-go/v3/service/docx/v1"
	larkdrive "github.com/larksuite/oapi-sdk-go/v3/service/drive/v1"
)

func (p *DocxMarkdownProcessor) DocxBlockMarkdown(ctx context.Context, item *ProcessItem) (md string) {
	if item == nil {
		return ""
	}

	switch *item.BlockType {
	case Page:
		md = p.BlockPageMarkdown(ctx, item.Normal)
	case Text:
		md = p.BlockTextMarkdown(ctx, item.Normal)
	case Heading1, Heading2, Heading3, Heading4, Heading5, Heading6, Heading7, Heading8, Heading9:
		md = p.BlockHeadingMarkdown(ctx, item.Normal)
	case Bullet:
		md = p.BlockBulletMarkdown(ctx, item.Normal)
	case Ordered:
		md = p.BlockOrderedMarkdown(ctx, item.Normal)
	case Code:
		md = p.BlockCodeMarkdown(ctx, item.Normal)
	case Quote:
		md = p.BlockQuoteMarkdown(ctx, item.Normal)
	case Todo:
		md = p.BlockTodoMarkdown(ctx, item.Normal)
	case Callout:
		md = p.BlockCalloutMarkdown(ctx, item.Callout)
	case Divider:
		md = p.BlockDividerMarkdown(ctx)
	case Image:
		md = p.BlockImageMarkdown(ctx, item.Normal)
	case Table:
		md = p.BlockTableMarkdown(ctx, item.Table)
	case QuoteContainer:
		md = p.BlockQuoteContainerMarkdown(ctx, item.QuoteContainer)
	default:
		md = fmt.Sprintf("<!-- not support block type %d -->", *item.BlockType)
	}

	return md + "\n\n"
}

func (p *DocxMarkdownProcessor) BlockPageMarkdown(ctx context.Context, block *larkdocx.Block) string {
	return "# " + p.TextMarkdown(ctx, block.Page)
}

func (p *DocxMarkdownProcessor) BlockTextMarkdown(ctx context.Context, block *larkdocx.Block) string {
	return p.TextMarkdown(ctx, block.Text)
}

func (p *DocxMarkdownProcessor) BlockHeadingMarkdown(ctx context.Context, block *larkdocx.Block) string {
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
	return strings.Repeat("#", *block.BlockType-2) + " " + p.TextMarkdown(ctx, heading)
}

func (p *DocxMarkdownProcessor) BlockBulletMarkdown(ctx context.Context, block *larkdocx.Block) string {
	return "- " + p.TextMarkdown(ctx, block.Bullet)
}

func (p *DocxMarkdownProcessor) BlockOrderedMarkdown(ctx context.Context, block *larkdocx.Block) string {
	return "1. " + p.TextMarkdown(ctx, block.Ordered)
}

func (p *DocxMarkdownProcessor) BlockCodeMarkdown(ctx context.Context, block *larkdocx.Block) string {
	return fmt.Sprintf("```%s\n%s\n```", languageMap[*block.Code.Style.Language], p.TextMarkdown(ctx, block.Code))
}

func (p *DocxMarkdownProcessor) BlockQuoteMarkdown(ctx context.Context, block *larkdocx.Block) string {
	return "> " + p.TextMarkdown(ctx, block.Quote)
}

func (p *DocxMarkdownProcessor) BlockTodoMarkdown(ctx context.Context, block *larkdocx.Block) string {
	if *block.Todo.Style.Done {
		return "[x] " + p.TextMarkdown(ctx, block.Todo)
	}
	return "[] " + p.TextMarkdown(ctx, block.Todo)
}

func (p *DocxMarkdownProcessor) BlockCalloutMarkdown(ctx context.Context, blocks []*larkdocx.Block) string {
	buf := new(strings.Builder)

	emoji := emojiMap[*blocks[0].Callout.EmojiId]

	// 没有颜色转成普通文本
	if blocks[0].Callout.BackgroundColor == nil {
		// 加入高亮块 emoji
		buf.WriteString(emoji)
		buf.WriteString(" ")
		for _, b := range blocks[1:] {
			buf.WriteString(p.BlockTextMarkdown(ctx, b))
			buf.WriteString("\n\n")
		}
		return buf.String()
	}

	if p.UseGhCallout {
		buf.WriteString("> ")
		buf.WriteString(backgroundColorMap[*blocks[0].Callout.BackgroundColor])
		buf.WriteString("\n>\n")
	}

	for i, b := range blocks[1:] {
		buf.WriteString("> ")
		// 加入高亮块 emoji
		if i == 0 {
			buf.WriteString(emoji)
			buf.WriteString(" ")
		}
		buf.WriteString(p.TextMarkdown(ctx, b.Text))
		buf.WriteString("\n>\n")
	}

	return buf.String()
}

func (p *DocxMarkdownProcessor) TextMarkdown(ctx context.Context, text *larkdocx.Text) string {
	buf := new(strings.Builder)

	preStyle := larkdocx.NewTextElementStyleBuilder().Bold(false).InlineCode(false).Italic(false).Strikethrough(false).Underline(false).Build()
	for _, e := range text.Elements {
		// 将链接和@文档都转成普通文字处理
		textRun := e.TextRun
		if textRun != nil && textRun.TextElementStyle.Link != nil {
			textRun = larkdocx.NewTextRunBuilder().
				Content(fmt.Sprintf("[%s](%s)", *textRun.Content, UnescapeUrl(*textRun.TextElementStyle.Link.Url))).
				TextElementStyle(textRun.TextElementStyle).
				Build()
		}
		if e.MentionDoc != nil {
			textRun = larkdocx.NewTextRunBuilder().
				Content(fmt.Sprintf("[%s](%s)", *e.MentionDoc.Title, UnescapeUrl(*e.MentionDoc.Url))).
				TextElementStyle(e.MentionDoc.TextElementStyle).
				Build()
		}
		// 处理文本
		if textRun != nil {
			// 相邻文本样式相同则统一加样式，不同则开启新样式
			if *textRun.TextElementStyle.Bold != *preStyle.Bold ||
				*textRun.TextElementStyle.InlineCode != *preStyle.InlineCode ||
				*textRun.TextElementStyle.Italic != *preStyle.Italic ||
				*textRun.TextElementStyle.Strikethrough != *preStyle.Strikethrough ||
				*textRun.TextElementStyle.Underline != *preStyle.Underline {
				// 结束上一个样式
				if *preStyle.Bold {
					buf.WriteString("**")
				}
				if *preStyle.InlineCode {
					buf.WriteString("`")
				}
				if *preStyle.Italic {
					buf.WriteString("*")
				}
				if *preStyle.Strikethrough {
					buf.WriteString("~~")
				}
				if *preStyle.Underline {
					buf.WriteString("</u>")
				}
				// 开启下一个样式
				if *textRun.TextElementStyle.Bold {
					buf.WriteString("**")
				}
				if *textRun.TextElementStyle.InlineCode {
					buf.WriteString("`")
				}
				if *textRun.TextElementStyle.Italic {
					buf.WriteString("*")
				}
				if *textRun.TextElementStyle.Strikethrough {
					buf.WriteString("~~")
				}
				if *textRun.TextElementStyle.Underline {
					buf.WriteString("<u>")
				}
			}
			buf.WriteString(*textRun.Content)
			preStyle = textRun.TextElementStyle
		}
	}
	// 最后一个文本的结束样式
	if *preStyle.Bold {
		buf.WriteString("**")
	}
	if *preStyle.InlineCode {
		buf.WriteString("`")
	}
	if *preStyle.Italic {
		buf.WriteString("*")
	}
	if *preStyle.Strikethrough {
		buf.WriteString("~~")
	}
	if *preStyle.Underline {
		buf.WriteString("</u>")
	}

	return buf.String()
}

func (p *DocxMarkdownProcessor) BlockDividerMarkdown(ctx context.Context) string {
	return "---"
}

func (p *DocxMarkdownProcessor) BlockImageMarkdown(ctx context.Context, block *larkdocx.Block) string {
	if p.StaticAsURL {
		// 创建请求对象
		req := larkdrive.NewBatchGetTmpDownloadUrlMediaReqBuilder().
			FileTokens([]string{*block.Image.Token}).
			Build()
		// 发起请求
		resp, err := p.LarkClient.Drive.V1.Media.BatchGetTmpDownloadUrl(ctx, req)
		if err != nil {
			log.Printf("lark get drive media tmp url %s fail: %s", *block.Image.Token, err)
			return ""
		}
		if !resp.Success() {
			log.Printf("lark get drive media tmp url %s fail: code:%d, msg:%s, requestId:%s", *block.Image.Token, resp.Code, resp.Msg, resp.RequestId())
			return ""
		}
		for _, v := range resp.Data.TmpDownloadUrls {
			return fmt.Sprintf("<img src=%q width=\"%d\" height=\"%d\"/>", *v.TmpDownloadUrl, *block.Image.Width, *block.Image.Height)
		}
		return ""
	} else {
		req := larkdrive.NewDownloadMediaReqBuilder().
			FileToken(*block.Image.Token).
			Build()
		resp, err := p.LarkClient.Drive.Media.Download(ctx, req)
		if err != nil {
			log.Printf("lark download drive media %s fail: %s", *block.Image.Token, err)
			return ""
		}
		if !resp.Success() {
			log.Printf("lark download drive media %s fail: code:%d, msg:%s, requestId:%s", *block.Image.Token, resp.Code, resp.Msg, resp.RequestId())
			return ""
		}
		name := *block.Image.Token + ".jpg"
		filename := fmt.Sprintf("%s/%s", p.StaticDir, name)
		mdname := fmt.Sprintf("%s/%s", p.FilePrefix, name)
		_ = os.MkdirAll(filepath.Dir(filename), 0o755)
		f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0o666)
		if err != nil {
			log.Printf("open file %s fail: %s", filename, err)
			return ""
		}
		defer f.Close()

		_, _ = io.Copy(f, resp.File)
		return fmt.Sprintf("<img src=%q width=\"%d\" height=\"%d\"/>", mdname, *block.Image.Width, *block.Image.Height)
	}
}

func (p *DocxMarkdownProcessor) BlockTableMarkdown(ctx context.Context, table [][]*larkdocx.Block) string {
	buf := new(strings.Builder)

	/**
	| Tables        | Are           | Cool  |
	| ------------- |:-------------:| -----:|
	| col 3 is      | right-aligned | $1600 |
	| col 2 is      | centered      |   $12 |
	| zebra stripes | are neat      |    $1 |
	*/
	// 处理表头，以第一行的样式作为表格样式
	header := new(strings.Builder)
	header.WriteString("|")
	buf.WriteString("|")
	for _, col := range table[0] {
		switch *col.Text.Style.Align {
		case AlignLeft:
			header.WriteString("-|")
		case AlignMid:
			header.WriteString(":-:|")
		case AlignRight:
			header.WriteString("-:|")
		default:
			header.WriteString("-|")
		}
		buf.WriteString(p.TextMarkdown(ctx, col.Text))
		buf.WriteString("|")
	}
	buf.WriteString("\n")
	buf.WriteString(header.String())
	buf.WriteString("\n")

	if len(table) > 1 {
		// 跳过表头
		table = table[1:]
	}
	for _, row := range table {
		buf.WriteString("|")
		for _, col := range row {
			buf.WriteString(p.TextMarkdown(ctx, col.Text))
			buf.WriteString("|")
		}
		buf.WriteString("\n")
	}

	return buf.String()
}

func (p *DocxMarkdownProcessor) BlockQuoteContainerMarkdown(ctx context.Context, blocks []*larkdocx.Block) string {
	buf := new(strings.Builder)

	for _, b := range blocks {
		buf.WriteString("> ")
		buf.WriteString(p.TextMarkdown(ctx, b.Text))
		buf.WriteString("\n>\n")
	}

	return buf.String()
}
