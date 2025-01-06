package models

type Manager struct {
	Id              string `json:"id"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	Name            string `json:"name"`
	UserImageUri    string `json:"userImageUri"`
	CompanyName     string `json:"companyName"`
	CompanyImageUri string `json:"companyImageUri"`
	Token           string `json:"token"`
}
