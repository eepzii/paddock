package validate

type LoginResponse struct {
	SessionId           string `json:"SessionId"`
	PasswordIsTemporary bool   `json:"PasswordIsTemporary"`
	Subscriber          struct {
		FirstName   string `json:"FirstName"`
		LastName    string `json:"LastName"`
		HomeCountry string `json:"HomeCountry"`
		Id          int    `json:"Id"`
		Email       string `json:"Email"`
		Login       string `json:"Login"`
	} `json:"Subscriber"`
	Country string `json:"Country"`
	Data    struct {
		SubscriptionStatus string `json:"subscriptionStatus"`
		SubscriptionToken  string `json:"subscriptionToken"`
	} `json:"data"`
}
