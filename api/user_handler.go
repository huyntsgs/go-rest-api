package api

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/huyntsgs/go-rest-api/models"
	"github.com/huyntsgs/go-rest-api/utils"
)

type UserHandler struct {
	userStore UserStore
}

// Creates new UserHandler.
// UserHandler accepts interface UserStore.
// Any data store implements UserStore could be the input of the handle.
func NewUserHandler(userStore UserStore) UserHandler {
	return UserHandler{userStore}
}

// Login handles login router.
// Function validates parameters and call Login from UserStore.
func (h UserHandler) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		err := c.BindJSON(&user)
		if err != nil {
			GinAbort(c, http.StatusBadRequest, INVALID_PARAMS, "")
			return
		}
		// Validate user data
		if !user.ValidateLogin() {
			GinAbort(c, http.StatusBadRequest, INVALID_PARAMS, "")
			return
		}
		u, errc := h.userStore.Login(&user)
		if errc != nil {
			GinAbort(c, http.StatusInternalServerError, errc.(*models.ErrorC).ErrCode, errc.(*models.ErrorC).ErrMsg)
			return
		}
		// Return token for authorizing other apis
		claimInfo := map[string]interface{}{
			"userId":   u.Id,
			"userName": u.UserName,
			"email":    u.Email,
		}
		// Generates token with expire time is 24 hours
		token, _ := utils.GenToken(claimInfo, []byte(os.Getenv("TOKEN_KEY")), 24*60)
		u.Password = ""
		c.JSON(http.StatusOK, gin.H{
			"status": "success", "res": u, "token": token,
		})
	}
}

// Register handles register router.
// Function validates parameters and call Register from UserStore.
func (h UserHandler) Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		err := c.BindJSON(&user)
		if err != nil {
			GinAbort(c, http.StatusBadRequest, INVALID_PARAMS, "")
			return
		}

		// Validate user data
		if !(user.Validate()) {
			GinAbort(c, http.StatusBadRequest, INVALID_PARAMS, "")
			return
		}
		errc := h.userStore.Register(&user)
		if errc != nil {
			log.Println(errc)
			GinAbort(c, http.StatusSeeOther, errc.(*models.ErrorC).ErrCode, errc.(*models.ErrorC).ErrMsg)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "sucess",
		})
	}
}
