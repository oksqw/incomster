package core

type Session struct {
	ID     int    `json:"id"     required:"true"`
	UserID int    `json:"userId" required:"true"`
	Token  string `json:"token"  required:"true"`
}

type SessionCreateInput struct {
	UserID int    `json:"userId" required:"true"`
	Token  string `json:"token"  required:"true"`
}

type SessionUpdateInput struct {
	Id     int    `json:"id"     required:"true"`
	UserID int    `json:"userId" required:"true"`
	Token  string `json:"token"  required:"true"`
}

type SessionGetInput struct {
	Id     int    `json:"id"     required:"true"`
	UserID int    `json:"userId" required:"false"`
	Token  string `json:"token"  required:"false"`
}
