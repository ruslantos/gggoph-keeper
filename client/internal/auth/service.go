package auth

//
//import (
//	"context"
//	"encoding/json"
//	"fmt"
//	"net/http"
//
//	"goph-keeper/client/internal/client"
//)
//
//type Service struct {
//	apiClient *client.APIClient
//}
//
//func New(apiURL string) *Service {
//	return &Service{
//		apiClient: client.New(apiURL),
//	}
//}
//
//func (s *Service) Register(ctx context.Context, username, password string) error {
//	req := map[string]string{
//		"username": username,
//		"password": password,
//	}
//
//	resp, err := s.apiClient.Post(ctx, "/register", req)
//	if err != nil {
//		return fmt.Errorf("registration failed: %w", err)
//	}
//	defer resp.Body.Close()
//
//	if resp.StatusCode != http.StatusCreated {
//		var errResp struct{ Error string }
//		json.NewDecoder(resp.Body).Decode(&errResp)
//		return fmt.Errorf(errResp.Error)
//	}
//
//	return nil
//}
//
//func (s *Service) Login(ctx context.Context, username, password string) (string, error) {
//	req := map[string]string{
//		"username": username,
//		"password": password,
//	}
//
//	resp, err := s.apiClient.Post(ctx, "/login", req)
//	if err != nil {
//		return "", fmt.Errorf("login failed: %w", err)
//	}
//	defer resp.Body.Close()
//
//	var result struct {
//		Token string `json:"token"`
//	}
//	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
//		return "", err
//	}
//
//	return result.Token, nil
//}
