const express = require("express");
const router = express.Router();

const foodHandler = require("./inventory");

router
    .route("/")
    .get(foodHandler.getFood)
    .post(foodHandler.createFood)

router
    .route("/:id")
    .patch(foodHandler.updateFood)
    .delete(foodHandler.deleteFood)
    .get(foodHandler.getByFood)


module.exports = router;
