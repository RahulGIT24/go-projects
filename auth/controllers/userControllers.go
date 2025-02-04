package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	// "strconv"
	"time"

	"github.com/rahulgit24/golang-projects/auth/database"
	helper "github.com/rahulgit24/golang-projects/auth/helpers"
	"github.com/rahulgit24/golang-projects/auth/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
var validate = validator.New()

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true
	msg := ""

	if err != nil {
		msg = fmt.Sprintf("Invalid Password")
		check = false
	}

	return check, msg
}

func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(user)

		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})

		defer cancel()

		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while validating mail"})
			return
		}

		password := HashPassword(*user.Password)
		user.Password = &password
		count, phoneErr := userCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		defer cancel()

		if phoneErr != nil {
			log.Panic(phoneErr)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while verfying phone"})
			return
		}

		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Email or phone number already exist"})
			return
		}

		user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_id = user.ID.Hex()
		token, refreshToken, _ := helper.GenerateAllTokens(*user.Email, *user.First_name, *user.Last_name, *user.User_type, *&user.User_id)
		user.Token = &token
		user.Refresh_token = &refreshToken

		resultInsertionNumber, insertErr := userCollection.InsertOne(ctx, user)
		if insertErr != nil {
			msg := fmt.Sprintf("User item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, resultInsertionNumber)
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User
		var foundUser models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		passwordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
		defer cancel()

		if passwordIsValid == false {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		if foundUser.Email == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
			return
		}

		token, refreshToken, _ := helper.GenerateAllTokens(*foundUser.Email, *foundUser.First_name, *foundUser.Last_name, *foundUser.User_type, *&foundUser.User_id)
		helper.UpdateAllTokens(token, refreshToken, foundUser.User_id)

		userErr := userCollection.FindOne(ctx, bson.M{"user_id": foundUser.User_id}).Decode(&foundUser)

		if userErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": userErr.Error()})
			return
		}

		c.JSON(http.StatusAccepted,gin.H{"accesstoken":token,"refreshtoken":refreshToken})
	}
}

// func GetUsers() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		if err := helper.CheckUserType(c, "ADMIN"); err != nil {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
// 			return
// 		}
// 		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

// 		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))

// 		if err != nil || recordPerPage < 1 {
// 			recordPerPage = 10
// 		}

// 		page, pageErr := strconv.Atoi(c.Query("page"))
// 		if pageErr != nil || page < 1 {
// 			page = 1
// 		}

// 		startIndex := (page - 1) * recordPerPage
// 		startIndex, err = strconv.Atoi(c.Query("startIndex"))
// 		matchStage := bson.D{{"$match", bson.D{}}}
// 		groupStage := bson.D{
// 			{"$group", bson.D{
// 				{"_id", "null"},
// 				{"total_count", bson.D{{"$sum", 1}}},
// 				{"data", bson.D{{"$push": "$$ROOT"}}},
// 			}},
// 		}

// 	}
// }

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("userid")
		if err := helper.MatchUserTypeToUid(c, userId); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var user models.User

		err := userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user)
		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}
