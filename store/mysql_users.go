package store

import (
	//"log"

	"github.com/huyntsgs/go-rest-api/models"
	"golang.org/x/crypto/bcrypt"
)

// DeleteUsers deletes all users in table users
// Used for test only
func (mySql *MySqlDB) DeleteUsers() error {
	// Check whether email already registered
	rows, _ := mySql.DB.Query("DELETE from users where user_id > 0")
	defer rows.Close()
	return nil
}

// DeleteProducts deletes all products in table products
// Used for test only
func (mySql *MySqlDB) DeleteProducts() error {
	// Check whether email already registered
	rows, _ := mySql.DB.Query("DELETE from products where product_id > 0")
	defer rows.Close()
	return nil
}

// Register inserts user with given username, email and password to users table.
// Function checks the existence of email and hashing password.
// Function returns nil if register process is successful.
func (mySql *MySqlDB) Register(user *models.User) error {

	if !user.Validate() {
		return models.NewError("Invalid user data", INVALID_DATA)
	}

	// Check whether email already registered
	rows, err := mySql.DB.Query("SELECT count(user_id) as cnt from users where email = LOWER(?)", user.Email)
	if err != nil {
		return models.NewError("Internal server error", ERR_INTERNAL)
	}
	defer rows.Close()
	n := 0
	for rows.Next() {
		rows.Scan(&n)
		if n > 0 {
			return models.NewError("The email already registered", EMAIL_REGISTERED)
		}
	}
	// Hash password before insert to db
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.NewError("Internal server error", ERR_INTERNAL)
	}

	query, err := mySql.DB.Exec("INSERT into users (user_name, email, password) VALUES (?,?,?)", user.UserName, user.Email, string(hashedPwd))
	if err != nil {
		return models.NewError("Internal server error", ERR_INTERNAL)
	}
	rowsCnt, err := query.RowsAffected()
	if rowsCnt == 0 {
		return models.NewError("Internal server error", ERR_INTERNAL)
	}
	return nil
}

// Login signs user in with provide user email and password.
// Function returns error if email or password are not matched.
func (mySql *MySqlDB) Login(user *models.User) (*models.User, error) {
	if !user.ValidateLogin() {
		return nil, models.NewError("Invalid data", INVALID_DATA)
	}
	rows, err := mySql.DB.Query("SELECT user_id, user_name, email, password, created_at from users where email = LOWER(?) LIMIT 1", user.Email)
	if err != nil {
		return nil, models.NewError("Internal server error", ERR_INTERNAL)
	}
	defer rows.Close()
	var u *models.User
	for rows.Next() {
		u = new(models.User)
		err := rows.Scan(&u.Id, &u.UserName, &u.Email, &u.Password, &u.CreatedAt)
		if err != nil {
			return nil, models.NewError("Internal server error", ERR_INTERNAL)
		}
		break
	}
	if u != nil {
		if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(user.Password)); err != nil {
			return nil, models.NewError("Invalid user name or password", INVALID_EMAIL_OR_PWD)
		}
		return u, nil
	}
	return nil, models.NewError("Invalid user name or password", INVALID_EMAIL_OR_PWD)
}
