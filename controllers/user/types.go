package user

import "api.ethscrow/models"

type poolResponse struct {
	Active []models.Pool `json:"active,omitempty"`
	Inbox  []models.Pool `json:"inbox,omitempty"`
	Sent   []models.Pool `json:"sent,omitempty"`
}
