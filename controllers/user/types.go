package user

import "api.ethscrow/models"

type poolResponse struct {
	Active    []models.Pool `json:"active"`
	Inbox     []models.Pool `json:"inbox"`
	Sent      []models.Pool `json:"sent"`
	Resolve   []models.Pool `json:"resolve"`
	Completed []models.Pool `json:"completed"`
}
