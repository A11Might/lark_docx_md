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
	Page           = 1
	Text           = 2
	Heading1       = 3
	Heading2       = 4
	Heading3       = 5
	Heading4       = 6
	Heading5       = 7
	Heading6       = 8
	Heading7       = 9
	Heading8       = 10
	Heading9       = 11
	Bullet         = 12
	Ordered        = 13
	Code           = 14
	Quote          = 15
	Todo           = 17
	Divider        = 22
	Image          = 27
	Table          = 31
	TableCell      = 32
	QuoteContainer = 34
)

func DocxBlockMarkdown(ctx context.Context, client *lark.Client, item *ProcessItem) (md string) {
	if item == nil {
		return ""
	}

	switch *item.BlockType {
	case Page:
		md = BlockPageMarkdown(ctx, item.Normal)
	case Text:
		md = BlockTextMarkdown(ctx, item.Normal)
	case Heading1, Heading2, Heading3, Heading4, Heading5, Heading6, Heading7, Heading8, Heading9:
		md = BlockHeadingMarkdown(ctx, item.Normal)
	case Bullet:
		md = BlockBulletMarkdown(ctx, item.Normal)
	case Ordered:
		md = BlockOrderedMarkdown(ctx, item.Normal)
	case Code:
		md = BlockCodeMarkdown(ctx, item.Normal)
	case Quote:
		md = BlockQuoteMarkdown(ctx, item.Normal)
	case Todo:
		md = BlockTodoMarkdown(ctx, item.Normal)
	case Divider:
		md = BlockDividerMarkdown(ctx)
	case Image:
		md = BlockImageMarkdown(ctx, client, item.Normal)
	case Table:
		md = BlockTableMarkdown(ctx, item.Table)
	case QuoteContainer:
		md = BlockQuoteContainerMarkdown(ctx, item.QuoteContainer)
	default:
		md = fmt.Sprintf("not support block type:%d", *item.BlockType)
	}

	return md + "\n\n"
}

func BlockPageMarkdown(ctx context.Context, block *larkdocx.Block) string {
	return "# " + TextMarkdown(ctx, block.Page)
}

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
	return strings.Repeat("#", *block.BlockType-2) + " " + TextMarkdown(ctx, heading)
}

func BlockBulletMarkdown(ctx context.Context, block *larkdocx.Block) string {
	return "- " + TextMarkdown(ctx, block.Bullet)
}

func BlockOrderedMarkdown(ctx context.Context, block *larkdocx.Block) string {
	return "1. " + TextMarkdown(ctx, block.Ordered)
}

const (
	PlainText = iota + 1
	ABAP
	Ada
	Apache
	Apex
	AssemblyLanguage
	Bash
	CSharp
	Cpp
	C
	COBOL
	CSS
	CoffeeScript
	D
	Dart
	Delphi
	Django
	Dockerfile
	Erlang
	Fortran
	FoxPro
	Go
	Groovy
	HTML
	HTMLBars
	HTTP
	Haskell
	JSON
	Java
	JavaScript
	Julia
	Kotlin
	LateX
	Lisp
	Logo
	Lua
	MATLAB
	Makefile
	Markdown
	Nginx
	ObjectiveC
	OpenEdgeABL
	PHP
	Perl
	PostScript
	PowerShell
	Prolog
	ProtoBuf
	Python
	R
	RPG
	Ruby
	Rust
	SAS
	SCSS
	SQL
	Scala
	Scheme
	Scratch
	Shell
	Swift
	Thrift
	TypeScript
	VBScript
	VisualBasic
	XML
	YAML
	CMake
	Diff
	Gherkin
	GraphQL
	OpenGLShadingLanguage
	Properties
	Solidity
	TOML
)

