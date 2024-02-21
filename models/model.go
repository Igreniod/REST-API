package models

type User struct {
	Uid    string `json:"Uid,omitempty"`
	Nama   string `json:"Nama,omitempty" validate:"required"`
	Email  string `json:"Email,omitempty" validate:"required"`
	Alamat string `json:"Alamat,omitempty" validate:"required"`
}
