package models
 
type User struct {
	AccountID int    `json:"account_id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
}