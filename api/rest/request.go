package rest

type addUserRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
