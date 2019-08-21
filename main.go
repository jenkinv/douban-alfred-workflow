package main

import (
	"flag"
	"fmt"

	"github.com/jenkinv/douban-movie-alfredworkflow/alfred"
	"github.com/jenkinv/douban-movie-alfredworkflow/douban"
)

func main() {
	var query, querytype string
	var detail bool
	flag.StringVar(&querytype, "t", "movie", "`type` can be movie, music or book")
	flag.BoolVar(&detail, "d", false, "if fetch `detail` data extraly, set true")
	flag.Parse()

	query = flag.Arg(0)

	var items = &alfred.Items{}
	for _, entry := range douban.Search(query, querytype, detail) {
		item := alfred.Item{}
		item.Title = entry.Title + " " + entry.Rate
		item.SubTitle = entry.Desc
		item.Arg = entry.URL
		items.AppendItem(item)
	}
	//output xml format for alfred
	fmt.Println(items.ToXML())
	return
}
