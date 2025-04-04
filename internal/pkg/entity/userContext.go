package entity

type userContextKey string

const UserContextKey = userContextKey("user_context")

type UserContext struct {
	UserID int
	Role   int
}
