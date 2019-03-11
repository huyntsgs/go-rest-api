package store

import (
	"testing"

	"github.com/bmizerany/assert"
	"github.com/huyntsgs/go-rest-api/models"
)

// test with register funcs
func TestRegister(t *testing.T) {
	u := models.User{UserName: "user1", Email: "user1@gmail.com", Password: "pass123"}
	mysql := InitDB()
	mysql.DeleteUsers()
	err := mysql.Register(&u)

	assert.Equal(t, nil, err)
}

func TestRegisterWithEmailRegistered(t *testing.T) {
	u := models.User{UserName: "user1", Email: "user1@gmail.com", Password: "pass123"}

	mysql := InitDB()
	mysql.DeleteUsers()
	err := mysql.Register(&u)

	err = mysql.Register(&u)

	if err == nil {
		t.Fail()
	}

}

func TestRegisterInvalidUserInfo(t *testing.T) {
	u := models.User{UserName: "11", Email: "user1@gmail", Password: "pass123"}
	mysql := InitDB()
	err := mysql.Register(&u)
	if err == nil {
		t.Fail()
	}
}

// test with login funcs
func TestLogin(t *testing.T) {
	u := models.User{UserName: "user1", Email: "user1@gmail.com", Password: "pass123"}
	mysql := InitDB()
	mysql.DeleteUsers()
	err := mysql.Register(&u)
	assert.Equal(t, nil, err)
	_, err = mysql.Login(&u)
	assert.Equal(t, nil, err)
}

func TestLoginWithInvalidData(t *testing.T) {
	u := models.User{UserName: "user1", Email: "user1@gmail.com", Password: "pass123"}
	mysql := InitDB()
	mysql.DeleteUsers()
	err := mysql.Register(&u)
	assert.Equal(t, nil, err)
	u.Email = "user1@gmail."
	_, err = mysql.Login(&u)
	if err == nil {
		t.Fail()
	}
}
func TestLoginWithWrongPass(t *testing.T) {
	u := models.User{UserName: "user1", Email: "user1@gmail.com", Password: "pass123"}
	mysql := InitDB()

	mysql.DeleteUsers()
	err := mysql.Register(&u)
	assert.Equal(t, nil, err)
	u.Password = "pass12"
	_, err = mysql.Login(&u)
	if err == nil {
		t.Fail()
	}
}
