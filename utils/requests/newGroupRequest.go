package requests

type NewGroupRequest struct {
	Desc   string
	Name   string
	Id     string
	Admins []string
	Users  []string
}
