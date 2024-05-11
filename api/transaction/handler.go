package transaction

import (
	"database/sql"
	"net/http"

	"github.com/KKGo-Software-engineering/workshop-summer/api/config"
	"github.com/labstack/echo/v4"
)

type Err struct {
	Message string `json:"message"`
}
type handler struct {
	flag config.FeatureFlag
	db   *sql.DB
}

func New(cfg config.FeatureFlag, db *sql.DB) *handler {
	return &handler{cfg, db}
}

const (
	cStmt = `INSERT INTO transaction (date, amount, category, transaction_type, note, image_url) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;`
)

func (h handler) Create(c echo.Context) error {

	// logger := mlog.L(c)
	// ctx := c.Request().Context()

	var tranReq TransactionRequest
	if err := c.Bind(tranReq); err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: "Invalid transaction request"})
	}

	// var lastInsertId int64
	// err := h.db.QueryRowContext(
	// 	ctx,
	// 	cStmt,
	// 	tranReq.Date, tranReq.Amount, tranReq.Category, tranReq.TransactionType, tranReq.Note, tranReq.ImageUrl,
	// ).Scan(&lastInsertId)

	// if err != nil {
	// 	logger.Error("query row error", zap.Error(err))
	// 	return c.JSON(http.StatusInternalServerError, err.Error())
	// }

	return c.JSON(http.StatusOK, "OK")
}

// func GetAll() {

// }

// func GetByExpenseId() {

// }
