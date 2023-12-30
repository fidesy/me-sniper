package models

type Action struct {
	Token *Token
}

func (a *Action) String() string {
	return a.Token.String()
}
