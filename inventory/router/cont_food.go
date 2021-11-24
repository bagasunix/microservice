package router

import (
	"context"
	"inventory/config"
	"inventory/structs"
	"math"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

//connect to to the database and open a food collection
var foodCollection *mongo.Collection = config.OpenCollection(config.Client, "food")

// create a validator object
var validate = validator.New()

//this function rounds the price value down to 2 decimal places
func ToFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(Round(num*output)) / output
}
func Round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func InputFood(c *gin.Context) {
	//this is used to determine how long the API call should last
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	//declare a variable of type food
	var foodData structs.Food
	//bind the object that comes in with the declared varaible. thrrow an error if one occurs
	if err := c.BindJSON(&foodData); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "can't bind struct",
		})
		c.Abort()
		return
	}
	// use the validation packge to verify that all items coming in meet the requirements of the struct
	validationErr := validate.Struct(foodData)
	if validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		return
	}

	// assing the time stamps upon creation
	foodData.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	foodData.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	// assign the the auto generated ID to the primary key attribute
	foodData.Food_id = uuid.NewString()
	var num = ToFixed(*foodData.Price, 2)
	foodData.Price = &num
	//insert the newly created object into mongodb
	result, insertErr := foodCollection.InsertOne(ctx, foodData)
	if insertErr != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": insertErr.Error(),
		})
		c.Abort()
		return
	}
	defer cancel()
	//return the id of the created object to the frontend
	c.JSON(http.StatusOK, result)
}

func GetFood(c *gin.Context) {
	cursor, err := foodCollection.Find(context.Background(), bson.M{})
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		c.Abort()
		return
	}

	var foodData []structs.Food
	if err = cursor.All(context.TODO(), &foodData); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		c.Abort()
		return
	}

	if len(foodData) == 0 {
		foodData = []structs.Food{}
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   foodData,
	})
}

func EditFood(c *gin.Context) {
	id := c.Param("id")
	var foodData structs.Food

	if err := c.BindJSON(&foodData); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "can't bind struct",
		})
		c.Abort()
		return
	}

	// add updated_at, updated_by
	foodData.UpdatedAt = time.Now()

	// update process
	pByte, err := bson.Marshal(foodData)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		c.Abort()
		return
	}

	var update bson.M
	err = bson.Unmarshal(pByte, &update)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		c.Abort()
		return
	}

	_, err = foodCollection.UpdateOne(context.Background(), bson.M{"food_id": id}, bson.M{"$set": update})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "food edited successfully",
	})
}

func DeleteFood(c *gin.Context) {
	id := c.Param("id")

	_, err := foodCollection.DeleteOne(context.Background(), bson.M{"food_id": id})
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		c.Abort()
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "food deleted successfully",
	})
}

func GetFoodById(c *gin.Context) {
	id := c.Param("id")
	var foodData structs.Food

	if err := foodCollection.FindOne(context.Background(), bson.M{"food_id": id}).Decode(&foodData); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		c.Abort()
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   foodData,
	})
}
