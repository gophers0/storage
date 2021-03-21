package users

const (
	UserRoleAdmin = "admin"
)

type (
	CheckTokenRequest struct {
		UserId int    `json:"user_id"`
		Token  string `json:"token"`
	}
	CheckTokenResponse struct {
		Success bool    `json:"success"`
		Session Session `json:"session"`
	}
)

type Session struct {
	Id        int    `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	User      User   `json:"user"`
	Token     string `json:"token"`
}

type User struct {
	Id        int    `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Login     string `json:"login"`
	Role      string `json:"role"`
}

type (
	SearchUserRequest struct {
		Login string `json:"login"`
	}
	SearchUserResponse struct {
		Code    int    `json:"code"`
		Count   int    `json:"count"`
		Records []User `json:"records"`
	}
)
