package domain

import (
	"io/ioutil"
)

type Page struct {
	Title   string
	Content []byte
}

func (p *Page) Save() error {
	filename := p.Title + ".txt" // make text file
	return ioutil.WriteFile(filename, p.Content, 0600)
}

func LoadPage(title string) (*Page, error) {
	filename := title + ".txt"
	content, _ := ioutil.ReadFile(filename)
	return &Page{Title: title, Content: content}, nil
}
