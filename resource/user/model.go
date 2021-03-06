package user

const (
	Admin = "admin"
	Dev   = "dev"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
	Salt     string `json:"salt"`
}

type CreateCmd struct {
	Username string
	Password string
	Role     string
}

type UpdateCmd struct {
	UserID   string
	Username string
	Role     string
}

type GetCmd struct {
	UserID string
}

type DeleteCmd struct {
	UserID string
}

type GetAllCmd struct{}

type ValidateCmd struct {
	Username string
	Password string
}

var ValidUserID = "ae6ac8d6-0bcf-4671-a21a-49eab3167cbb"
var ValidUser = User{
	Username: "some username",
	Role:     Admin,
	Password: "some-hash",
	Salt:     "some-salt",
}
