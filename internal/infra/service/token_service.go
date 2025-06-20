package service

type TokenService interface {
	Generate(userID string) (string, error)
	Validate(token string) (string, error)
}
