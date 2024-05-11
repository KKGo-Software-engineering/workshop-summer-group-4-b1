package transaction

import (
	"database/sql"
	"net/http"
	"strconv"

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

func (h *handler) GetTransactionById(c echo.Context) error {
	// Retrieve spenderID as a string and convert to integer
	spenderIDStr := c.Param("id")
	spenderID, err := strconv.Atoi(spenderIDStr)
	if err != nil {
		// Return an error if conversion fails
		c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid spender ID"})
		return err
	}

	var transactions []TransactionResponse

	// Use the integer spenderID in the SQL query
	rows, err := h.db.Query(`
        SELECT id, date, amount, category, transaction_type, note, image_url
        FROM transaction
        WHERE spender_id = $1`, spenderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, echo.Map{"error": "Database error"})
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var t TransactionResponse
		if err := rows.Scan(&t.ID, &t.Date, &t.Amount, &t.Category, &t.TransactionType, &t.Note, &t.ImageUrl); err != nil {
			c.JSON(http.StatusInternalServerError, echo.Map{"error": "Error scanning database results"})
			return err
		}
		transactions = append(transactions, t)
	}

	if err = rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, echo.Map{"error": "Error iterating database results"})
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{"transactions": transactions})
}
