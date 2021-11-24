const express = require("express");
const router = express.Router();

const foodHandler = require("./inventory");

router
    .route("/")
    .get(foodHandler.getFood)
    .post(foodHandler.createFood)


module.exports = router;
