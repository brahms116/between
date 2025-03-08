package main

type User struct {
	Name      string     `json:"name"`
	Age       *int       `json:"age,omitEmpty"`
	Gender    *string    `json:"gender,omitEmpty"`
	FriendIds *[]*string `json:"friendIds,omitEmpty"`
	Array2d   []*[]*int  `json:"array2d"`
	Status    Status     `json:"status"`
	Data      Data       `json:"data"`
}
type Status string

const Status_NOT_HERE = "NOT_HERE"
const Status_REALLY = "REALLY"

type Data_Type string

const Data_Type_adminData = "adminData"
const Data_Type_userData = "userData"

type Data struct {
	Type      Data_Type  `json:"_type"`
	AdminData *AdminData `json:"adminData,omitEmpty"`
	UserData  *UserData  `json:"userData,omitEmpty"`
}
type AdminData struct {
	Level int `json:"level"`
}
type UserData struct {
	Id   int            `json:"id"`
	Data map[string]any `json:"data"`
}
