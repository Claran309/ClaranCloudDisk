package mysql

type TokenRepository interface {
	AddBlackList(token string) error
	CheckBlackList(token string) (string, error)
}
