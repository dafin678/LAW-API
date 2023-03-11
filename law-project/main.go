package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type ATM struct {
	AvailableCash float64
	Accounts      map[string]float64
}

func main() {
	atm := ATM{
		Accounts: map[string]float64{
			"0001": 500.0,
			"0002": 1000.0,
			"0003": 750.0,
		},
	}
	router := gin.Default()

	v1 := router.Group("/v1")
	{
		v1.GET("/accounts/:account", atm.checkBalance)
		v1.POST("/accounts/:account/withdraw", atm.withdraw)
		v1.POST("/accounts/:account/deposit", atm.deposit)
		v1.POST("/accounts/:account/transfer", atm.transfer)
	}
	router.Run(":9888")
}

func (atm *ATM) checkBalance(c *gin.Context) {
	account := c.Param("account")
	balance, exists := atm.Accounts[account]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"balance": balance})
}

func (atm *ATM) withdraw(c *gin.Context) {
	account := c.Param("account")
	var req struct {
		Amount float64 `form:"amount"`
	}
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	amount := req.Amount

	if amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Amount must be greater than zero"})
		return
	}

	balance, exists := atm.Accounts[account]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}

	if balance < amount {
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient balance"})
		return
	}

	atm.Accounts[account] -= amount

	c.JSON(http.StatusOK, gin.H{"message": "Withdrawal successful", "balance": atm.Accounts[account]})
}

func (atm *ATM) deposit(c *gin.Context) {
	account := c.Param("account")
	var req struct {
		Amount float64 `form:"amount"`
	}
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	amount := req.Amount

	if amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Amount must be greater than zero"})
		return
	}

	_, exists := atm.Accounts[account]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}

	atm.Accounts[account] += amount

	c.JSON(http.StatusOK, gin.H{"message": "Deposit successful", "balance": atm.Accounts[account]})
}

type TransferRequest struct {
	ToAccount string
	Amount    float64
}

func (atm *ATM) transfer(c *gin.Context) {
	fromAccount := c.Param("account")
	var req TransferRequest
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	toAccount := req.ToAccount
	amount := req.Amount
	if amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Amount must be greater than zero"})
		return
	}

	fromBalance, exists := atm.Accounts[fromAccount]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}

	if fromBalance < amount {
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient balance"})
		return
	}

	_, exists = atm.Accounts[toAccount]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}

	atm.Accounts[fromAccount] -= amount
	atm.Accounts[toAccount] += amount

	c.JSON(http.StatusOK, gin.H{"message": "Transfer successful", "balance": atm.Accounts[fromAccount]})
}
