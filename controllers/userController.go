package controllers

import (
	"gorm.io/gorm"
	"myapp/auth"
	"net/http"
	"time"

	"myapp/database"
	"myapp/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Register(c *gin.Context) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}

	user.ID = uuid.New()
	user.CreatedDate = time.Now()

	result := database.DB.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Phone Number already registered"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "SUCCESS",
		"result": user,
	})
}

func Login(c *gin.Context) {
	var request struct {
		PhoneNumber string `json:"phone_number"`
		PIN         string `json:"pin"`
	}

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}

	var user models.User
	err := database.DB.Where("phone_number = ?", request.PhoneNumber).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Phone Number and PIN doesn't match"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		}
		return
	}

	if user.PIN != request.PIN {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Phone Number and PIN doesn't match"})
		return
	}

	// Generate access token
	accessToken, err := auth.GenerateJWT(user.ID.String(), user.PhoneNumber, "access")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error generating access token"})
		return
	}

	// Generate refresh token
	refreshToken, err := auth.GenerateJWT(user.ID.String(), user.PhoneNumber, "refresh")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error generating refresh token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "SUCCESS",
		"result": gin.H{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		},
	})
}

func TopUp(c *gin.Context) {
	var request struct {
		Amount float64 `json:"amount"`
	}
	// Parse JWT token
	claims, err := auth.ParseJWT(c.Request.Header.Get("Authorization"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthenticated"})
		return
	}

	// Bind JSON request to struct
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request", "error": err.Error()})
		return
	}

	// Retrieve user from database
	var user models.User
	if err := database.DB.Where("phone_number = ?", claims.PhoneNumber).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		}
		return
	}

	// Perform top-up operation
	previousBalance := user.Balance
	user.Balance += request.Amount

	// Save updated user balance
	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update balance", "error": err.Error()})
		return
	}

	// Create top-up record
	topUp := models.TopUp{
		UserID:        user.ID,
		Amount:        request.Amount,
		BalanceBefore: previousBalance,
		BalanceAfter:  user.Balance,
		CreatedDate:   time.Now(),
	}
	if err := database.DB.Create(&topUp).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create top-up record", "error": err.Error()})
		return
	}

	// Prepare response
	response := gin.H{
		"status": "SUCCESS",
		"result": gin.H{
			"top_up_id":      topUp.ID,
			"amount_top_up":  topUp.Amount,
			"balance_before": previousBalance,
			"balance_after":  user.Balance,
			"created_date":   topUp.CreatedDate.Format("2006-01-02 15:04:05"),
		},
	}

	c.JSON(http.StatusOK, response)
}

func Payment(c *gin.Context) {
	var request struct {
		Amount  float64 `json:"amount"`
		Remarks string  `json:"remarks"`
	}

	// Parse JWT token
	claims, err := auth.ParseJWT(c.Request.Header.Get("Authorization"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthenticated"})
		return
	}

	// Bind JSON request to struct
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request", "error": err.Error()})
		return
	}

	// Retrieve user from database
	var user models.User
	err = database.DB.Where("phone_number = ?", claims.PhoneNumber).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		}
		return
	}

	// Check if user has enough balance
	if user.Balance < request.Amount {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Balance is not enough"})
		return
	}

	// Perform payment operation
	previousBalance := user.Balance
	user.Balance -= request.Amount

	// Save updated user balance
	err = database.DB.Save(&user).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update balance"})
		return
	}

	// Create payment record
	payment := models.Payment{
		UserID:        user.ID,
		Amount:        request.Amount,
		Remarks:       request.Remarks,
		BalanceBefore: previousBalance,
		BalanceAfter:  user.Balance,
		CreatedDate:   time.Now(),
	}
	err = database.DB.Create(&payment).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create payment record"})
		return
	}

	// Prepare response
	response := gin.H{
		"status": "SUCCESS",
		"result": gin.H{
			"payment_id":     payment.ID,
			"amount":         payment.Amount,
			"remarks":        payment.Remarks,
			"balance_before": payment.BalanceBefore,
			"balance_after":  payment.BalanceAfter,
			"created_date":   payment.CreatedDate.Format("2006-01-02 15:04:05"),
		},
	}

	c.JSON(http.StatusOK, response)
}

func Transfer(c *gin.Context) {
	var request struct {
		TargetUser uuid.UUID `json:"target_user"`
		Amount     float64   `json:"amount"`
		Remarks    string    `json:"remarks"`
	}

	// Parse JWT token
	claims, err := auth.ParseJWT(c.Request.Header.Get("Authorization"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthenticated"})
		return
	}

	// Bind JSON request to struct
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request", "error": err.Error()})
		return
	}

	// Create response channel
	responseChan := make(chan gin.H)

	// Process transfer in the background
	go processTransfer(request, claims.PhoneNumber, responseChan)

	// Receive the response from the background goroutine
	response := <-responseChan
	c.JSON(http.StatusOK, response)
}

