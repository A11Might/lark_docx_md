package lark_docx_md

import (
	"context"
	"strings"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkdocx "github.com/larksuite/oapi-sdk-go/v3/service/docx/v1"
	"github.com/samber/lo"
)

type ProcessItem struct {
	BlockType      *int                // block 类型
	Normal         *larkdocx.Block     // 普通块
	QuoteContainer []*larkdocx.Block   // 引用容器
	Table          [][]*larkdocx.Block // 表格
}

func DocxMarkdown(ctx context.Context, client *lark.Client, documentId string) (string, error) {
	req := larkdocx.NewListDocumentBlockReqBuilder().DocumentId(documentId).Build()
	iterator, _ := client.Docx.V1.DocumentBlock.ListByIterator(ctx, req)

	// 读出所有块
	var (
		hasMore = true
		block   *larkdocx.Block
		err     error

		allBlock []*larkdocx.Block
	)
	for ; hasMore && err == nil; hasMore, block, err = iterator.Next() {
		allBlock = append(allBlock, block)
	}
	if err != nil {
		return "", err
	}
	if len(allBlock) > 1 {
		allBlock = allBlock[1:] // 去除第一个空块
	}

	// 分组：普通块、引用容器和表格
	var (
		processItems []*ProcessItem
		allBlockMap  map[string]*larkdocx.Block
		st           = make(map[string]struct{}) // see twice
	)
	allBlockMap = lo.SliceToMap(allBlock, func(item *larkdocx.Block) (string, *larkdocx.Block) {
		return *item.BlockId, item
	})
	for _, b := range allBlock {
		if _, ok := st[*b.BlockId]; ok {
			continue
		}
		st[*b.BlockId] = struct{}{}
		processItem := &ProcessItem{
			BlockType: b.BlockType,
		}
		switch *b.BlockType {
		case QuoteContainer:
			for _, cBlockId := range b.Children {
				st[cBlockId] = struct{}{}
				processItem.QuoteContainer = append(processItem.QuoteContainer, allBlockMap[cBlockId])
			}

		case Table:
			var children []*larkdocx.Block
			for _, cBlockId := range b.Children {
				// todo 处理一个单元格中包含多个块的情况
				st[cBlockId] = struct{}{}
				st[allBlockMap[cBlockId].Children[0]] = struct{}{}
				children = append(children, allBlockMap[allBlockMap[cBlockId].Children[0]])
			}
			processItem.Table = lo.Chunk[*larkdocx.Block](children, *b.Table.Property.ColumnSize)

		default:
			processItem.Normal = b
		}
		processItems = append(processItems, processItem)
	}

	// 转为 Markdown
	var buf = new(strings.Builder)
	for _, item := range processItems {
		buf.WriteString(DocxBlockMarkdown(ctx, client, item))
	}
	return buf.String(), nil
}
