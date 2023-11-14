package controller

import (
	"encoding/json"
	"fmt"

	"frostik.com/auth/constants"
	"frostik.com/auth/mapper"
	"frostik.com/auth/model"
	"frostik.com/auth/util"
	"github.com/allegro/bigcache/v3"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetUserByEmail(ctx *gin.Context, mongoClient *mongo.Client, cacheClient *bigcache.BigCache, email *string, role *string, noCache bool) (*model.StudentPopulated, *string) {
	var student model.Student
	var studentPopulated model.StudentPopulated

	// Check if copy is there in the cache
	if !noCache {
		studentBytes, _ := cacheClient.Get(*email)
		if err := json.Unmarshal(studentBytes, &studentPopulated); err == nil {
			fmt.Println("Retreiving the student data from the cache")
			return &studentPopulated, nil
		}
	}

	// Query to DB
	fmt.Println("Queriying the DB for User Details")
	mongoClient.Database(constants.DB).Collection(constants.COLLECTION_STUDENT).FindOne(ctx, bson.M{
		"email": *email,
	}).Decode(&student)
	studentPopulated = mapper.TransformStudentToStudentPopulated(student)

	var groupIds = []primitive.ObjectID{}
	var groupDetails = []model.Group{}
	for _, id := range student.Groups {
		groupIds = append(groupIds, id)
	}

	cursor, _ := mongoClient.Database(constants.DB).Collection(constants.COLLECTION_GROUP).Find(ctx, bson.M{
		"_id": bson.M{"$in": groupIds},
	})
	cursor.All(ctx, &groupDetails)
	studentPopulated.Groups = groupDetails

	// Now check if it is actually a student by the ROLES
	if !util.CheckRoleExists(&groupDetails, *role) {
		return nil, &constants.ERROR_NOT_A_STUDENT
	}

	// Set to bigCache
	studentBytes, _ := json.Marshal(studentPopulated)
	if err := cacheClient.Set(*email, studentBytes); err == nil {
		fmt.Println("Successfully set UserDetails in cache")
	}
	return &studentPopulated, nil
}
