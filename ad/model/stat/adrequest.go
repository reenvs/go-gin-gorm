package stat

type AdRequest struct {
	Id             uint32 `json:"id"`
	AdId           string `json:"ad_id"`
	AdName         string `json:"ad_name"`
	SpotCode       string `json:"spot_code"`
	SportName      string `json:"sport_name"`
	OrderId        string `json:"order_id"`
	DspId          string `json:"dsp_id"`
	TimeSlice      string `json:"time_slice"`
	RequestAll     uint32 `json:"request_all"`
	RequestSuccess uint32 `json:"request_success"`
	RequestFailed  uint32 `json:"request_failed"`
	Start          uint32 `json:"start"`
	FirstQuarter   uint32 `json:"first_quarter"`
	Half           uint32 `json:"half"`
	ThirdQuarter   uint32 `json:"third_quarter"`
	End            uint32 `json:"end"`
}
