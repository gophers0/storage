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
		Success bool `json:"success"`
		Session struct {
			Id        int    `json:"id"`
			CreatedAt string `json:"created_at"`
			UpdatedAt string `json:"updated_at"`
			User      User   `json:"user"`
		} `json:"session"`
	}
)

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
