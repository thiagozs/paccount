package server

import (
	"net/http"
	"paccount/pkg/account"
	"paccount/pkg/transaction"
	"strconv"

	"github.com/gin-gonic/gin"
)

// NoRoute wildcard for no route
func (s *Server) NoRoute(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{"code": http.StatusNotFound, "message": "route not found"})
}

// Welcome for home access api
func (s *Server) Welcome(c *gin.Context) {
	c.JSON(200, map[string]string{"message": "welcome to code challenge"})
}

// Ping healthcheck api endpoint
func (s *Server) Ping(c *gin.Context) {
	c.String(200, "pong")
}

// CreateAccount controller create Account
// @Summary Create a new Account
// @Description This method you will create a new Account
// @Accept json
// @Produce json
// @Success 201 {object} models.Account
// @Router /accounts [post]
// @Param account body models.Account true "Account"
func (s *Server) CreateAccount(c *gin.Context) {
	var input account.Entity
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.DocNumber <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid number"})
		return
	}

	a := account.New(s.DB)
	acc, err := a.Create(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "fail on create account"})
		return
	}

	c.JSON(http.StatusCreated, acc)
}

// FindAccount controller find account by ID
// @Summary Get information about accounts
// @Description Get a JSON with search by ID
// @Accept  json
// @Produce  json
// @Param id path int true "ID"
// @Success 200 {object} models.Account
// @Router /accounts/{id} [get]
func (s *Server) FindAccount(c *gin.Context) {

	uit, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid params"})
		return
	}

	a := account.New(s.DB)
	acc, err := a.Find(uit)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "record not found"})
		return

	}

	c.JSON(http.StatusOK, acc)
}

// CreateTx controller create new transaction
// @Summary Create a new Account
// @Description This method you will create a new transaction
// @Accept json
// @Produce json
// @Success 201 {object} models.Account
// @Router /transactions [post]
// @Param transaction body models.Transaction true "Transaction"
func (s *Server) CreateTx(c *gin.Context) {
	var input transaction.Entity
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tx := transaction.New(s.DB)

	newtx, err := tx.ProcessTx(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, newtx)
}

// FindAllTxsByAccount controller find all txs by account ID
// @Summary Get information about transactions by account ID
// @Description Get a JSON with search by ID
// @Accept  json
// @Produce  json
// @Param id path int true "ID"
// @Success 200 {array} models.Transaction
// @Router /transactions/account/{id} [get]
func (s *Server) FindAllTxsByAccount(c *gin.Context) {

	uit, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid params"})
		return
	}

	tx := transaction.New(s.DB)
	txs, err := tx.GetAllTxs(uit)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, txs)
}