func processTransfer(request struct {
	TargetUser uuid.UUID `json:"target_user"`
	Amount     float64   `json:"amount"`
	Remarks    string    `json:"remarks"`
}, senderPhoneNumber string, responseChan chan gin.H) {
	var transfer models.Transfer

	// Retrieve sender (from user) from database
	var fromUser models.User
	err := database.DB.Where("phone_number = ?", senderPhoneNumber).First(&fromUser).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			responseChan <- gin.H{"message": "User not found"}
		} else {
			responseChan <- gin.H{"message": "Database error"}
		}
		return
	}

	// Check if sender has enough balance
	if fromUser.Balance < request.Amount {
		responseChan <- gin.H{"message": "Balance is not enough"}
		return
	}

	// Retrieve receiver (to user) from database
	var toUser models.User
	err = database.DB.Where("id = ?", request.TargetUser).First(&toUser).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			responseChan <- gin.H{"message": "Target user not found"}
		} else {
			responseChan <- gin.H{"message": "Database error"}
		}
		return
	}

	// Perform transfer operation
	previousBalanceFrom := fromUser.Balance
	fromUser.Balance -= request.Amount
	toUser.Balance += request.Amount

	// Save updated balances
	err = database.DB.Transaction(func(tx *gorm.DB) error {
		if err = tx.Save(&fromUser).Error; err != nil {
			return err
		}
		if err = tx.Save(&toUser).Error; err != nil {
			return err
		}
		// Create transfer record
		transfer = models.Transfer{
			FromUserID:    fromUser.ID,
			ToUserID:      toUser.ID,
			Amount:        request.Amount,
			Remarks:       request.Remarks,
			BalanceBefore: previousBalanceFrom,
			BalanceAfter:  fromUser.Balance,
			CreatedDate:   time.Now(),
		}
		if err = tx.Create(&transfer).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		responseChan <- gin.H{"message": "Failed to process transfer", "error": err.Error()}
		return
	}

	// Prepare response
	response := gin.H{
		"status": "SUCCESS",
		"result": gin.H{
			"transfer_id":    transfer.ID,
			"amount":         request.Amount,
			"remarks":        request.Remarks,
			"balance_before": previousBalanceFrom,
			"balance_after":  fromUser.Balance,
			"created_date":   time.Now().Format("2006-01-02 15:04:05"),
		},
	}

	responseChan <- response
}

func Transactions(c *gin.Context) {
	// Parse JWT token
	claims, err := auth.ParseJWT(c.Request.Header.Get("Authorization"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthenticated"})
		return
	}

	// Get transfers for the user
	var transfers []models.Transfer
	database.DB.Where("from_user_id = ?", claims.UserID).Find(&transfers)

	// Get payments for the user
	var payments []models.Payment
	database.DB.Where("user_id = ?", claims.UserID).Find(&payments)

	// Get top-ups for the user
	var topUps []models.TopUp
	database.DB.Where("user_id = ?", claims.UserID).Find(&topUps)

	// Prepare result array
	var result []gin.H

	// Process transfers
	for _, t := range transfers {
		entry := gin.H{
			"transfer_id":      t.ID,
			"status":           "SUCCESS",
			"user_id":          claims.UserID,
			"transaction_type": "DEBIT",
			"amount":           t.Amount,
			"remarks":          t.Remarks,
			"balance_before":   t.BalanceBefore,
			"balance_after":    t.BalanceAfter,
			"created_date":     t.CreatedDate.Format("2006-01-02 15:04:05"),
		}
		result = append(result, entry)
	}

	// Process payments
	for _, p := range payments {
		entry := gin.H{
			"payment_id":       p.ID,
			"status":           "SUCCESS",
			"user_id":          claims.UserID,
			"transaction_type": "DEBIT",
			"amount":           p.Amount,
			"remarks":          p.Remarks,
			"balance_before":   p.BalanceBefore,
			"balance_after":    p.BalanceAfter,
			"created_date":     p.CreatedDate.Format("2006-01-02 15:04:05"),
		}
		result = append(result, entry)
	}

	// Process top-ups
	for _, tu := range topUps {
		entry := gin.H{
			"top_up_id":        tu.ID,
			"status":           "SUCCESS",
			"user_id":          claims.UserID,
			"transaction_type": "CREDIT",
			"amount":           tu.Amount,
			"balance_before":   tu.BalanceBefore,
			"balance_after":    tu.BalanceAfter,
			"created_date":     tu.CreatedDate.Format("2006-01-02 15:04:05"),
		}
		result = append(result, entry)
	}

	// Return the combined results
	response := gin.H{
		"status": "SUCCESS",
		"result": result,
	}

	c.JSON(http.StatusOK, response)
}

func UpdateProfile(c *gin.Context) {
	// Parse JWT token
	claims, err := auth.ParseJWT(c.Request.Header.Get("Authorization"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthenticated"})
		return
	}

	// Retrieve current user
	var user models.User
	err = database.DB.Where("id = ?", claims.UserID).First(&user).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}

	// Bind JSON request body to struct
	var request struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Address   string `json:"address"`
	}

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}

	// Update user fields
	user.FirstName = request.FirstName
	user.LastName = request.LastName
	user.Address = request.Address
	user.UpdatedDate = time.Now()

	// Save updated user profile to database
	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update profile"})
		return
	}

	// Prepare response
	response := gin.H{
		"status": "SUCCESS",
		"result": gin.H{
			"user_id":      user.ID,
			"first_name":   user.FirstName,
			"last_name":    user.LastName,
			"address":      user.Address,
			"updated_date": user.UpdatedDate.Format("2006-01-02 15:04:05"),
		},
	}

	c.JSON(http.StatusOK, response)
}
