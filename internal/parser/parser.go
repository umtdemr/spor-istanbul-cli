package parser

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/umtdemr/spor-istanbul-cli/internal/session"
	"golang.org/x/net/html"
	"io"
	"regexp"
	"strings"
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

func (p *Parser) GetSubscriptions(buffer *bytes.Buffer) ([]*session.Subscription, string) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(buffer.Bytes()))

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

	viewState, _ := p.GetViewState(bytes.NewReader(buffer.Bytes()))
	return subscriptions, viewState
}

func (p *Parser) ParseSessionsDoc(r io.Reader) []*session.Collection {
	doc, _ := goquery.NewDocumentFromReader(r)

	var sessionCollections []*session.Collection

	// find all the panels
	sessionQuery := doc.Find("#dvScheduler > div")

	sessionQuery.Each(func(index int, panelNode *goquery.Selection) {
		if index == sessionQuery.Length()-1 {
			return
		}

		// get the panel heading
		heading := strings.TrimSpace(
			strings.ReplaceAll(
				strings.Replace(
					panelNode.Find(".panel-heading").Text(),
					"\n",
					"",
					2,
				),
				" ",
				"",
			),
		)

		// since the format is DAY\nDATE, we need to split it
		headingParts := strings.Split(heading, "\n")

		if len(headingParts) != 2 {
			return
		}

		// create the collection
		sessionCollection := &session.Collection{
			Day:      headingParts[0],
			Date:     headingParts[1],
			Sessions: []*session.Session{},
		}

		// get the sessions for this collection
		panelNode.Find(".panel-body .well").Each(func(sessionGroupIdx int, sessionNode *goquery.Selection) {
			sessionId, ok := sessionNode.Attr("id")
			if !ok {
				return
			}
			input := sessionNode.Find("input[type='checkbox']")

			sessionCollection.Sessions = append(sessionCollection.Sessions, &session.Session{
				Available:  strings.TrimSpace(sessionNode.Find(".label-success").Text()),
				Time:       strings.TrimSpace(sessionNode.Find("span[id*='lblSeansSaat']").Text()),
				Id:         sessionId,
				Applicable: input.Length() > 0,
			})
		})

		// do not add if there is no session
		if len(sessionCollection.Sessions) == 0 {
			return
		}

		sessionCollections = append(sessionCollections, sessionCollection)
	})

	return sessionCollections
}

// GetViewState gets view state value from input
func (p *Parser) GetViewState(r io.Reader) (string, error) {
	var err error

	doc, _ := goquery.NewDocumentFromReader(r)

	viewState := ""
	selector := fmt.Sprintf("input[name=%s]", "\"__VIEWSTATE\"")

	doc.Find(selector).Each(func(i int, selection *goquery.Selection) {
		value, exists := selection.Attr("value")
		if exists {
			viewState = value
		}
	})

	if viewState == "" {
		err = errors.New("could not extract viewstate")
	}

	return viewState, err
}
