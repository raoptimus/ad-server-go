package openrtbex

type BidStatusExt struct {
	Status        ConfirmStatus `json:"status"`
	ClearingPrice *float32      `json:"clearingPrice,emitempty"`
	WinningPrice  *float32      `json:"winningPrice,emitempty"`
	LossReason    LossReason    `json:"lossReason,emitempty"`
	Error         string        `json:"error,emitempty"`
}
