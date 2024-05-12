package service

import (
	"github.com/umtdemr/spor-istanbul-cli/internal/client"
	"github.com/umtdemr/spor-istanbul-cli/internal/parser"
	"github.com/umtdemr/spor-istanbul-cli/internal/session"
	"strings"
)

type Service struct {
	client *client.Client
	parser *parser.Parser
}

func NewService() *Service {
	return &Service{
		client: client.NewClient(),
		parser: parser.NewParser(),
	}
}

func (s *Service) Login(id string, password string) bool {
	body := s.client.Login(id, password)
	title, ok := s.parser.GetTitle(body)

	if !ok {
		return false
	}

	return !strings.Contains(title, "Giriş Yap")

}

func (s *Service) GetSubscriptions() []*session.Subscription {
	body := s.client.GetSubscriptionsPage()
	return s.parser.GetSubscriptions(body)
}

func (s *Service) GetSessions(postRequestId string) []*session.Collection {
	body := s.client.GetSessions(postRequestId)
	return s.parser.ParseSessionsDoc(body)
}

func (s *Service) CheckSessionApplicable(postRequestId string, sessionId string) bool {
	sessions := s.GetSessions(postRequestId)

	for _, collection := range sessions {
		for _, singleSession := range collection.Sessions {
			if singleSession.Id == sessionId && singleSession.Applicable {
				return true
			}
		}
	}

	return false
}
