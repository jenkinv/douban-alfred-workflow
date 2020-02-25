package main

import (
	"flag"
	"fmt"
	"strconv"

	"github.com/jenkinv/douban-alfred-workflow/alfred"
	"github.com/jenkinv/douban-alfred-workflow/douban"
)

func main() {
	var query, querytype string
	var detail bool
	flag.StringVar(&querytype, "t", "movie", "`type` can be movie, music or book")
	flag.BoolVar(&detail, "d", false, "if fetch `detail` data extraly, set true")
	flag.Parse()

	query = flag.Arg(0)

	var items = &alfred.Items{}
	var stars = "★★★★★☆☆☆☆☆"
	for _, entry := range douban.Search(query, querytype, detail) {
		item := alfred.Item{}
		item.Title = fmt.Sprintf("%s（%s）", entry.Title, entry.Rate)
		item.SubTitle = entry.Desc
		if rate, err := strconv.ParseFloat(entry.Rate, 32); err == nil && entry.Desc == "" {
			var halfRate = int((rate + 0.5) / 2)
			item.SubTitle = stars[(5-halfRate)*3 : (10-halfRate)*3]
		}
		item.Arg = entry.URL
		item.Icon = entry.Type + ".png"
		items.AppendItem(item)
	}
	if items.Length() == 0 {
		items.AppendItem(alfred.Item{
			Title:    "没有搜索到结果",
			SubTitle: "有可能是爬虫不稳定，请联系作者github",
			Arg:      "",
			Icon:     "",
		})
	}
	//output xml format for alfred
	fmt.Println(items.ToXML())
	return
}
