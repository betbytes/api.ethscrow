package broker

type stateChangeRequest struct {
	NewState          int     `json:"new_state,omitempty"`
	ThresholdKey      *string `json:"threshold_key,omitempty"`
	PlainThresholdKey *string `json:"plain_threshold_key,omitempty"`
}

type resolveConflictRequest struct {
	WinnerUsername string `json:"winner_username,omitempty"`
	ThresholdKey   string `json:"threshold_key,omitempty"`
}
