package structs

import (
	"time"
)

type Food struct {
	Food_id    string    `bson:"food_id,omitempty" json:"food_id,omitempty"`
	Name       *string   `bson:"name,omitempty" json:"name" validate:"required,min=2,max=100"`
	Price      *float64  `bson:"price,omitempty" json:"price" validate:"required"`
	Food_image *string   `bson:"food_image,omitempty" json:"food_image" validate:"required"`
	CreatedAt  time.Time `bson:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt  time.Time `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}
