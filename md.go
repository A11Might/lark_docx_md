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
	"github.com/samber/lo"
)

func (p *DocxMarkdownProcessor) DocxBlockMarkdown(ctx context.Context, root *Node) (md string) {
	if root == nil || root.Block == nil {
		return ""
	}
	curBlock := root.Block

	// 先处理子块
	var subBlockTexts []string
	for _, childNode := range root.ChildrenNode {
		subBlockTexts = append(subBlockTexts, p.DocxBlockMarkdown(ctx, childNode))
	}

	// 再处理父块
	switch *curBlock.BlockType {
	case Page:
		md = p.BlockPageMarkdown(ctx, curBlock)
	case Text:
		md = p.BlockTextMarkdown(ctx, curBlock)
	case Heading1, Heading2, Heading3, Heading4, Heading5, Heading6, Heading7, Heading8, Heading9:
		md = p.BlockHeadingMarkdown(ctx, curBlock)
	case Bullet:
		md = p.BlockBulletMarkdown(ctx, curBlock)
	case Ordered:
		md = p.BlockOrderedMarkdown(ctx, curBlock)
	case Code:
		md = p.BlockCodeMarkdown(ctx, curBlock)
	case Quote:
		md = p.BlockQuoteMarkdown(ctx, curBlock)
	case Todo:
		md = p.BlockTodoMarkdown(ctx, curBlock)
	case Callout:
		return p.BlockCalloutMarkdown(ctx, curBlock, subBlockTexts)
	case Divider:
		md = p.BlockDividerMarkdown(ctx)
	case Image:
		md = p.BlockImageMarkdown(ctx, curBlock)
	case TableCell:
		return strings.Join(subBlockTexts, "")
	case Table:
		return p.BlockTableMarkdown(ctx, curBlock, subBlockTexts)
	case QuoteContainer:
		return p.BlockQuoteContainerMarkdown(ctx, subBlockTexts)
	default:
		md = fmt.Sprintf("<!-- not support block type %d -->", *curBlock.BlockType)
	}

	// 合并父块和子块
	return md + "\n\n" + strings.Join(subBlockTexts, "")
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

	// block type: [3, 11] -> heading: [1, 9] -> markdown [1, 6]
	cnt := *block.BlockType - 2
	cnt = lo.Ternary(cnt > 6, 6, cnt)
	return strings.Repeat("#", cnt) + " " + p.TextMarkdown(ctx, heading)
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
		return "- [x] " + p.TextMarkdown(ctx, block.Todo)
	}
	return "- [ ] " + p.TextMarkdown(ctx, block.Todo)
}

func (p *DocxMarkdownProcessor) BlockCalloutMarkdown(ctx context.Context, block *larkdocx.Block, subBlockTexts []string) string {
	buf := new(strings.Builder)

	emoji := emojiMap[*block.Callout.EmojiId]
	text := strings.Join(FixTexts(subBlockTexts), "")

	// 没有颜色转成普通文本
	if block.Callout.BackgroundColor == nil {
		// 加入高亮块 emoji
		buf.WriteString(emoji)
		buf.WriteString(" ")
		buf.WriteString(text)
		buf.WriteString("\n\n")
		return buf.String()
	}

	if p.UseGhCallout && backgroundColorMap[*block.Callout.BackgroundColor] != "" {
		buf.WriteString("> ")
		buf.WriteString(backgroundColorMap[*block.Callout.BackgroundColor])
		buf.WriteString("\n>\n")
	}

	buf.WriteString("> ")
	// 加入高亮块 emoji
	buf.WriteString(emoji)
	buf.WriteString(" ")
	buf.WriteString(text)
	buf.WriteString("\n>\n\n")

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
			// return fmt.Sprintf("<img src=%q width=\"%d\" height=\"%d\"/>", *v.TmpDownloadUrl, *block.Image.Width, *block.Image.Height)
			return fmt.Sprintf("![%s](%s)", *v.FileToken, *v.TmpDownloadUrl)
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
		// return fmt.Sprintf("<img src=%q width=\"%d\" height=\"%d\"/>", mdname, *block.Image.Width, *block.Image.Height)
		return fmt.Sprintf("![%s](%s)", name, mdname)
	}
}

func (p *DocxMarkdownProcessor) BlockTableMarkdown(ctx context.Context, block *larkdocx.Block, subBlockTexts []string) string {
	buf := new(strings.Builder)

	/**
	| Tables        | Are           | Cool  |
	| ------------- |:-------------:| -----:|
	| col 3 is      | right-aligned | $1600 |
	| col 2 is      | centered      |   $12 |
	| zebra stripes | are neat      |    $1 |
	*/
	rows := lo.FromPtr(block.Table.Property.RowSize)
	cols := lo.FromPtr(block.Table.Property.ColumnSize)
	for row := 0; row < rows; row++ {
		if row == 1 {
			// TODO 处理表头，以第一行的样式作为表格样式
			buf.WriteString(strings.ReplaceAll(strings.Repeat("|:-:|", cols), "||", "|"))
			buf.WriteString("\n")
		}
		buf.WriteString("|")
		for col := 0; col < cols; col++ {
			buf.WriteString(FixText(subBlockTexts[row*cols+col]))
			buf.WriteString("|")
		}
		buf.WriteString("\n")
	}

	return buf.String()
}

func (p *DocxMarkdownProcessor) BlockQuoteContainerMarkdown(ctx context.Context, subBlockTexts []string) string {
	buf := new(strings.Builder)

	subBlockTexts = FixTexts(subBlockTexts)
	for _, line := range subBlockTexts {
		buf.WriteString("> ")
		buf.WriteString(line)
		buf.WriteString("\n>\n\n")
	}

	return buf.String()
}
