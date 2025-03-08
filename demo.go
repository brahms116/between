package demo

type User struct {
	Age      int       `json:"age"`
	Name     string    `json:"name"`
	Email    *string   `json:"email,omitEmpty"`
	Hobbies  *[]string `json:"hobbies,omitEmpty"`
	Status   Status    `json:"status"`
	UserData UserData  `json:"userData"`
}
type Status string

const Status_Active = "Active"
const Status_Disabled = "Disabled"
const Status_Pending = "Pending"

type UserData_Type string

const UserData_Type_adminData = "adminData"
const UserData_Type_customerData = "customerData"

type UserData struct {
	Type         UserData_Type `json:"_type"`
	AdminData    *AdminData    `json:"adminData,omitEmpty"`
	CustomerData *CustomerData `json:"customerData,omitEmpty"`
}
type AdminData struct {
	AccessLevel int `json:"accessLevel"`
}
type CustomerData struct {
	Attributes map[string]any `json:"attributes"`
}
