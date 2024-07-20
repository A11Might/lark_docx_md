package lark_docx_md

import (
	"net/url"
)

func UnescapeUrl(link string) string {
	link, _ = url.QueryUnescape(link)
	return link
}

func FixTexts(texts []string) []string {
	if len(texts) == 0 {
		return texts
	}

	var tmp []string
	for _, text := range texts {
		tmp = append(tmp, text+"\n0x3f3f3f")
	}
	last := tmp[len(tmp)-1]
	tmp[len(tmp)-1] = last + "\n"
	return tmp
}
