package main

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
)

func main() {
	// 文章正文防xss注入
	var md = `

# 这是一级标题

> 这是应用

![xxx](jjdfnghjdf)
 <img src="x" onerror="alert(2)" alt="">
<script>alert(123)</script>
<Script>alert(121)</script>

`
	contentDoc, err := goquery.NewDocumentFromReader(bytes.NewReader([]byte(md)))
	if err != nil {
		fmt.Println(err)
		return
	}
	contentDoc.Find("script").Remove()
	contentDoc.Find("img").Remove()
	contentDoc.Find("iframe").Remove()

	content := contentDoc.Text()
	fmt.Println(content)
}
