const express = require("express");
const router = express.Router();

const foodHandler = require("./inventory");

router.get("/", foodHandler.getFood);


module.exports = router;
