package dto

import "oceanlearn.learn/ginessential/model"

type UserDto struct {
	Name string `json:"name"`
	Telephone string `json:"telephone"`
}

//只返回用户名和密码，其他字段省略
func ToUserDto(user model.User) UserDto{
	return UserDto{
		Name: user.Name,
		Telephone: user.Telephone,
	}
}