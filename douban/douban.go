package douban

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

//SearchResult 豆瓣移动版H5搜索结果
type SearchResult struct {
	Img   string
	Title string
	URL   string
	Type  string
	Rate  string
	Desc  string
}

//Detail 豆瓣电影详情数据
type Detail struct {
	Desc string
	URL  string
}

// Search 搜索豆瓣数据
func Search(keyword string, subjecttype string, fetchDetail bool) []SearchResult {
	url := fmt.Sprintf("https://m.douban.com/j/search/?q=%s&t=%s&p=0", keyword, subjecttype)
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Set("User-Agent", "Golang_Spider_Bot/3.0")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("network err:", err)
		return nil
	}
	buff, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	respJSON := make(map[string]interface{})
	doubanResults := make([]SearchResult, 0)
	err = json.Unmarshal(buff, &respJSON)
	if err != nil {
		fmt.Println("json Unmarshal error:", err)
		return nil
	}
	respHTML := respJSON["html"].(string)

	r := regexp.MustCompile(`(?sU:<li>.+href="(.+)".+src="(.+)".+class="subject-title">(.+)</span>.+<span>(.+)</span>.+</li>)`)
	results := r.FindAllStringSubmatch(respHTML, -1)
	for _, result := range results {
		doubanResult := SearchResult{Type: subjecttype, URL: "https://m.douban.com" + result[1], Img: result[2], Title: result[3], Rate: result[4]}
		doubanResults = append(doubanResults, doubanResult)
	}

	if fetchDetail { //爬取详情页面的描述信息
		fetchDetailForList(doubanResults)
	}

	return doubanResults
}

func fetchDetailForList(doubanResults []SearchResult) {
	detailChan := make(chan Detail, 0)
	for _, entry := range doubanResults { //启动协程并发获取
		go fetchDetail(entry.URL, detailChan)
	}
	detailMap := make(map[string]Detail)
	for i := 0; i < len(doubanResults); i++ {
		detail := <-detailChan
		detailMap[detail.URL] = detail
	}

	for i := 0; i < len(doubanResults); i++ {
		detail := detailMap[doubanResults[i].URL]
		doubanResults[i].Desc = detail.Desc
	}
}

func fetchDetail(detailURL string, detailChan chan Detail) {
	detail := Detail{URL: detailURL}
	resp, err := http.Get(detailURL)
	if err != nil {
		fmt.Println("获取详情页网络连接失败：", detailURL)
		detailChan <- detail
		return
	}
	buff, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	body := string(buff)

	r := regexp.MustCompile(`<meta name="description" content="(?s:(.+?))">`)
	results := r.FindStringSubmatch(body)
	if len(results) != 2 { //部分电影没有豆瓣评分
		detailChan <- detail
		return
	}
	detail.Desc = results[1]
	detailChan <- detail
}
