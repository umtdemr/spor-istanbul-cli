package parser

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/umtdemr/spor-istanbul-cli/internal/session"
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

func (p *Parser) GetSubscriptions(r io.Reader) []*session.Subscription {
	doc, _ := goquery.NewDocumentFromReader(r)

	var subscriptions []*session.Subscription

	doc.Find("#dtUyeSpor tr").Each(func(index int, row *goquery.Selection) {
		// ignore header of the table
		if index == 0 {
			return
		}

		// there should be a select session button.
		if row.Find("a").Length() < 2 {
			return
		}

		subscriptions = append(subscriptions, &session.Subscription{
			Name:      row.Find("td").Get(2).FirstChild.Data,
			Date:      row.Find("td").Get(3).FirstChild.Data,
			Remaining: row.Find("td").Get(5).FirstChild.Data,
		})
	})

	return subscriptions
}
