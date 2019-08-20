package main

import (
	"douban-crawer/alfred"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
)

//DoubanSuggestResultEntry 豆瓣搜索框返回的suggest数据结构
type DoubanSuggestResultEntry struct {
	Episode   string
	Img       string
	Title     string
	URL       string
	EntryType string `json:"type"`
	Year      string
	SubTitle  string `json:"sub_title"`
	ID        string
	Rate      string
	Desc      string
}

//DoubanMovieDetail 豆瓣电影详情数据
type DoubanMovieDetail struct {
	Rate string
	Desc string
	ID   string
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("请输入要搜索的电影名字")
		return
	}
	var query = os.Args[1]
	var items *alfred.Items = &alfred.Items{}

	for _, entry := range searchDouban(query) {
		items.AppendItem(alfred.Item{Title: entry.Title + " " + entry.Rate, SubTitle: "上映时间：" + entry.Year + " " + entry.SubTitle + " " + entry.Desc, Arg: entry.URL})
	}
	xmlBytes, _ := xml.MarshalIndent(items, "", "	")
	fmt.Println(string(xmlBytes))
}

func searchDouban(keyword string) []DoubanSuggestResultEntry {
	url := "https://movie.douban.com/j/subject_suggest?q="
	resp, err := http.Get(url + keyword)
	if err != nil {
		fmt.Println("network err:", err)
		return nil
	}
	buff, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	entrys := make([]DoubanSuggestResultEntry, 0)
	err = json.Unmarshal(buff, &entrys)
	if err != nil {
		fmt.Println("json Unmarshal error,", err)
	}

	details := make(chan DoubanMovieDetail, len(entrys))
	for _, entry := range entrys {
		go getEntryDetail(entry.ID, details)
	}
	detailMap := make(map[string]DoubanMovieDetail)
	for i := 0; i < len(entrys); i++ {
		detail := <-details
		detailMap[detail.ID] = detail
	}

	for i := 0; i < len(entrys); i++ {
		detail := detailMap[entrys[i].ID]
		entrys[i].Desc = detail.Desc
		entrys[i].Rate = detail.Rate
	}
	return entrys
}

func getEntryDetail(movieID string, details chan DoubanMovieDetail) {
	detail := DoubanMovieDetail{ID: movieID}
	detailURL := "https://m.douban.com/movie/subject/" + movieID + "/"
	resp, err := http.Get(detailURL)
	if err != nil {
		fmt.Println("获取详情页网络连接失败：", detailURL)
		details <- detail
		return
	}
	buff, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	body := string(buff)

	r := regexp.MustCompile(`<meta name="description" content=".*：(\S+)\s*简介：(.+)">`)
	results := r.FindStringSubmatch(body)
	if len(results) != 3 { //部分电影没有豆瓣评分
		details <- detail
		return
	}
	detail.Rate = results[1]
	detail.Desc = results[2]
	//fmt.Println(detail)
	details <- detail
	return
}
