package trackermodel

type Tracker struct {
	Id        int64  `json:"id"`
	UserEmail string `json:"user_email"`
	CoinId    string `json:"coin_id"`
}
