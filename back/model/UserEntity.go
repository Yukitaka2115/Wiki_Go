package model

type User struct {
	ID       int
	Role     int //1,admin; 2,commonUsers
	Username string
	Pwd      string
}
