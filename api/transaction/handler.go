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

func (h handler) Create(c echo.Context) error {
	var tranReq TransactionRequest
	if err := c.Bind(tranReq); err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: "Invalid transaction request"})
	}

	return c.JSON(http.StatusOK, "OK")
}

func GetAll() {

}

func GetByExpenseId() {

}
