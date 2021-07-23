package domain

type StoreAccessToken string

func (t StoreAccessToken) String() string {
	return string(t)
}
