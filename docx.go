package lark_docx_md

import (
	"context"
	"strings"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkdocx "github.com/larksuite/oapi-sdk-go/v3/service/docx/v1"
	larkwiki "github.com/larksuite/oapi-sdk-go/v3/service/wiki/v2"
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
	Typ        string       // 文档类型，eg. docx, wiki
	Token      string       // 文档 token
}

func NewDocxMarkdownProcessor(client *lark.Client, typ, token string, opts ...Option) *DocxMarkdownProcessor {
	processor := DocxMarkdownProcessor{
		Config: &Config{
			StaticAsURL:  true,  // 默认不下载静态文件
			UseGhCallout: false, // 默认不使用 github 高亮块样式
		},
		LarkClient: client,
		Typ:        typ,
		Token:      token,
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
	switch p.Typ {
	case Docx:
		p.DocumentId = p.Token
	case Wiki:
		req := larkwiki.NewGetNodeSpaceReqBuilder().Token(p.Token).Build()
		resp, err := p.LarkClient.Wiki.V2.Space.GetNode(ctx, req)
		if err != nil {
			return "", err
		}
		p.DocumentId = *resp.Data.Node.ObjToken
	default:
	}

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

	allBlockMap := lo.SliceToMap(allBlock, func(item *larkdocx.Block) (string, *larkdocx.Block) {
		return *item.BlockId, item
	})
	root := p.listTransformToTree(ctx, allBlock[0], allBlockMap)

	// 转为 Markdown
	var buf = new(strings.Builder)
	buf.WriteString(strings.Join(p.DocxBlockMarkdown(ctx, root), "\n\n"))
	// 广告位
	buf.WriteString("\n\n***\n")
	buf.WriteString("_This MARKDOWN was generated with ❤️ by [lark_docx_md](https://github.com/A11Might/lark_docx_md)_")
	return strings.ReplaceAll(buf.String(), "0x3f3f3f\n\n", ""), nil
}

type Node struct {
	*larkdocx.Block
	ChildrenNode []*Node
}

// listTransformToTree 递归填充节点的子节点
// https://open.feishu.cn/document/ukTMukTMukTM/uUDN04SN0QjL1QDN/document-docx/docx-v1/faq#83024c0a
func (p *DocxMarkdownProcessor) listTransformToTree(ctx context.Context, rootLarkBlock *larkdocx.Block, larkBlockMap map[string]*larkdocx.Block) *Node {
	if rootLarkBlock == nil {
		return nil
	}

	root := &Node{Block: rootLarkBlock}

	for _, c := range root.Children {
		lb := larkBlockMap[c]
		root.ChildrenNode = append(root.ChildrenNode, p.listTransformToTree(ctx, lb, larkBlockMap))
	}

	return root
}
