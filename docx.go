package lark_docx_md

import (
	"context"
	"strings"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkdocx "github.com/larksuite/oapi-sdk-go/v3/service/docx/v1"
	"github.com/samber/lo"
)

type Config struct {
	StaticDir    string // 如果需要下载静态文件，那么需要指定静态文件的目录
	FilePrefix   string // 针对静态文件，需要指定文件在 Markdown 中的前缀
	StaticAsURL  bool   // 不下载静态文件，直接把静态文件的 URL 插入到 Markdown 中
	UseGhCallout bool   // 高亮块使用 github 样式
}

type Option func(*DocxMarkdownProcessor)

// DownloadStatic 下载图片等静态文件
func DownloadStatic(staticDir, filePrefix string) Option {
	return func(p *DocxMarkdownProcessor) {
		p.StaticDir = staticDir
		p.FilePrefix = filePrefix
		p.StaticAsURL = false
	}
}

// UseGhCalloutStyle 使用 github 高亮块样式
func UseGhCalloutStyle() Option {
	return func(p *DocxMarkdownProcessor) {
		p.UseGhCallout = true
	}
}

type DocxMarkdownProcessor struct {
	*Config
	LarkClient *lark.Client // lark 客户端
	DocumentId string       // docx 文档 token
}

func NewDocxMarkdownProcessor(client *lark.Client, documentId string, opts ...Option) *DocxMarkdownProcessor {
	processor := DocxMarkdownProcessor{
		Config: &Config{
			StaticAsURL:  true,  // 默认不下载静态文件
			UseGhCallout: false, // 默认不使用 github 高亮块样式
		},
		LarkClient: client,
		DocumentId: documentId,
	}

	for _, opt := range opts {
		opt(&processor)
	}

	return &processor
}

type ProcessItem struct {
	BlockType      *int                // block 类型
	Normal         *larkdocx.Block     // 普通块
	Callout        []*larkdocx.Block   // 高亮块
	Table          [][]*larkdocx.Block // 表格
	QuoteContainer []*larkdocx.Block   // 引用容器
}

func (p *DocxMarkdownProcessor) DocxMarkdown(ctx context.Context) (string, error) {
	req := larkdocx.NewListDocumentBlockReqBuilder().DocumentId(p.DocumentId).Build()
	iterator, _ := p.LarkClient.Docx.V1.DocumentBlock.ListByIterator(ctx, req)

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
		case Callout:
			// 放入自己用于解析 emoji
			processItem.Callout = append(processItem.Callout, b)
			for _, cBlockId := range b.Children {
				st[cBlockId] = struct{}{}
				processItem.Callout = append(processItem.Callout, allBlockMap[cBlockId])
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

		case QuoteContainer:
			for _, cBlockId := range b.Children {
				st[cBlockId] = struct{}{}
				processItem.QuoteContainer = append(processItem.QuoteContainer, allBlockMap[cBlockId])
			}

		default:
			processItem.Normal = b
		}
		processItems = append(processItems, processItem)
	}

	// 转为 Markdown
	var buf = new(strings.Builder)
	for _, item := range processItems {
		buf.WriteString(p.DocxBlockMarkdown(ctx, item))
	}
	return buf.String(), nil
}

type Block struct {
	*larkdocx.Block
	ChildrenBlock []*Block
}

// listTransformToTree 递归填充节点的子节点
// https://open.feishu.cn/document/ukTMukTMukTM/uUDN04SN0QjL1QDN/document-docx/docx-v1/faq#424e7e9b
func listTransformToTree(ctx context.Context, rootLarkBlock *larkdocx.Block, larkBlockMap map[string]*larkdocx.Block) *Block {
	if rootLarkBlock == nil {
		return nil
	}

	root := &Block{Block: rootLarkBlock}

	for _, c := range root.Children {
		lb := larkBlockMap[c]
		b := &Block{
			Block: lb,
		}
		for _, bc := range b.Children {
			b.ChildrenBlock = append(b.ChildrenBlock, listTransformToTree(ctx, larkBlockMap[bc], larkBlockMap))
		}
		root.ChildrenBlock = append(root.ChildrenBlock, b)
	}

	return root
}
