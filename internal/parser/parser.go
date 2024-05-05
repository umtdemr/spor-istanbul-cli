package parser

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
	"io"
)

type Parser struct {
}

func NewParser() *Parser {
	return &Parser{}
}

func isTitleElement(n *html.Node) bool {
	return n.Type == html.ElementNode && n.Data == "title"
}

func traverse(n *html.Node) (string, bool) {
	if isTitleElement(n) {
		return n.FirstChild.Data, true
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result, ok := traverse(c)
		if ok {
			return result, ok
		}
	}

	return "", false
}

func (p *Parser) GetTitle(r io.Reader) (string, bool) {
	doc, err := html.Parse(r)
	if err != nil {
		panic("Fail to parse html")
	}

	return traverse(doc)
}

func (p *Parser) GetSubscriptions(r io.Reader) {
	doc, _ := goquery.NewDocumentFromReader(r)

	doc.Find("#dtUyeSpor tr").Each(func(index int, d *goquery.Selection) {
		fmt.Println(index)
	})
}
