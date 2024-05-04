package service

import "github.com/umtdemr/spor-istanbul-cli/internal/client"

type Service struct {
	client *client.Client
}

func NewService() *Service {
	return &Service{
		client: client.NewClient(),
	}
}

func (s *Service) Login(id string, password string) bool {
	return s.client.Login(id, password)
}

func (s *Service) GetSubscriptions() {
	s.client.GetSubscriptionsPage()
}
