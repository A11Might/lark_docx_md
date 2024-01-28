package lark_docx_md

import "net/url"

func UnescapeUrl(link string) string {
	link, _ = url.QueryUnescape(link)
	return link
}
