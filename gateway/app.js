require('dotenv').config()
const express = require('express');
const path = require('path');
const cookieParser = require('cookie-parser');
const logger = require('morgan');
const rateLimit = require("express-rate-limit");
const helmet = require("helmet");
const xss = require("xss-clean");


const indexRouter = require('./routes/index');
const foodRouter = require('./routes/food');

// setup route middlewares
// Limit requests from same API
const limiter = rateLimit({
    max: 100,
    windowMs: 60 * 60 * 1000,
    message: "Too many requests from this IP, please try again in an hour!",
});
const hpp = require('hpp');
const app = express();

if (process.env.NODE_ENV === "development") {
    app.use(logger("dev"));
}
app.use(express.json());
app.use(express.urlencoded({ extended: false }));
app.use(cookieParser());
app.use(express.static(path.join(__dirname, 'public')));

app.use(hpp());
app.use(helmet());
app.use(xss());
app.use("/", limiter);


app.use('/', indexRouter);
app.use('/food', foodRouter);

module.exports = app;
