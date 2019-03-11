package api

import (
	"encoding/json"

	"github.com/huyntsgs/go-rest-api/store"

	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bmizerany/assert"

	"time"

	"github.com/huyntsgs/go-rest-api/models"
)

var uid = 4
var users []models.User = []models.User{
	{Id: 1, UserName: "usr1", Password: "pass123", Email: "usr1@gmail.com", CreatedAt: time.Now()},
	{Id: 2, UserName: "usr2", Password: "pass123", Email: "usr2@gmail.com", CreatedAt: time.Now()},
	{Id: 3, UserName: "usr3", Password: "pass123", Email: "usr3@gmail.com", CreatedAt: time.Now()},
	{Id: 4, UserName: "usr4", Password: "pass123", Email: "usr4@gmail.com", CreatedAt: time.Now()},
}

type UserMock struct{}
type ErrUserMock struct{}

func (u *UserMock) Register(user *models.User) error {
	for _, usr := range users {
		if usr.Email == user.Email {
			return models.NewError("Email already registered", store.EMAIL_REGISTERED)
		}
	}
	users = append(users, *user)
	return nil
}
func (u *UserMock) Login(user *models.User) (*models.User, error) {
	for _, usr := range users {
		if usr.Email == user.Email && usr.Password == user.Password {
			return user, nil
		}
	}
	return nil, models.NewError("Invalid user email or password", store.INVALID_EMAIL_OR_PWD)
}
func (mock *ErrUserMock) Register(u *models.User) error {
	return models.NewError("Can not register with this email", ERR_INTERNAL)
}
func (mock *ErrUserMock) Login(u *models.User) (*models.User, error) {
	return nil, models.NewError("Invalid user email or password", store.INVALID_EMAIL_OR_PWD)
}

// register test funcs
func TestRegister(t *testing.T) {
	r := SetupUserRouter(&UserMock{})
	user := models.User{UserName: "usr10", Email: "usr10@gmail.com", Password: "pass10"}
	userJson, err := json.Marshal(user)
	req, err := MakeRequest("POST", "/api/v1/users/register", userJson, false)
	if err != nil {
		t.Fail()
	}
	res := httptest.NewRecorder()

	r.ServeHTTP(res, req)
	assert.Equal(t, http.StatusOK, res.Code)
}
func TestRegisterWithInvalidJson(t *testing.T) {
	r := SetupUserRouter(&UserMock{})
	userJson := `{UserNameNew: "usr10", Email: "usr10@gmail.com", Password: "pass123"}`

	req, err := MakeRequest("POST", "/api/v1/users/register", []byte(userJson), false)
	if err != nil {
		t.Fail()
	}
	res := httptest.NewRecorder()

	r.ServeHTTP(res, req)
	assert.Equal(t, http.StatusBadRequest, res.Code)
}
func TestRegisterWithInvalidInfo(t *testing.T) {
	r := SetupUserRouter(&UserMock{})
	userJson := `{UserName: "", Email: "usr10@gmail.com", Password: "pass10"}`

	req, err := MakeRequest("POST", "/api/v1/users/register", []byte(userJson), false)
	if err != nil {
		t.Fail()
	}
	res := httptest.NewRecorder()

	r.ServeHTTP(res, req)
	assert.Equal(t, http.StatusBadRequest, res.Code)
}
func TestRegisterWithInvalidEmail(t *testing.T) {
	r := SetupUserRouter(&UserMock{})
	userJson := `{UserName: "user11", Email: "usr10@gmail", Password: "pass10"}`

	req, err := MakeRequest("POST", "/api/v1/users/register", []byte(userJson), false)
	if err != nil {
		t.Fail()
	}
	res := httptest.NewRecorder()

	r.ServeHTTP(res, req)
	assert.Equal(t, http.StatusBadRequest, res.Code)
}

func TestRegisterWithErrorUserMock(t *testing.T) {
	r := SetupUserRouter(&ErrUserMock{})
	user := models.User{UserName: "usr1", Email: "usr1@gmail.com", Password: "pass10"}
	userJson, err := json.Marshal(user)
	req, err := MakeRequest("POST", "/api/v1/users/register", userJson, false)
	if err != nil {
		t.Fail()
	}
	res := httptest.NewRecorder()

	r.ServeHTTP(res, req)
	assert.Equal(t, http.StatusSeeOther, res.Code)
}

// login test funcs
func TestLogin(t *testing.T) {
	r := SetupUserRouter(&UserMock{})
	user := models.User{Email: "usr1@gmail.com", Password: "pass123"}
	userJson, _ := json.Marshal(user)
	req, err := MakeRequest("POST", "/api/v1/users/login", []byte(userJson), false)
	if err != nil {
		t.Fail()
	}
	res := httptest.NewRecorder()

	r.ServeHTTP(res, req)
	assert.Equal(t, http.StatusOK, res.Code)
}
func TestLoginWithInvalidInfo(t *testing.T) {
	r := SetupUserRouter(&UserMock{})
	userJson := `{Email1: "usr1@gmail.com", Password: "pass123"}`
	req, err := MakeRequest("POST", "/api/v1/users/login", []byte(userJson), false)
	if err != nil {
		t.Fail()
	}
	res := httptest.NewRecorder()

	r.ServeHTTP(res, req)
	assert.Equal(t, http.StatusBadRequest, res.Code)
}
func TestLoginWithInvalidPassword(t *testing.T) {
	r := SetupUserRouter(&UserMock{})
	user := models.User{Email: "usr1@gmail.com", Password: "1pass123"}
	userJson, _ := json.Marshal(user)
	req, err := MakeRequest("POST", "/api/v1/users/login", userJson, false)
	if err != nil {
		t.Fail()
	}
	res := httptest.NewRecorder()

	r.ServeHTTP(res, req)
	assert.Equal(t, http.StatusInternalServerError, res.Code)
}
func TestLoginWithErrorUserMock(t *testing.T) {
	r := SetupUserRouter(&ErrUserMock{})
	user := models.User{Email: "usr1@gmail.com", Password: "1pass123"}
	userJson, _ := json.Marshal(user)
	req, err := MakeRequest("POST", "/api/v1/users/login", userJson, false)
	if err != nil {
		t.Fail()
	}
	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)
	assert.Equal(t, http.StatusInternalServerError, res.Code)
}
