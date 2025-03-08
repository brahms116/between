# Between

A small language to describe DTOs and sync them between languages.

> [!NOTE]  
> Very early prototype, many and all things subject to change as well as many feature gaps

## Example

Given the following file `demo.bt`

```bt
prod User {
  age Int,
  name Str,
  email Str?,
  hobbies []?Str,
  Status,
  UserData,
}

sumstr Status {
  "Active",
  "Disabled",
  "Pending",
}

sum UserData {
  AdminData,
  CustomerData,
}

prod AdminData {
  accessLevel Int,
}

prod CustomerData {
  attributes Object,
}
```

Generates the following TypeScript code

```ts
interface User {
  age: number;
  name: string;
  email?: string;
  hobbies?: string[];
  status: Status;
  userData: UserData;
}
type Status = "Active" | "Disabled" | "Pending";
type UserData =
  | { _type: "adminData"; adminData: AdminData }
  | { _type: "customerData"; customerData: CustomerData };
interface AdminData {
  accessLevel: number;
}
interface CustomerData {
  attributes: Record<string, unknown>;
}
```

and the following Go code

```go
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
```
