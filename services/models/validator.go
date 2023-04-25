package models

type User struct {
	UserID     string       `json:"user_id"`
	Validators []*Validator `json:"validators"`
}

type Validator struct {
	Name        string `json:"name"`
	Address     string `json:"address"`
	NotifyEvery int64  `json:"notify_every"`
}

func NewValidator(name string, addess string) *Validator {
	return &Validator{
		Name:        name,
		Address:     addess,
		NotifyEvery: DefaultNotifyEvery,
	}
}
