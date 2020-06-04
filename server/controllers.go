package server

import (
	"net/http"
	"paccount/models"
	"strconv"
	"time"

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
	var input models.Account
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.DocNumber <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid number"})
		return

	}

	timeStamp := time.Now()
	model := &models.Account{
		DocNumber: input.DocNumber,
		Limit:     input.Limit,
		CreatedAt: int32(timeStamp.Unix()),
		//UpdatedAt: int32(timeStamp.Unix()),
	}

	if err := s.DB.Create(model); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "can't create account"})
		return
	}

	c.JSON(http.StatusCreated, model)
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

	var account models.Account
	if err := s.DB.FindOne(models.Account{ID: uit}, &account); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "record not found"})
		return
	}

	c.JSON(http.StatusOK, account)
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
	var input models.Transaction
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var acc models.Account
	if err := s.DB.FindOne(models.Account{ID: input.AccountID}, &acc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.OperationID < 4 && input.Amount > acc.Limit {
		c.JSON(http.StatusBadRequest, gin.H{"error": "account not have a limit"})
		return
	}

	//fmt.Printf("%#v\n", acc)

	if _, ok := s.OprType[input.OperationID]; !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "operation type not exist"})
		return
	}

	switch input.OperationID {
	case 1, 2, 3:
		// -
		acc.Limit -= input.Amount
	case 4:
		// +
		acc.Limit += input.Amount
	}
	timeStamp := time.Now()
	//acc.UpdatedAt = int32(timeStamp.Unix())

	if err := s.DB.Update(acc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	timeStamp = time.Now()
	model := &models.Transaction{
		Amount:      input.Amount,
		OperationID: input.OperationID,
		AccountID:   input.AccountID,
		CreatedAt:   int32(timeStamp.Unix()),
	}

	if err := s.DB.Create(model); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "can't create transaction"})
		return
	}

	c.JSON(http.StatusCreated, model)
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

	var model models.Transaction
	var txs []models.Transaction

	if err := s.DB.GetDB().Table(model.TableName()).
		Where("account_id = ?", uit).
		Find(&txs).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "record(s) not found"})
		return
	}

	c.JSON(http.StatusOK, txs)
}
