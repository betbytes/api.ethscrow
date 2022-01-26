package broker

type stateChangeRequest struct {
	NewState     int     `json:"new_state,omitempty"`
	ThresholdKey *string `json:"threshold_key,omitempty"`
	Conflict     *bool   `json:"conflict,omitempty"`
}

type resolveConflictRequest struct {
	WinnerUsername string `json:"winner_username,omitempty"`
	ThresholdKey   string `json:"threshold_key,omitempty"`
}
