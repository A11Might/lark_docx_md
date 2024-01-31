package lark_docx_md

import (
	"context"
	"io"
	"os"
	"testing"

	"github.com/bytedance/mockey"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkdocx "github.com/larksuite/oapi-sdk-go/v3/service/docx/v1"
	larkdrive "github.com/larksuite/oapi-sdk-go/v3/service/drive/v1"
	"github.com/samber/lo"
)

func TestDocxMarkdownProcessor_BlockPageMarkdown(t *testing.T) {
	type fields struct {
		Config     *Config
		LarkClient *lark.Client
		DocumentId string
	}
	type args struct {
		ctx   context.Context
		block *larkdocx.Block
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			"page",
			fields{},
			args{
				context.Background(),
				larkdocx.NewBlockBuilder().
					BlockType(Page).
					Page(
						larkdocx.NewTextBuilder().Elements(
							[]*larkdocx.TextElement{
								larkdocx.NewTextElementBuilder().TextRun(
									larkdocx.NewTextRunBuilder().Content(
										"文章标题",
									).TextElementStyle(
										larkdocx.NewTextElementStyleBuilder().
											Bold(false).
											InlineCode(false).
											Italic(false).
											Strikethrough(false).
											Underline(false).
											Build(),
									).Build(),
								).Build(),
							},
						).Style(
							larkdocx.NewTextStyleBuilder().Align(AlignLeft).Build(),
						).Build(),
					).Build(),
			},
			"# 文章标题",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &DocxMarkdownProcessor{
				Config:     tt.fields.Config,
				LarkClient: tt.fields.LarkClient,
				DocumentId: tt.fields.DocumentId,
			}
			if got := p.BlockPageMarkdown(tt.args.ctx, tt.args.block); got != tt.want {
				t.Errorf("BlockPageMarkdown() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDocxMarkdownProcessor_BlockTextMarkdown(t *testing.T) {
	type fields struct {
		Config     *Config
		LarkClient *lark.Client
		DocumentId string
	}
	type args struct {
		ctx   context.Context
		block *larkdocx.Block
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			"text",
			fields{},
			args{
				context.Background(),
				larkdocx.NewBlockBuilder().
					BlockType(Text).
					Text(
						larkdocx.NewTextBuilder().Elements(
							[]*larkdocx.TextElement{
								larkdocx.NewTextElementBuilder().TextRun(
									larkdocx.NewTextRunBuilder().Content(
										"文本",
									).TextElementStyle(
										larkdocx.NewTextElementStyleBuilder().
											Bold(false).
											InlineCode(false).
											Italic(false).
											Strikethrough(false).
											Underline(false).
											Build(),
									).Build(),
								).Build(),
							},
						).Style(
							larkdocx.NewTextStyleBuilder().Align(AlignLeft).Build(),
						).Build(),
					).Build(),
			},
			"文本",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &DocxMarkdownProcessor{
				Config:     tt.fields.Config,
				LarkClient: tt.fields.LarkClient,
				DocumentId: tt.fields.DocumentId,
			}
			if got := p.BlockTextMarkdown(tt.args.ctx, tt.args.block); got != tt.want {
				t.Errorf("BlockTextMarkdown() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDocxMarkdownProcessor_BlockHeadingMarkdown(t *testing.T) {
	type fields struct {
		Config     *Config
		LarkClient *lark.Client
		DocumentId string
	}
	type args struct {
		ctx   context.Context
		block *larkdocx.Block
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			"heading1",
			fields{},
			args{
				context.Background(),
				larkdocx.NewBlockBuilder().
					BlockType(Heading1).
					Heading1(
						larkdocx.NewTextBuilder().Elements(
							[]*larkdocx.TextElement{
								larkdocx.NewTextElementBuilder().TextRun(
									larkdocx.NewTextRunBuilder().Content(
										"一级标题",
									).TextElementStyle(
										larkdocx.NewTextElementStyleBuilder().
											Bold(false).
											InlineCode(false).
											Italic(false).
											Strikethrough(false).
											Underline(false).
											Build(),
									).Build(),
								).Build(),
							},
						).Style(
							larkdocx.NewTextStyleBuilder().Align(AlignLeft).Build(),
						).Build(),
					).Build(),
			},
			"# 一级标题",
		},
		{
			"heading2",
			fields{},
			args{
				context.Background(),
				larkdocx.NewBlockBuilder().
					BlockType(Heading2).
					Heading2(
						larkdocx.NewTextBuilder().Elements(
							[]*larkdocx.TextElement{
								larkdocx.NewTextElementBuilder().TextRun(
									larkdocx.NewTextRunBuilder().Content(
										"二级标题",
									).TextElementStyle(
										larkdocx.NewTextElementStyleBuilder().
											Bold(false).
											InlineCode(false).
											Italic(false).
											Strikethrough(false).
											Underline(false).
											Build(),
									).Build(),
								).Build(),
							},
						).Style(
							larkdocx.NewTextStyleBuilder().Align(AlignLeft).Build(),
						).Build(),
					).Build(),
			},
			"## 二级标题",
		},
		{
			"heading3",
			fields{},
			args{
				context.Background(),
				larkdocx.NewBlockBuilder().
					BlockType(Heading3).
					Heading3(
						larkdocx.NewTextBuilder().Elements(
							[]*larkdocx.TextElement{
								larkdocx.NewTextElementBuilder().TextRun(
									larkdocx.NewTextRunBuilder().Content(
										"三级标题",
									).TextElementStyle(
										larkdocx.NewTextElementStyleBuilder().
											Bold(false).
											InlineCode(false).
											Italic(false).
											Strikethrough(false).
											Underline(false).
											Build(),
									).Build(),
								).Build(),
							},
						).Style(
							larkdocx.NewTextStyleBuilder().Align(AlignLeft).Build(),
						).Build(),
					).Build(),
			},
			"### 三级标题",
		},
		{
			"heading4",
			fields{},
			args{
				context.Background(),
				larkdocx.NewBlockBuilder().
					BlockType(Heading4).
					Heading4(
						larkdocx.NewTextBuilder().Elements(
							[]*larkdocx.TextElement{
								larkdocx.NewTextElementBuilder().TextRun(
									larkdocx.NewTextRunBuilder().Content(
										"四级标题",
									).TextElementStyle(
										larkdocx.NewTextElementStyleBuilder().
											Bold(false).
											InlineCode(false).
											Italic(false).
											Strikethrough(false).
											Underline(false).
											Build(),
									).Build(),
								).Build(),
							},
						).Style(
							larkdocx.NewTextStyleBuilder().Align(AlignLeft).Build(),
						).Build(),
					).Build(),
			},
			"#### 四级标题",
		},
		{
			"heading5",
			fields{},
			args{
				context.Background(),
				larkdocx.NewBlockBuilder().
					BlockType(Heading5).
					Heading5(
						larkdocx.NewTextBuilder().Elements(
							[]*larkdocx.TextElement{
								larkdocx.NewTextElementBuilder().TextRun(
									larkdocx.NewTextRunBuilder().Content(
										"五级标题",
									).TextElementStyle(
										larkdocx.NewTextElementStyleBuilder().
											Bold(false).
											InlineCode(false).
											Italic(false).
											Strikethrough(false).
											Underline(false).
											Build(),
									).Build(),
								).Build(),
							},
						).Style(
							larkdocx.NewTextStyleBuilder().Align(AlignLeft).Build(),
						).Build(),
					).Build(),
			},
			"##### 五级标题",
		},
		{
			"heading6",
			fields{},
			args{
				context.Background(),
				larkdocx.NewBlockBuilder().
					BlockType(Heading6).
					Heading6(
						larkdocx.NewTextBuilder().Elements(
							[]*larkdocx.TextElement{
								larkdocx.NewTextElementBuilder().TextRun(
									larkdocx.NewTextRunBuilder().Content(
										"六级标题",
									).TextElementStyle(
										larkdocx.NewTextElementStyleBuilder().
											Bold(false).
											InlineCode(false).
											Italic(false).
											Strikethrough(false).
											Underline(false).
											Build(),
									).Build(),
								).Build(),
							},
						).Style(
							larkdocx.NewTextStyleBuilder().Align(AlignLeft).Build(),
						).Build(),
					).Build(),
			},
			"###### 六级标题",
		},
		{
			"heading7",
			fields{},
			args{
				context.Background(),
				larkdocx.NewBlockBuilder().
					BlockType(Heading7).
					Heading7(
						larkdocx.NewTextBuilder().Elements(
							[]*larkdocx.TextElement{
								larkdocx.NewTextElementBuilder().TextRun(
									larkdocx.NewTextRunBuilder().Content(
										"七级标题",
									).TextElementStyle(
										larkdocx.NewTextElementStyleBuilder().
											Bold(false).
											InlineCode(false).
											Italic(false).
											Strikethrough(false).
											Underline(false).
											Build(),
									).Build(),
								).Build(),
							},
						).Style(
							larkdocx.NewTextStyleBuilder().Align(AlignLeft).Build(),
						).Build(),
					).Build(),
			},
			"####### 七级标题",
		},
		{
			"heading8",
			fields{},
			args{
				context.Background(),
				larkdocx.NewBlockBuilder().
					BlockType(Heading8).
					Heading8(
						larkdocx.NewTextBuilder().Elements(
							[]*larkdocx.TextElement{
								larkdocx.NewTextElementBuilder().TextRun(
									larkdocx.NewTextRunBuilder().Content(
										"八级标题",
									).TextElementStyle(
										larkdocx.NewTextElementStyleBuilder().
											Bold(false).
											InlineCode(false).
											Italic(false).
											Strikethrough(false).
											Underline(false).
											Build(),
									).Build(),
								).Build(),
							},
						).Style(
							larkdocx.NewTextStyleBuilder().Align(AlignLeft).Build(),
						).Build(),
					).Build(),
			},
			"######## 八级标题",
		},
		{
			"heading9",
			fields{},
			args{
				context.Background(),
				larkdocx.NewBlockBuilder().
					BlockType(Heading9).
					Heading9(
						larkdocx.NewTextBuilder().Elements(
							[]*larkdocx.TextElement{
								larkdocx.NewTextElementBuilder().TextRun(
									larkdocx.NewTextRunBuilder().Content(
										"九级标题",
									).TextElementStyle(
										larkdocx.NewTextElementStyleBuilder().
											Bold(false).
											InlineCode(false).
											Italic(false).
											Strikethrough(false).
											Underline(false).
											Build(),
									).Build(),
								).Build(),
							},
						).Style(
							larkdocx.NewTextStyleBuilder().Align(AlignLeft).Build(),
						).Build(),
					).Build(),
			},
			"######### 九级标题",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &DocxMarkdownProcessor{
				Config:     tt.fields.Config,
				LarkClient: tt.fields.LarkClient,
				DocumentId: tt.fields.DocumentId,
			}
			if got := p.BlockHeadingMarkdown(tt.args.ctx, tt.args.block); got != tt.want {
				t.Errorf("BlockHeadingMarkdown() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDocxMarkdownProcessor_BlockBulletMarkdown(t *testing.T) {
	type fields struct {
		Config     *Config
		LarkClient *lark.Client
		DocumentId string
	}
	type args struct {
		ctx   context.Context
		block *larkdocx.Block
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			"bullet",
			fields{},
			args{
				context.Background(),
				larkdocx.NewBlockBuilder().
					BlockType(Bullet).
					Bullet(
						larkdocx.NewTextBuilder().Elements(
							[]*larkdocx.TextElement{
								larkdocx.NewTextElementBuilder().TextRun(
									larkdocx.NewTextRunBuilder().Content(
										"无序列表",
									).TextElementStyle(
										larkdocx.NewTextElementStyleBuilder().
											Bold(false).
											InlineCode(false).
											Italic(false).
											Strikethrough(false).
											Underline(false).
											Build(),
									).Build(),
								).Build(),
							},
						).Style(
							larkdocx.NewTextStyleBuilder().Align(AlignLeft).Build(),
						).Build(),
					).Build(),
			},
			"- 无序列表",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &DocxMarkdownProcessor{
				Config:     tt.fields.Config,
				LarkClient: tt.fields.LarkClient,
				DocumentId: tt.fields.DocumentId,
			}
			if got := p.BlockBulletMarkdown(tt.args.ctx, tt.args.block); got != tt.want {
				t.Errorf("BlockBulletMarkdown() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
					larkdocx.NewTextStyleBuilder().Align(AlignLeft).Build(),
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
					larkdocx.NewTextStyleBuilder().Align(AlignLeft).Build(),
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
					larkdocx.NewTextStyleBuilder().Align(AlignLeft).Build(),
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

func TestDocxMarkdownProcessor_BlockOrderedMarkdown(t *testing.T) {
	type fields struct {
		Config     *Config
		LarkClient *lark.Client
		DocumentId string
	}
	type args struct {
		ctx   context.Context
		block *larkdocx.Block
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			"ordered",
			fields{},
			args{
				context.Background(),
				larkdocx.NewBlockBuilder().
					BlockType(Ordered).
					Ordered(
						larkdocx.NewTextBuilder().Elements(
							[]*larkdocx.TextElement{
								larkdocx.NewTextElementBuilder().TextRun(
									larkdocx.NewTextRunBuilder().Content(
										"有序列表",
									).TextElementStyle(
										larkdocx.NewTextElementStyleBuilder().
											Bold(false).
											InlineCode(false).
											Italic(false).
											Strikethrough(false).
											Underline(false).
											Build(),
									).Build(),
								).Build(),
							},
						).Style(
							larkdocx.NewTextStyleBuilder().Align(AlignLeft).Build(),
						).Build(),
					).Build(),
			},
			"1. 有序列表",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &DocxMarkdownProcessor{
				Config:     tt.fields.Config,
				LarkClient: tt.fields.LarkClient,
				DocumentId: tt.fields.DocumentId,
			}
			if got := p.BlockOrderedMarkdown(tt.args.ctx, tt.args.block); got != tt.want {
				t.Errorf("BlockOrderedMarkdown() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDocxMarkdownProcessor_BlockCalloutMarkdown(t *testing.T) {
	type fields struct {
		Config     *Config
		LarkClient *lark.Client
		DocumentId string
	}
	type args struct {
		ctx    context.Context
		blocks []*larkdocx.Block
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			"callout haven't background color",
			fields{
				Config: &Config{},
			},
			args{
				context.Background(),
				[]*larkdocx.Block{
					larkdocx.NewBlockBuilder().
						BlockType(Callout).
						Callout(
							larkdocx.NewCalloutBuilder().
								EmojiId("dog").
								Build(),
						).Build(),
					larkdocx.NewBlockBuilder().
						BlockType(Text).
						Text(
							larkdocx.NewTextBuilder().Elements(
								[]*larkdocx.TextElement{
									larkdocx.NewTextElementBuilder().TextRun(
										larkdocx.NewTextRunBuilder().Content(
											"高亮块",
										).TextElementStyle(
											larkdocx.NewTextElementStyleBuilder().
												Bold(false).
												InlineCode(false).
												Italic(false).
												Strikethrough(false).
												Underline(false).
												Build(),
										).Build(),
									).Build(),
								},
							).Style(
								larkdocx.NewTextStyleBuilder().Language(Go).Build(),
							).Build(),
						).Build(),
				},
			},
			"🐶 高亮块\n\n",
		},
		{
			"callout use github callout style",
			fields{
				Config: &Config{
					UseGhCallout: true,
				},
			},
			args{
				context.Background(),
				[]*larkdocx.Block{
					larkdocx.NewBlockBuilder().
						BlockType(Callout).
						Callout(
							larkdocx.NewCalloutBuilder().
								BackgroundColor(Blue).
								EmojiId("dog").
								Build(),
						).Build(),
					larkdocx.NewBlockBuilder().
						BlockType(Text).
						Text(
							larkdocx.NewTextBuilder().Elements(
								[]*larkdocx.TextElement{
									larkdocx.NewTextElementBuilder().TextRun(
										larkdocx.NewTextRunBuilder().Content(
											"高亮块",
										).TextElementStyle(
											larkdocx.NewTextElementStyleBuilder().
												Bold(false).
												InlineCode(false).
												Italic(false).
												Strikethrough(false).
												Underline(false).
												Build(),
										).Build(),
									).Build(),
								},
							).Style(
								larkdocx.NewTextStyleBuilder().Language(Go).Build(),
							).Build(),
						).Build(),
				},
			},
			"> [!NOTE]\n>\n> 🐶 高亮块\n>\n",
		},
		{
			"callout don't use github callout style",
			fields{
				Config: &Config{
					UseGhCallout: false,
				},
			},
			args{
				context.Background(),
				[]*larkdocx.Block{
					larkdocx.NewBlockBuilder().
						BlockType(Callout).
						Callout(
							larkdocx.NewCalloutBuilder().
								BackgroundColor(Blue).
								EmojiId("dog").
								Build(),
						).Build(),
					larkdocx.NewBlockBuilder().
						BlockType(Text).
						Text(
							larkdocx.NewTextBuilder().Elements(
								[]*larkdocx.TextElement{
									larkdocx.NewTextElementBuilder().TextRun(
										larkdocx.NewTextRunBuilder().Content(
											"高亮块",
										).TextElementStyle(
											larkdocx.NewTextElementStyleBuilder().
												Bold(false).
												InlineCode(false).
												Italic(false).
												Strikethrough(false).
												Underline(false).
												Build(),
										).Build(),
									).Build(),
								},
							).Style(
								larkdocx.NewTextStyleBuilder().Language(Go).Build(),
							).Build(),
						).Build(),
				},
			},
			"> 🐶 高亮块\n>\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &DocxMarkdownProcessor{
				Config:     tt.fields.Config,
				LarkClient: tt.fields.LarkClient,
				DocumentId: tt.fields.DocumentId,
			}
			if got := p.BlockCalloutMarkdown(tt.args.ctx, tt.args.blocks); got != tt.want {
				t.Errorf("BlockCalloutMarkdown() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDocxMarkdownProcessor_BlockCodeMarkdown(t *testing.T) {
	type fields struct {
		Config     *Config
		LarkClient *lark.Client
		DocumentId string
	}
	type args struct {
		ctx   context.Context
		block *larkdocx.Block
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			"code",
			fields{},
			args{
				context.Background(),
				larkdocx.NewBlockBuilder().
					BlockType(Code).
					Code(
						larkdocx.NewTextBuilder().Elements(
							[]*larkdocx.TextElement{
								larkdocx.NewTextElementBuilder().TextRun(
									larkdocx.NewTextRunBuilder().Content(
										"fmt.Println(\"hello world\")",
									).TextElementStyle(
										larkdocx.NewTextElementStyleBuilder().
											Bold(false).
											InlineCode(false).
											Italic(false).
											Strikethrough(false).
											Underline(false).
											Build(),
									).Build(),
								).Build(),
							},
						).Style(
							larkdocx.NewTextStyleBuilder().Language(Go).Build(),
						).Build(),
					).Build(),
			},
			"```go\nfmt.Println(\"hello world\")\n```",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &DocxMarkdownProcessor{
				Config:     tt.fields.Config,
				LarkClient: tt.fields.LarkClient,
				DocumentId: tt.fields.DocumentId,
			}
			if got := p.BlockCodeMarkdown(tt.args.ctx, tt.args.block); got != tt.want {
				t.Errorf("BlockCodeMarkdown() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDocxMarkdownProcessor_BlockDividerMarkdown(t *testing.T) {
	type fields struct {
		Config     *Config
		LarkClient *lark.Client
		DocumentId string
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			"divider",
			fields{},
			args{
				context.Background(),
			},
			"---",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &DocxMarkdownProcessor{
				Config:     tt.fields.Config,
				LarkClient: tt.fields.LarkClient,
				DocumentId: tt.fields.DocumentId,
			}
			if got := p.BlockDividerMarkdown(tt.args.ctx); got != tt.want {
				t.Errorf("BlockDividerMarkdown() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDocxMarkdownProcessor_BlockImageMarkdown(t *testing.T) {
	client := lark.NewClient("appId", "secret")
	type fields struct {
		Config     *Config
		LarkClient *lark.Client
		DocumentId string
	}
	type args struct {
		ctx   context.Context
		block *larkdocx.Block
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
		mock   func()
	}{
		{
			"image static as url",
			fields{
				Config: &Config{
					StaticAsURL: true,
				},
				LarkClient: client,
			},
			args{
				context.Background(),
				larkdocx.NewBlockBuilder().
					BlockType(Image).
					Image(
						larkdocx.NewImageBuilder().
							Token("image-token").
							Width(100).
							Height(100).
							Build(),
					).Build(),
			},
			"<img src=\"image-token-url\" width=\"100\" height=\"100\"/>",
			func() {
				mockey.Mock(mockey.GetMethod(client.Drive.V1.Media, "BatchGetTmpDownloadUrl")).Return(
					&larkdrive.BatchGetTmpDownloadUrlMediaResp{
						Data: &larkdrive.BatchGetTmpDownloadUrlMediaRespData{
							TmpDownloadUrls: []*larkdrive.TmpDownloadUrl{
								{
									TmpDownloadUrl: lo.ToPtr("image-token-url"),
								},
							},
						},
					},
					nil,
				).Build()
			},
		},
		{
			"image download static",
			fields{
				Config: &Config{
					StaticDir:   "static",
					FilePrefix:  "static",
					StaticAsURL: false,
				},
				LarkClient: client,
			},
			args{
				context.Background(),
				larkdocx.NewBlockBuilder().
					BlockType(Image).
					Image(
						larkdocx.NewImageBuilder().
							Token("image-token").
							Width(100).
							Height(100).
							Build(),
					).Build(),
			},
			"<img src=\"static/image-token.jpg\" width=\"100\" height=\"100\"/>",
			func() {
				mockey.Mock(mockey.GetMethod(client.Drive.V1.Media, "Download")).Return(
					&larkdrive.DownloadMediaResp{},
					nil,
				).Build()
				mockey.Mock(os.MkdirAll).Return(nil).Build()
				mockey.Mock(os.OpenFile).Return(
					&os.File{},
					nil,
				).Build()
				mockey.Mock(io.Copy).Return(
					int64(1),
					nil,
				).Build()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &DocxMarkdownProcessor{
				Config:     tt.fields.Config,
				LarkClient: tt.fields.LarkClient,
				DocumentId: tt.fields.DocumentId,
			}
			mockey.PatchConvey(tt.name, t, func() {
				tt.mock()
				if got := p.BlockImageMarkdown(tt.args.ctx, tt.args.block); got != tt.want {
					t.Errorf("BlockImageMarkdown() = %v, want %v", got, tt.want)
				}
			})
		})
	}
}

func TestDocxMarkdownProcessor_BlockQuoteContainerMarkdown(t *testing.T) {
	type fields struct {
		Config     *Config
		LarkClient *lark.Client
		DocumentId string
	}
	type args struct {
		ctx    context.Context
		blocks []*larkdocx.Block
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			"quote container",
			fields{},
			args{
				context.Background(),
				[]*larkdocx.Block{
					larkdocx.NewBlockBuilder().
						BlockType(Text).
						Text(
							larkdocx.NewTextBuilder().Elements(
								[]*larkdocx.TextElement{
									larkdocx.NewTextElementBuilder().TextRun(
										larkdocx.NewTextRunBuilder().Content(
											"引用内容1",
										).TextElementStyle(
											larkdocx.NewTextElementStyleBuilder().
												Bold(false).
												InlineCode(false).
												Italic(false).
												Strikethrough(false).
												Underline(false).
												Build(),
										).Build(),
									).Build(),
								},
							).Style(
								larkdocx.NewTextStyleBuilder().Align(AlignLeft).Build(),
							).Build(),
						).Build(),
					larkdocx.NewBlockBuilder().
						BlockType(Text).
						Text(
							larkdocx.NewTextBuilder().Elements(
								[]*larkdocx.TextElement{
									larkdocx.NewTextElementBuilder().TextRun(
										larkdocx.NewTextRunBuilder().Content(
											"引用内容2",
										).TextElementStyle(
											larkdocx.NewTextElementStyleBuilder().
												Bold(false).
												InlineCode(false).
												Italic(false).
												Strikethrough(false).
												Underline(false).
												Build(),
										).Build(),
									).Build(),
								},
							).Style(
								larkdocx.NewTextStyleBuilder().Align(AlignLeft).Build(),
							).Build(),
						).Build(),
				},
			},
			"> 引用内容1\n>\n> 引用内容2\n>\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &DocxMarkdownProcessor{
				Config:     tt.fields.Config,
				LarkClient: tt.fields.LarkClient,
				DocumentId: tt.fields.DocumentId,
			}
			if got := p.BlockQuoteContainerMarkdown(tt.args.ctx, tt.args.blocks); got != tt.want {
				t.Errorf("BlockQuoteContainerMarkdown() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDocxMarkdownProcessor_BlockQuoteMarkdown(t *testing.T) {
	type fields struct {
		Config     *Config
		LarkClient *lark.Client
		DocumentId string
	}
	type args struct {
		ctx   context.Context
		block *larkdocx.Block
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			"quote",
			fields{},
			args{
				context.Background(),
				larkdocx.NewBlockBuilder().
					BlockType(Quote).
					Quote(
						larkdocx.NewTextBuilder().Elements(
							[]*larkdocx.TextElement{
								larkdocx.NewTextElementBuilder().TextRun(
									larkdocx.NewTextRunBuilder().Content(
										"引用内容",
									).TextElementStyle(
										larkdocx.NewTextElementStyleBuilder().
											Bold(false).
											InlineCode(false).
											Italic(false).
											Strikethrough(false).
											Underline(false).
											Build(),
									).Build(),
								).Build(),
							},
						).Style(
							larkdocx.NewTextStyleBuilder().Align(AlignLeft).Build(),
						).Build(),
					).Build(),
			},
			"> 引用内容",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &DocxMarkdownProcessor{
				Config:     tt.fields.Config,
				LarkClient: tt.fields.LarkClient,
				DocumentId: tt.fields.DocumentId,
			}
			if got := p.BlockQuoteMarkdown(tt.args.ctx, tt.args.block); got != tt.want {
				t.Errorf("BlockQuoteMarkdown() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDocxMarkdownProcessor_BlockTableMarkdown(t *testing.T) {
	type fields struct {
		Config     *Config
		LarkClient *lark.Client
		DocumentId string
	}
	type args struct {
		ctx   context.Context
		table [][]*larkdocx.Block
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			"table",
			fields{},
			args{
				context.Background(),
				[][]*larkdocx.Block{
					{
						larkdocx.NewBlockBuilder().
							BlockType(Text).
							Text(
								larkdocx.NewTextBuilder().Elements(
									[]*larkdocx.TextElement{
										larkdocx.NewTextElementBuilder().TextRun(
											larkdocx.NewTextRunBuilder().Content(
												"表头第一个单元格",
											).TextElementStyle(
												larkdocx.NewTextElementStyleBuilder().
													Bold(false).
													InlineCode(false).
													Italic(false).
													Strikethrough(false).
													Underline(false).
													Build(),
											).Build(),
										).Build(),
									},
								).Style(
									larkdocx.NewTextStyleBuilder().Align(AlignLeft).Build(),
								).Build(),
							).Build(),
						larkdocx.NewBlockBuilder().
							BlockType(Text).
							Text(
								larkdocx.NewTextBuilder().Elements(
									[]*larkdocx.TextElement{
										larkdocx.NewTextElementBuilder().TextRun(
											larkdocx.NewTextRunBuilder().Content(
												"表头第二个单元格",
											).TextElementStyle(
												larkdocx.NewTextElementStyleBuilder().
													Bold(false).
													InlineCode(false).
													Italic(false).
													Strikethrough(false).
													Underline(false).
													Build(),
											).Build(),
										).Build(),
									},
								).Style(
									larkdocx.NewTextStyleBuilder().Align(AlignMid).Build(),
								).Build(),
							).Build(),
					},
					{
						larkdocx.NewBlockBuilder().
							BlockType(Text).
							Text(
								larkdocx.NewTextBuilder().Elements(
									[]*larkdocx.TextElement{
										larkdocx.NewTextElementBuilder().TextRun(
											larkdocx.NewTextRunBuilder().Content(
												"第三个单元格",
											).TextElementStyle(
												larkdocx.NewTextElementStyleBuilder().
													Bold(false).
													InlineCode(false).
													Italic(false).
													Strikethrough(false).
													Underline(false).
													Build(),
											).Build(),
										).Build(),
									},
								).Style(
									larkdocx.NewTextStyleBuilder().Align(AlignLeft).Build(),
								).Build(),
							).Build(),
						larkdocx.NewBlockBuilder().
							BlockType(Text).
							Text(
								larkdocx.NewTextBuilder().Elements(
									[]*larkdocx.TextElement{
										larkdocx.NewTextElementBuilder().TextRun(
											larkdocx.NewTextRunBuilder().Content(
												"第四个单元格",
											).TextElementStyle(
												larkdocx.NewTextElementStyleBuilder().
													Bold(false).
													InlineCode(false).
													Italic(false).
													Strikethrough(false).
													Underline(false).
													Build(),
											).Build(),
										).Build(),
									},
								).Style(
									larkdocx.NewTextStyleBuilder().Align(AlignLeft).Build(),
								).Build(),
							).Build(),
					},
				},
			},
			"|表头第一个单元格|表头第二个单元格|\n|-|:-:|\n|第三个单元格|第四个单元格|\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &DocxMarkdownProcessor{
				Config:     tt.fields.Config,
				LarkClient: tt.fields.LarkClient,
				DocumentId: tt.fields.DocumentId,
			}
			if got := p.BlockTableMarkdown(tt.args.ctx, tt.args.table); got != tt.want {
				t.Errorf("BlockTableMarkdown() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDocxMarkdownProcessor_BlockTodoMarkdown(t *testing.T) {
	type fields struct {
		Config     *Config
		LarkClient *lark.Client
		DocumentId string
	}
	type args struct {
		ctx   context.Context
		block *larkdocx.Block
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			"todo undone",
			fields{},
			args{
				context.Background(),
				larkdocx.NewBlockBuilder().
					BlockType(Todo).
					Todo(
						larkdocx.NewTextBuilder().Elements(
							[]*larkdocx.TextElement{
								larkdocx.NewTextElementBuilder().TextRun(
									larkdocx.NewTextRunBuilder().Content(
										"未完成待办事项",
									).TextElementStyle(
										larkdocx.NewTextElementStyleBuilder().
											Bold(false).
											InlineCode(false).
											Italic(false).
											Strikethrough(false).
											Underline(false).
											Build(),
									).Build(),
								).Build(),
							},
						).Style(
							larkdocx.NewTextStyleBuilder().Align(AlignLeft).Done(false).Build(),
						).Build(),
					).Build(),
			},
			"[] 未完成待办事项",
		},
		{
			"todo done",
			fields{},
			args{
				context.Background(),
				larkdocx.NewBlockBuilder().
					BlockType(Todo).
					Todo(
						larkdocx.NewTextBuilder().Elements(
							[]*larkdocx.TextElement{
								larkdocx.NewTextElementBuilder().TextRun(
									larkdocx.NewTextRunBuilder().Content(
										"已完成待办事项",
									).TextElementStyle(
										larkdocx.NewTextElementStyleBuilder().
											Bold(false).
											InlineCode(false).
											Italic(false).
											Strikethrough(false).
											Underline(false).
											Build(),
									).Build(),
								).Build(),
							},
						).Style(
							larkdocx.NewTextStyleBuilder().Align(AlignLeft).Done(true).Build(),
						).Build(),
					).Build(),
			},
			"[x] 已完成待办事项",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &DocxMarkdownProcessor{
				Config:     tt.fields.Config,
				LarkClient: tt.fields.LarkClient,
				DocumentId: tt.fields.DocumentId,
			}
			if got := p.BlockTodoMarkdown(tt.args.ctx, tt.args.block); got != tt.want {
				t.Errorf("BlockTodoMarkdown() = %v, want %v", got, tt.want)
			}
		})
	}
}
