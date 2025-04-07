package referral_dto

type ReferralResponse struct {
	UserID        int                    `json:"user_id"`
	Name          string                 `json:"name"`
	Surname       string                 `json:"surname"`
	Telegram      string                 `json:"telegram"`
	PhotoPath     string                 `json:"photo_path"`
	ReferredUsers []ReferredUserResponse `json:"referred_users"`
}

type ReferredUserResponse struct {
	UserID    int    `json:"user_id"`
	Name      string `json:"name"`
	Surname   string `json:"surname"`
	Telegram  string `json:"telegram"`
	PhotoPath string `json:"photo_path"`
	CreatedAt string `json:"createdAt"`
}
