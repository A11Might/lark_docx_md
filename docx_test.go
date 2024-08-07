package lark_docx_md

import (
	"context"
	"reflect"
	"testing"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkdocx "github.com/larksuite/oapi-sdk-go/v3/service/docx/v1"
	"github.com/samber/lo"
)

func TestListTransformToTree(t *testing.T) {
	larkBlocks := []*larkdocx.Block{
		larkdocx.NewBlockBuilder().
			BlockId("标题").
			BlockType(Page).
			Children([]string{"文本", "引用容器", "表格"}).
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
		larkdocx.NewBlockBuilder().
			BlockId("文本").
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
		larkdocx.NewBlockBuilder().
			BlockId("引用容器").
			BlockType(QuoteContainer).
			Children([]string{"引用内容1", "引用内容2"}).
			QuoteContainer(&larkdocx.QuoteContainer{}).
			Build(),
		larkdocx.NewBlockBuilder().
			BlockId("引用容器1").
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
			BlockId("引用容器2").
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
		larkdocx.NewBlockBuilder().
			BlockId("表格").
			BlockType(Table).
			Children([]string{"单元格一", "单元格二", "单元格三", "单元格四"}).
			Build(),
		larkdocx.NewBlockBuilder().
			BlockId("单元格一").
			BlockType(TableCell).
			Children([]string{"表头第一个单元格"}).
			Build(),
		larkdocx.NewBlockBuilder().
			BlockId("表头第一个单元格").
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
			BlockId("单元格二").
			BlockType(TableCell).
			Children([]string{"表头第二个单元格"}).
			Build(),
		larkdocx.NewBlockBuilder().
			BlockId("表头第二个单元格").
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
		larkdocx.NewBlockBuilder().
			BlockId("单元格三").
			BlockType(TableCell).
			Children([]string{"第三个单元格"}).
			Build(),
		larkdocx.NewBlockBuilder().
			BlockId("第三个单元格").
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
			BlockId("单元格四").
			BlockType(TableCell).
			Children([]string{"第四个单元格"}).
			Build(),
		larkdocx.NewBlockBuilder().
			BlockId("第四个单元格").
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
	}
	larkBlockMap := lo.SliceToMap(larkBlocks, func(item *larkdocx.Block) (string, *larkdocx.Block) {
		return *item.BlockId, item
	})

	type fields struct {
		Config     *Config
		LarkClient *lark.Client
		DocumentId string
	}
	type args struct {
		ctx           context.Context
		rootLarkBlock *larkdocx.Block
		larkBlockMap  map[string]*larkdocx.Block
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Node
	}{
		{
			"docx",
			fields{},
			args{
				context.Background(),
				larkBlocks[0],
				larkBlockMap,
			},
			&Node{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &DocxMarkdownProcessor{
				Config:     tt.fields.Config,
				LarkClient: tt.fields.LarkClient,
				DocumentId: tt.fields.DocumentId,
			}
			if got := p.listTransformToTree(tt.args.ctx, tt.args.rootLarkBlock, tt.args.larkBlockMap); !reflect.DeepEqual(got, tt.want) {
				larkcore.Prettify(got)
				// t.Errorf("listTransformToTree() = %v, want %v", larkcore.Prettify(got), larkcore.Prettify(tt.want))
			}
		})
	}
}
