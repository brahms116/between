package demo

type User struct {
	Age      int       `json:"age"`
	Name     string    `json:"$name"`
	Email    *string   `json:"email,omitEmpty"`
	Hobbies  *[]string `json:"hobbies,omitEmpty"`
	Status   Status    `json:"status"`
	UserData UserData  `json:"userData"`
}
type Status string

const Status_Active Status = "Active"
const Status_Disabled Status = "Disabled"
const Status_Pending Status = "pending activation"

type UserData_Type string

const UserData_Type_adminData UserData_Type = "adminData"
const UserData_Type_customerData UserData_Type = "customerData"

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
