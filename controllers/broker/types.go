package broker

import "encoding/json"

type stateChangeRequest struct {
	NewState          int     `json:"new_state,omitempty"`
	ThresholdKey      *string `json:"threshold_key,omitempty"`
	PlainThresholdKey *string `json:"plain_threshold_key,omitempty"`
}

type resolveConflictRequest struct {
	WinnerUsername string `json:"winner_username,omitempty"`
	ThresholdKey   string `json:"threshold_key,omitempty"`
}

type transactionRequest struct {
	To string `json:"to,omitempty"`
}

type transactionResponse struct {
	Transaction json.RawMessage `json:"transaction,omitempty"`
	NetworkID   int64           `json:"network_id,omitempty"`
}

type rawTransactionRequest struct {
	Data string `json:"transaction"`
}

type transactionProcessingResponse struct {
	Hash string `json:"hash,omitempty"`
}
