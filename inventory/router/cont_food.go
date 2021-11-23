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
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	foodData.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	foodData.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	//generate new ID for the object to be created
	foodData.ID = primitive.NewObjectID()

	// assign the the auto generated ID to the primary key attribute
	foodData.Food_id = foodData.ID.Hex()
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

	var careers []structs.Food
	if err = cursor.All(context.TODO(), &careers); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		c.Abort()
		return
	}

	if len(careers) == 0 {
		careers = []structs.Food{}
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   careers,
	})

}
