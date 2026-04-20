package main

import (
	"errors"
	"fmt"
	"time"
)

// ============================================
// STRUCT - Like a PHP class without methods
// ============================================
type User struct {
	ID        int    `json:"id"`   // `json:"..."` = serialization tag
	Name      string `json:"name"` // Exported (public)
	email     string `json:"-"`    // Unexported (private), - = ignore in JSON
	CreatedAt int64  `json:"created_at"`
}

// ============================================
// CONSTRUCTOR - Go doesn't have constructors
// ============================================
func NewUser(name string) *User {
	return &User{
		Name:      name,
		CreatedAt: time.Now().Unix(),
	}
	// & means returning pointer (reference)
	// Without &, returns a COPY of the struct
}

// ============================================
// METHOD - Function attached to struct
// ============================================
// (u *User) = receiver - like $this in PHP
// * means pointer receiver - can modify the struct
func (u *User) SetName(name string) {
	u.Name = name
}

// Without * = value receiver - works on a copy
func (u User) GetName() string {
	return u.Name
}

// ============================================
// INTERFACE - Implicit implementation
// ============================================
// In PHP: class User implements UserRepository
// In Go: Any type with these methods automatically satisfies the interface
type Repository interface {
	GetByID(id int) (*User, error)
	Save(user *User) error
	Delete(id int) error
}

// ============================================
// ERROR HANDLING - No exceptions!
// ============================================
func GetUser(id int) (*User, error) {
	if id <= 0 {
		return nil, errors.New("invalid id")
		// nil = null in Go
		// errors.New() = throw new Exception() in PHP
	}

	user := &User{ID: id, Name: "John"}
	return user, nil
	// Always return (result, error) - caller must check
}

// ============================================
// CALLING CODE - Must handle errors explicitly
// ============================================
func main() {
	user, err := GetUser(1)
	if err != nil {
		// Like try-catch, but explicit
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(user.Name)
}
