package alfred

import "encoding/xml"

type Item struct {
	XMLName  xml.Name `xml:"item"`
	Title    string   `xml:"title"`
	SubTitle string   `xml:"subtitle"`
	Arg      string   `xml:"arg,attr"`
}

type Items struct {
	XMLName xml.Name `xml:"items"`
	Items   []Item   `xml:"item"`
}

func (items *Items) AppendItem(item ...Item) {
	items.Items = append(items.Items, item...)
}
