package referral_dto

type ReferrerResponse struct {
	UserID        int                    `json:"user_id"`
	Name          string                 `json:"name"`
	Surname       string                 `json:"surname"`
	Telegram      string                 `json:"telegram"`
	PhotoPath     string                 `json:"photo_path"`
	ReferrerID    int                    `json:"refer_id,omitempty"`
	WalletAddress string                 `json:"wallet_address,omitempty"`
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

type ChangeBalanceUserRequest struct {
	Amount int `json:"amount"`
}

type ChangeBalanceUserResponse struct {
	Success    bool   `json:"success"`
	Message    string `json:"message"`
	NewBalance int    `json:"newBalance"`
	Amount     int    `json:"amount"`
}

type BalanceResponse struct {
	Balance int `json:"balance"`
}

type DebitBalanceUserRequest struct {
	Balance int `json:"balance"`
}
