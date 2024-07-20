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

func (p *DocxMarkdownProcessor) DocxBlockMarkdown(ctx context.Context, root *Node) (texts []string) {
	if root == nil || root.Block == nil {
		return nil
	}
	curBlock := root.Block

	// 先处理子块
	var subBlockTexts []string
	for _, childNode := range root.ChildrenNode {
		subBlockTexts = append(subBlockTexts, p.DocxBlockMarkdown(ctx, childNode)...)
	}

	// 再处理父块
	parentText := ""
	switch *curBlock.BlockType {
	case Page:
		parentText = p.BlockPageMarkdown(ctx, curBlock)
	case Text:
		parentText = p.BlockTextMarkdown(ctx, curBlock)
	case Heading1, Heading2, Heading3, Heading4, Heading5, Heading6, Heading7, Heading8, Heading9:
		parentText = p.BlockHeadingMarkdown(ctx, curBlock)
	case Bullet:
		parentText = p.BlockBulletMarkdown(ctx, curBlock)
	case Ordered:
		parentText = p.BlockOrderedMarkdown(ctx, curBlock)
	case Code:
		return p.BlockCodeMarkdown(ctx, curBlock)
	case Quote:
		parentText = p.BlockQuoteMarkdown(ctx, curBlock)
	case Todo:
		parentText = p.BlockTodoMarkdown(ctx, curBlock)
	case Callout:
		return p.BlockCalloutMarkdown(ctx, curBlock, subBlockTexts)
	case Divider:
		parentText = p.BlockDividerMarkdown(ctx)
	case Image:
		parentText = p.BlockImageMarkdown(ctx, curBlock)
	case TableCell:
		return subBlockTexts
	case Table:
		return p.BlockTableMarkdown(ctx, curBlock, subBlockTexts)
	case QuoteContainer:
		return p.BlockQuoteContainerMarkdown(ctx, subBlockTexts)
	default:
		parentText = fmt.Sprintf("<!-- not support block type %d -->", *curBlock.BlockType)
	}

	// 合并父块和子块
	var tmp []string
	if parentText != "" {
		tmp = append(tmp, parentText)
	}
	for _, text := range subBlockTexts {
		switch *curBlock.BlockType {
		case Page, Heading1, Heading2, Heading3, Heading4, Heading5, Heading6, Heading7, Heading8, Heading9:
			tmp = append(tmp, text)
		default:
			tmp = append(tmp, "    "+text)
		}
	}
	return tmp
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

func (p *DocxMarkdownProcessor) BlockCodeMarkdown(ctx context.Context, block *larkdocx.Block) (texts []string) {
	texts = append(texts, fmt.Sprintf("```%s", languageMap[*block.Code.Style.Language]))
	texts = append(texts, strings.Split(p.TextMarkdown(ctx, block.Code, true), "\n")...)
	texts = append(texts, "```")
	return FixTexts(texts)
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

func (p *DocxMarkdownProcessor) BlockCalloutMarkdown(ctx context.Context, block *larkdocx.Block, subBlockTexts []string) (texts []string) {
	// 加入高亮块 emoji
	emoji := emojiMap[*block.Callout.EmojiId]
	subBlockTexts[0] = emoji + " " + subBlockTexts[0]

	// 没有颜色转成普通文本
	if block.Callout.BackgroundColor == nil {
		return subBlockTexts
	}

	if p.UseGhCallout && backgroundColorMap[*block.Callout.BackgroundColor] != "" {
		texts = append(texts, fmt.Sprintf("> %s\n>", backgroundColorMap[*block.Callout.BackgroundColor]))
	}

	for _, text := range subBlockTexts {
		texts = append(texts, "> "+text+"\n>")
	}

	return FixTexts(texts)
}

func (p *DocxMarkdownProcessor) TextMarkdown(ctx context.Context, text *larkdocx.Text, withoutStyle1 ...bool) string {
	if text == nil {
		return ""
	}

	withoutStyle := false
	if len(withoutStyle1) != 0 {
		withoutStyle = withoutStyle1[0]
	}

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
			if !withoutStyle &&
				(*textRun.TextElementStyle.Bold != *preStyle.Bold ||
					*textRun.TextElementStyle.InlineCode != *preStyle.InlineCode ||
					*textRun.TextElementStyle.Italic != *preStyle.Italic ||
					*textRun.TextElementStyle.Strikethrough != *preStyle.Strikethrough ||
					*textRun.TextElementStyle.Underline != *preStyle.Underline) {
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
	if withoutStyle {
		return buf.String()
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

func (p *DocxMarkdownProcessor) BlockTableMarkdown(ctx context.Context, block *larkdocx.Block, subBlockTexts []string) (texts []string) {
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
			texts = append(texts, fmt.Sprintf("%s", strings.ReplaceAll(strings.Repeat("|:-:|", cols), "||", "|")))
		}
		var tmp []string
		for col := 0; col < cols; col++ {
			tmp = append(tmp, subBlockTexts[row*cols+col])
		}
		texts = append(texts, fmt.Sprintf("|%s|", strings.Join(tmp, "|")))
	}

	return FixTexts(texts)
}

func (p *DocxMarkdownProcessor) BlockQuoteContainerMarkdown(ctx context.Context, subBlockTexts []string) (texts []string) {
	for _, line := range subBlockTexts {
		texts = append(texts, fmt.Sprintf("> %s\n>", line))
	}
	return texts
}
