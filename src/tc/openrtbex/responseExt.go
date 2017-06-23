package openrtbex

type ResponseExt struct {
	ConfirmType ConfirmType `json:"confirmType,omitempty"`
	UserId      *string     `json:"userId,omitempty"`
	DeviceIp    *string     `json:"deviceIp,omitempty"`
}
