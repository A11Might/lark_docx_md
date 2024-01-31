package lark_docx_md

import (
	"context"
	"testing"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkdocx "github.com/larksuite/oapi-sdk-go/v3/service/docx/v1"
)

func TestDocxMarkdownProcessor_TextMarkdown(t *testing.T) {
	type fields struct {
		Config     *Config
		LarkClient *lark.Client
		DocumentId string
	}
	type args struct {
		ctx  context.Context
		text *larkdocx.Text
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			"case 1",
			fields{},
			args{
				context.Background(),
				larkdocx.NewTextBuilder().Elements(
					[]*larkdocx.TextElement{
						larkdocx.NewTextElementBuilder().TextRun(
							larkdocx.NewTextRunBuilder().Content(
								"删除线",
							).TextElementStyle(
								larkdocx.NewTextElementStyleBuilder().
									Bold(false).
									InlineCode(false).
									Italic(false).
									Strikethrough(true).
									Underline(false).
									Build(),
							).Build(),
						).Build(),
						larkdocx.NewTextElementBuilder().TextRun(
							larkdocx.NewTextRunBuilder().Content(
								"链接",
							).TextElementStyle(
								larkdocx.NewTextElementStyleBuilder().
									Bold(false).
									InlineCode(false).
									Italic(false).
									Link(
										larkdocx.NewLinkBuilder().Url("https://github.com/").Build(),
									).
									Strikethrough(true).
									Underline(false).
									Build(),
							).Build(),
						).Build(),
						larkdocx.NewTextElementBuilder().TextRun(
							larkdocx.NewTextRunBuilder().Content(
								"删除线",
							).TextElementStyle(
								larkdocx.NewTextElementStyleBuilder().
									Bold(false).
									InlineCode(false).
									Italic(false).
									Strikethrough(true).
									Underline(false).
									Build(),
							).Build(),
						).Build(),
					},
				).Style(
					larkdocx.NewTextStyleBuilder().Align(1).Build(),
				).Build(),
			},
			"~~删除线[链接](https://github.com/)删除线~~",
		},
		{
			"case 2",
			fields{},
			args{
				context.Background(),
				larkdocx.NewTextBuilder().Elements(
					[]*larkdocx.TextElement{
						larkdocx.NewTextElementBuilder().TextRun(
							larkdocx.NewTextRunBuilder().Content(
								"删除线",
							).TextElementStyle(
								larkdocx.NewTextElementStyleBuilder().
									Bold(false).
									InlineCode(false).
									Italic(false).
									Strikethrough(true).
									Underline(false).
									Build(),
							).Build(),
						).Build(),
						larkdocx.NewTextElementBuilder().TextRun(
							larkdocx.NewTextRunBuilder().Content(
								"链接",
							).TextElementStyle(
								larkdocx.NewTextElementStyleBuilder().
									Bold(false).
									InlineCode(false).
									Italic(false).
									Link(
										larkdocx.NewLinkBuilder().Url("https://github.com/").Build(),
									).
									Strikethrough(false).
									Underline(false).
									Build(),
							).Build(),
						).Build(),
						larkdocx.NewTextElementBuilder().TextRun(
							larkdocx.NewTextRunBuilder().Content(
								"删除线",
							).TextElementStyle(
								larkdocx.NewTextElementStyleBuilder().
									Bold(false).
									InlineCode(false).
									Italic(false).
									Strikethrough(true).
									Underline(false).
									Build(),
							).Build(),
						).Build(),
					},
				).Style(
					larkdocx.NewTextStyleBuilder().Align(1).Build(),
				).Build(),
			},
			"~~删除线~~[链接](https://github.com/)~~删除线~~",
		},
		{
			"case 3",
			fields{},
			args{
				context.Background(),
				larkdocx.NewTextBuilder().Elements(
					[]*larkdocx.TextElement{
						larkdocx.NewTextElementBuilder().TextRun(
							larkdocx.NewTextRunBuilder().Content(
								"删除线",
							).TextElementStyle(
								larkdocx.NewTextElementStyleBuilder().
									Bold(false).
									InlineCode(false).
									Italic(false).
									Strikethrough(true).
									Underline(false).
									Build(),
							).Build(),
						).Build(),
						larkdocx.NewTextElementBuilder().TextRun(
							larkdocx.NewTextRunBuilder().Content(
								"链接",
							).TextElementStyle(
								larkdocx.NewTextElementStyleBuilder().
									Bold(false).
									InlineCode(false).
									Italic(false).
									Link(
										larkdocx.NewLinkBuilder().Url("https://github.com/").Build(),
									).
									Strikethrough(true).
									Underline(false).
									Build(),
							).Build(),
						).Build(),
						larkdocx.NewTextElementBuilder().TextRun(
							larkdocx.NewTextRunBuilder().Content(
								"加粗",
							).TextElementStyle(
								larkdocx.NewTextElementStyleBuilder().
									Bold(true).
									InlineCode(false).
									Italic(false).
									Strikethrough(false).
									Underline(false).
									Build(),
							).Build(),
						).Build(),
					},
				).Style(
					larkdocx.NewTextStyleBuilder().Align(1).Build(),
				).Build(),
			},
			"~~删除线[链接](https://github.com/)~~**加粗**",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &DocxMarkdownProcessor{
				Config:     tt.fields.Config,
				LarkClient: tt.fields.LarkClient,
				DocumentId: tt.fields.DocumentId,
			}
			if got := p.TextMarkdown(tt.args.ctx, tt.args.text); got != tt.want {
				t.Errorf("TextMarkdown() = %v, want %v", got, tt.want)
			}
		})
	}
}
