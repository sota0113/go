package main

import (
	"fmt"
	"io/ioutil"
	"encoding/json"
)

type Page struct {
    Title    string  `json:"Title"`
    Body  []byte `json:"Body,string"`
}

type PageForJson struct {
    Title    string  `json:"Title"`
    Body  string `json:"Body"`
}


func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page,error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func main () {

	p1 := &Page{Title: "TestPage",Body: []byte("this is body")}
	p1.save()
	//p2,_ := loadPage(p1.Title)

	json_str := `
{
"Title": "`+p1.Title+`",
"Body" : "`+string(p1.Body)+`"
}
`
	var p3 PageForJson
	err2 := json.Unmarshal([]byte(json_str),&p3)
	if err2 != nil {
		fmt.Println(err2)
		return
	}
	str_json,_ := json.Marshal(p3)
	fmt.Println(string(str_json))
}
