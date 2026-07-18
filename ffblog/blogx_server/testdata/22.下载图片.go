package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	//u := "https://secure2.wostatic.cn/static/tfLmFmcEqzBedocYFMSXuX/image.png?auth_key=1731249538-7BM6hJcunn2nKrE2YwcmiC-0-4de4827a24ee57dc07081a6c2a1e3384"
	//res, err := http.Get(u)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//fmt.Println(res.Status)

	//x := "介绍\n![](https://i1.hdslb.com/bfs/archive/832be5e18a614f7bcd93bf34c20ce690a044d656.jpg@672w_378h_1c_!web-home-common-cover.avif)\n介绍\n![](https://i1.hdslb.com/bfs/archive/832be5e18a614f7bcd93bf34c20ce690a044d656.jpg@672w_378h_1c_!web-home-common-cover.avif)\n"
	//regex := regexp.MustCompile(`!\[.*?\]\((.*?)\)`)
	//xx := regex.ReplaceAllStringFunc(x, func(s string) string {
	//	src := regex.FindStringSubmatch(s)
	//	fmt.Println(ss[1])
	//	return s
	//})
	//fmt.Println(xx)
	getImage()
}

func getImage() {
	url := "https://i1.hdslb.com/bfs/archive/832be5e18a614f7bcd93bf34c20ce690a044d656.jpg@672w_378h_1c_!web-home-common-cover.avif"

	response, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	byteData, _ := io.ReadAll(response.Body)
	os.WriteFile("xx.jpg", byteData, 0666)

}
