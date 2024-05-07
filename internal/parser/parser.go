package parser

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/umtdemr/spor-istanbul-cli/internal/session"
	"golang.org/x/net/html"
	"io"
	"regexp"
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
		buttons := row.Find("a")
		if buttons.Length() < 2 {
			return
		}

		// get session post request id
		var postHref string

		for _, attr := range buttons.Get(0).Attr {
			if attr.Key != "href" {
				continue
			}

			re := regexp.MustCompile(`__doPostBack\('([^']+)'`)
			matches := re.FindStringSubmatch(attr.Val)

			if len(matches) < 2 {
				continue
			}

			postHref = matches[1]
		}

		// if we couldn't get the post request id, we can't go further
		if postHref == "" {
			return
		}

		subscriptions = append(subscriptions, &session.Subscription{
			Name:          row.Find("td").Get(2).FirstChild.Data,
			Date:          row.Find("td").Get(3).FirstChild.Data,
			Remaining:     row.Find("td").Get(5).FirstChild.Data,
			PostRequestId: postHref,
		})
	})

	return subscriptions
}
