package alfred

import "encoding/xml"

// Item 为 Alfred 输出格式中的元素
type Item struct {
	XMLName  xml.Name `xml:"item"`
	Title    string   `xml:"title"`
	SubTitle string   `xml:"subtitle"`
	Arg      string   `xml:"arg,attr"`
}

// Items 为 Alfred 输出格式中的数组
type Items struct {
	XMLName xml.Name `xml:"items"`
	Items   []Item   `xml:"item"`
}

// AppendItem 增加列表元素
func (items *Items) AppendItem(item ...Item) {
	items.Items = append(items.Items, item...)
}

// ToXML 输出为XML格式
func (items *Items) ToXML() string {
	xmlBytes, _ := xml.MarshalIndent(items, "", "	")
	return string(xmlBytes)
}
