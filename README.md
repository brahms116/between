# Between

A small language to describe DTOs and sync them between languages.

> [!WARNING]  
> Very early prototype, many and all things subject to change as well as many feature gaps

## Example

Given the following file `demo.bt`

```bt
prod User {
  age Int,
  name "$name" Str,
  email Str?,
  hobbies []?Str,
  dateOfBirth Date,
  Status,
  UserData,
}

sumstr Status {
  Active,
  Disabled,
  Pending "pending activation",
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
export interface User {
  age: number;
  $name: string;
  email?: string;
  hobbies?: string[];
  dateOfBirth: string;
  status: Status;
  userData: UserData;
}
export type Status = "Active" | "Disabled" | "pending activation";
export type UserData =
  | { adminData: AdminData }
  | { customerData: CustomerData };
export interface AdminData {
  accessLevel: number;
}
export interface CustomerData {
  attributes: Record<string, unknown>;
}
```

and the following Go code

```go
package demo

import (
	"time"
)

type User struct {
	Age         int       `json:"age"`
	Name        string    `json:"$name"`
	Email       *string   `json:"email,omitEmpty"`
	Hobbies     *[]string `json:"hobbies,omitEmpty"`
	DateOfBirth time.Time `json:"dateOfBirth"`
	Status      Status    `json:"status"`
	UserData    UserData  `json:"userData"`
}
type Status string

const Status_Active Status = "Active"
const Status_Disabled Status = "Disabled"
const Status_Pending Status = "pending activation"

type UserData struct {
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

## TODOs

- [x] Use `runes` instead of `bytes`
- [x] Split ast and codegen input. Tokens should be part of the AST so that it can contain sugar and locations so it can be used for formatting.
- [ ] Add support for doc comments
- [ ] Add type checking
- [ ] Add LSP support
- [ ] Add support for more languages

## Installation

```sh
go install github.com/brahms116/between/cmd/bt@latest
```

## Usage

```sh
bt --input ./demo.bt --output ./result.go && gofmt -w ./result.go
```

or

```sh
bt --input ./demo.bt --output ./result.ts && prettier --write ./result.ts
``` 
