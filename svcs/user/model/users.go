package model

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"strconv"
	"time"
)

var (
	ErrNoCustomerInResponse = errors.New("Response has no matching customer")
	ErrMissingField         = "Error missing %v"
)

// User 用户结构
type User struct {
	FirstName string        `json:"firstName" bson:"firstName"`
	LastName  string        `json:"lastName" bson:"lastName"`
	Email     string        `json:"-" bson:"email"`
	Username  string        `json:"username" bson:"username"`
	Password  string        `json:"-" bson:"password,omitempty"`
	UserID    string        `json:"id" bson:"-"`
	Salt      string        `json:"-" bson:"salt"`
	Authority UserAuthority `json:"authority" bson:"authority"`
}

// New ..
func New() User {
	u := User{}
	u.NewSalt()
	return u
}

// Validate ..
func (u *User) Validate() error {
	if u.FirstName == "" {
		return fmt.Errorf(ErrMissingField, "FirstName")
	}
	if u.LastName == "" {
		return fmt.Errorf(ErrMissingField, "LastName")
	}
	if u.Username == "" {
		return fmt.Errorf(ErrMissingField, "Username")
	}
	if u.Password == "" {
		return fmt.Errorf(ErrMissingField, "Password")
	}
	return nil
}

// NewSalt ..
func (u *User) NewSalt() {
	h := sha1.New()
	io.WriteString(h, strconv.Itoa(int(time.Now().UnixNano())))
	u.Salt = fmt.Sprintf("%x", h.Sum(nil))
}

// Failer is an interface that should be implemented by response types.
// Response encoders can check if responses are Failer, and if so if they've
// failed, and if so encode them using a separate write path based on the error.
type Failer interface {
	Failed() error
}

// GetUserRequest collects the request parameters for the GetProducts method.
type GetUserRequest struct {
	A string
}

// GetUserResponse collects the response values for the GetProducts method.
type GetUserResponse struct {
	V   User  `json:"v"`
	Err error `json:"err,omitempty"` // should be intercepted by Failed/errorEncoder
}

// Failed implements Failer.
func (r GetUserResponse) Failed() error { return r.Err }

// RegisterRequest struct
type RegisterRequest struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

// RegisterUserResponse ...
type RegisterUserResponse struct {
	ID  string `json:"id"`
	Err error  `json:"-"`
}

// LoginRequest ..
type LoginRequest struct {
	Username string
	Password string
}

// LoginResponse ..
type LoginResponse struct {
	User  *User  `json:"user,omitempty"`
	Token string `json:"token,omitempty"`
	Err   error  `json:"err,omitempty"`
}
