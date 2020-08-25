package response

import "ginweb/model"

type UserResponse struct {
	Name string `json:"name"`
	Telephone string `json:"telephone"`
}

func ToUserResponse(user model.User) UserResponse{
	return UserResponse{
		Name:	user.Name,
		Telephone: user.Telephone,
	}
}