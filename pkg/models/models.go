package models

type Accounts struct {
	Account_id int    `json:"accountid" gorm:"primaryKey"`
	First_name string `json:"firstname" validate:"required"`
	Last_name  string `json:"lastname" validate:"required"`
	Username   string `json:"username" validate:"required"`
	Password   string `json:"password" validate:"required"`
}

type Contacts struct {
	Id         int    `json:"id"`
	Account_id int    `json:"accountid"`
	Email      string `json:"email"  validate:"omitempty,email"`
	Contact    string `json:"contact"  validate:"omitempty,number,max=11"`
}

type Address struct {
	Id          int    `json:"id"`
	Account_id  int    `json:"accountid"`
	House_no    string `json:"houseno"`
	Street      string `json:"street"`
	Subdivision string `json:"subdivision"`
	Barangay    string `json:"barangay"`
	City        string `json:"city"`
	Province    string `json:"province"`
	Country     string `json:"country"`
	Zip_code    string `json:"zipcode"`
}
