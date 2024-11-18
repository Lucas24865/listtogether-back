package requests

type UserUpdateRequest struct {
	User    string
	NewPass string
	Color   string
	Name    string
	OldPass string
}
