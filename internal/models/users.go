package models

import (
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	HashedPassword string `json:"hashed_password"`
	CreatedAt      string `json:"created_at"`
}

type UserModel struct {
	DB *sql.DB
}

func (user *UserModel) Insert(name, email, password string) (int, error) {
	query := `insert into users (name, email, hashed_password, created_at) 
			  values (?, ?, ?, UTC_TIMESTAMP())`

	result, errQuery := user.DB.Exec(query, name, email, password)
	id, errResult := result.LastInsertId()

	if err := errors.Join(errQuery, errResult); err != nil {
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			// 1062 is the error number for duplicate entry
			if mySQLError.Number == 1062 {
				return 0, ErrDuplicateEmail
			}
		}
		return 0, err
	}

	return int(id), nil
}

func (user *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPassword []byte

	usr := &User{}

	query := `select id, hashed_password from users where email = ?`
	err := user.DB.QueryRow(query, email).Scan(&usr.Email, &usr.HashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		}
		return 0, err
	}

	hashedPassword = []byte(usr.HashedPassword)

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		}
		return 0, err
	}

	return id, nil
}

func (user *UserModel) Get(id int) (*User, error) {
	return nil, nil
}

func (user *UserModel) Exist(id int) (bool, error) {
	return false, nil
}
