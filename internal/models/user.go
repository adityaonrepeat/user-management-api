package models

import "time"

const DateLayout = "2006-01-02"

type User struct {
	ID   int32
	Name string
	DOB  time.Time
}

type CreateUserRequest struct {
	Name string `json:"name" validate:"required"`
	DOB  string `json:"dob"  validate:"required,datetime=2006-01-02"`
}

type UpdateUserRequest struct {
	Name string `json:"name" validate:"required"`
	DOB  string `json:"dob"  validate:"required,datetime=2006-01-02"`
}

type ListParams struct {
	Paginated bool
	Limit     int32
	Offset    int32
}

type UserResponse struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
	DOB  string `json:"dob"`
}

type UserWithAgeResponse struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
	DOB  string `json:"dob"`
	Age  int    `json:"age"`
}