var lmap = map[int]string{
	PlainText:             "plaintext",
	ABAP:                  "abap",
	Ada:                   "ada",
	Apache:                "apache",
	Apex:                  "apex",
	AssemblyLanguage:      "assemblylanguage",
	Bash:                  "bash",
	CSharp:                "csharp",
	Cpp:                   "cpp",
	C:                     "c",
	COBOL:                 "cobol",
	CSS:                   "css",
	CoffeeScript:          "coffeescript",
	D:                     "d",
	Dart:                  "dart",
	Delphi:                "delphi",
	Django:                "django",
	Dockerfile:            "dockerfile",
	Erlang:                "erlang",
	Fortran:               "fortran",
	FoxPro:                "foxpro",
	Go:                    "go",
	Groovy:                "groovy",
	HTML:                  "html",
	HTMLBars:              "htmlbars",
	HTTP:                  "http",
	Haskell:               "haskell",
	JSON:                  "json",
	Java:                  "java",
	JavaScript:            "javascript",
	Julia:                 "julia",
	Kotlin:                "kotlin",
	LateX:                 "latex",
	Lisp:                  "lisp",
	Logo:                  "logo",
	Lua:                   "lua",
	MATLAB:                "matlab",
	Makefile:              "makefile",
	Markdown:              "markdown",
	Nginx:                 "nginx",
	ObjectiveC:            "objectivec",
	OpenEdgeABL:           "openedgeabl",
	PHP:                   "php",
	Perl:                  "perl",
	PostScript:            "postscript",
	PowerShell:            "powershell",
	Prolog:                "prolog",
	ProtoBuf:              "protobuf",
	Python:                "python",
	R:                     "r",
	RPG:                   "rpg",
	Ruby:                  "ruby",
	Rust:                  "rust",
	SAS:                   "sas",
	SCSS:                  "scss",
	SQL:                   "sql",
	Scala:                 "scala",
	Scheme:                "scheme",
	Scratch:               "scratch",
	Shell:                 "shell",
	Swift:                 "swift",
	Thrift:                "thrift",
	TypeScript:            "typescript",
	VBScript:              "vbscript",
	VisualBasic:           "visualbasic",
	XML:                   "xml",
	YAML:                  "yaml",
	CMake:                 "cmake",
	Diff:                  "diff",
	Gherkin:               "gherkin",
	GraphQL:               "graphql",
	OpenGLShadingLanguage: "openglshadinglanguage",
	Properties:            "properties",
	Solidity:              "solidity",
	TOML:                  "toml",
}

func BlockCodeMarkdown(ctx context.Context, block *larkdocx.Block) string {
	return fmt.Sprintf("```%s\n%s\n```", lmap[*block.Code.Style.Language], TextMarkdown(ctx, block.Code))
}

func BlockQuoteMarkdown(ctx context.Context, block *larkdocx.Block) string {
	return "> " + TextMarkdown(ctx, block.Code)
}

func BlockTodoMarkdown(ctx context.Context, block *larkdocx.Block) string {
	if *block.Todo.Style.Done {
		return "[x] " + TextMarkdown(ctx, block.Todo)
	}
	return "[]" + TextMarkdown(ctx, block.Todo)
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

func BlockDividerMarkdown(ctx context.Context) string {
	return "---"
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

const (
	AlignLeft = iota + 1
	AlignMid
	AlignRight
)

func BlockTableMarkdown(ctx context.Context, table [][]*larkdocx.Block) string {
	buf := new(strings.Builder)

	/**
	| Tables        | Are           | Cool  |
	| ------------- |:-------------:| -----:|
	| col 3 is      | right-aligned | $1600 |
	| col 2 is      | centered      |   $12 |
	| zebra stripes | are neat      |    $1 |
	*/
	// 处理表头
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
		buf.WriteString(TextMarkdown(ctx, col.Text))
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
			buf.WriteString(TextMarkdown(ctx, col.Text))
			buf.WriteString("|")
		}
		buf.WriteString("\n")
	}

	return buf.String()
}

func BlockQuoteContainerMarkdown(ctx context.Context, blocks []*larkdocx.Block) string {
	buf := new(strings.Builder)

	for _, b := range blocks {
		buf.WriteString("> ")
		buf.WriteString(TextMarkdown(ctx, b.Text))
		buf.WriteString("\n>\n")
	}

	return buf.String()
}
