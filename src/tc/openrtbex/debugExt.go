package openrtbex

type DebugExt struct {
	Error               string `json:"error,omitempty"`
	HttpStatusCode      int    `json:"status,omitempty"`
	HttpRequestContent  string `json:"request,omitempty"`
	HttpResponseContent string `json:"response,omitempty"`
	Subscriber          string `json:"subscriber,omitempty"`
	CustomData          string `json:"customdata,omitempty"`
}
