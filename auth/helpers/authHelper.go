package helpers

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
)

func CheckUserType(c *gin.Context, role string) (err error) {
	userType := c.GetString("user_type")
	err = nil
	if userType != role {
		err = errors.New("Unauthorized")
		return err
	}
	return err
}

func MatchUserTypeToUid(c *gin.Context, userId string) (err error) {
	userType := c.GetString("user_type")
	fmt.Println(userType)
	uid := c.GetString("uid")
	fmt.Println(uid)
	fmt.Println(userId)

	if userType == "USER" && uid != userId {
		err = errors.New("Unauthorized")
		return err
	}
	err = CheckUserType(c, userType)
	return err
}
