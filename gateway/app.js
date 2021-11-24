const express = require('express');
const path = require('path');
const cookieParser = require('cookie-parser');
const logger = require('morgan');
const csrf = require('csurf')

const indexRouter = require('./routes/index');
const usersRouter = require('./routes/users');

// setup route middlewares
const hpp = require('hpp');
const csrfProtection = csrf({ cookie: true })
const parseForm = bodyParser.urlencoded({ extended: false })

const app = express();

app.use(logger('dev'));
app.use(express.json());
app.use(express.urlencoded({ extended: false }));
app.use(cookieParser());
app.use(express.static(path.join(__dirname, 'public')));

app.use(hpp());

app.use('/', indexRouter);
app.use('/users', usersRouter);

module.exports = app;
