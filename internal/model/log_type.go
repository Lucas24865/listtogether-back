package model

const (
	LoginType = "login"
)

type LogType string

func (n LogType) GetLogTypes() []string {
	return []string{LoginType}
}

func (n LogType) ToString() string {
	return string(n)
}
