package lark_docx_md

import (
	"net/url"
	"strings"
)

func UnescapeUrl(link string) string {
	link, _ = url.QueryUnescape(link)
	return link
}

// 删除字符串后两个\n\n
func FixText(text string) string {
	if len(text) < 2 {
		return text
	}

	return strings.TrimSuffix(text, "\n\n")
}

// 删除最后一个元素的后两个\n\n
func FixTexts(texts []string) []string {
	if len(texts) == 0 {
		return texts
	}

	length := len(texts)
	last := texts[length-1]
	last = FixText(last)
	texts[length-1] = last
	return texts
}
