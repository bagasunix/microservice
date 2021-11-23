package router

import (
	"context"
	"fmt"
	"inventory/config"
	"inventory/structs"
	"math"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

func GetFood(c *gin.Context) {

	//this is used to determine how long the API call should last
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	//declare a variable of type food
	var food structs.Food

	//bind the object that comes in with the declared varaible. thrrow an error if one occurs
	if err := c.BindJSON(&food); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// use the validation packge to verify that all items coming in meet the requirements of the struct
	validationErr := validate.Struct(food)
	if validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		return
	}

	// assing the time stamps upon creation
	food.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	food.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	//generate new ID for the object to be created
	food.ID = primitive.NewObjectID()

	// assign the the auto generated ID to the primary key attribute
	food.Food_id = food.ID.Hex()
	var num = ToFixed(*food.Price, 2)
	food.Price = &num

	//insert the newly created object into mongodb
	result, insertErr := foodCollection.InsertOne(ctx, food)
	if insertErr != nil {
		msg := fmt.Sprintf("Food item was not created")
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}
	defer cancel()

	//return the id of the created object to the frontend
	c.JSON(http.StatusOK, result)
}
