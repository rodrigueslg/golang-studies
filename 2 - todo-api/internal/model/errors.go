package model

import "errors"

var ErrConnectDatabase = errors.New("failed to connect to the database")
var ErrPreparingStatemant = errors.New("failed preparing statemant")
var ErrExecuteQuery = errors.New("failed to execute query")
var ErrInsertingRow = errors.New("failed inserting row")
var ErrScanningRows = errors.New("failed scanning row")

var ErrInvalidTaskStatus = errors.New("invalid task status")
var ErrInvalidTaskId = errors.New("invalid task id")
var ErrInvalidTaskName = errors.New("invalid task name")
var ErrTaskAlreadyExists = errors.New("task already exists")
var ErrTaskNotFound = errors.New("task not found")

var ErrUnauthorized = errors.New("invalid token")
var ErrInvalidRequestBody = errors.New("invalid request body")
var ErrInternalServerError = errors.New("internal server error")
