package translate

type DeviceInfo struct {
	DeviceId     string `json:"deviceId"`
	Ip           string `json:"ip"`
	Brand        string `json:"brand"`
	model        string `json:"model"`
	is4K         bool   `json:"is4K"`
	osVersion    string `json:"osVersion"`
	screenWidth  int    `json:"screenWidth"`
	screenHeight int    `json:"screenHeight"`
	carrier      string `json:"carrier"`
	deviceMac    string `json:"deviceMac"`
	ssid         string `json:"ssid"`
	bssid        string `json:"bssid"`
	city         string `json:"city"`
}

type Profile struct {
	gender      string `json:"gender"`
	age         int    `json:"age"`
	childrenNum int    `json:"childrenNum"`
}

type UserInfo struct {
	uid     string   `json:"uid"`
	profile Profile  `json:"profile"`
	tags    []string `json:"tags"`
}
type AppInfo struct {
	platform     string `json:"platform"`
	platformType string `json:"platformType"`
	version      string `json:"version"`
}
type ImpRequests struct {
	impId       string   `json:"impId"`
	plotId      string   `json:"plotId"`
	dealId      string   `json:"dealId"`
	plotType    string   `json:"plotType"`
	path        string   `json:"path"`
	channel     string   `json:"channel"`
	tags        []string `json:"tags"`
	title       string   `json:"title"`
	function    string   `json:"function"`
	width       int      `json:"width"`
	height      int      `json:"height"`
	minDuration int      `json:"minDuration"`
	maxDuration int      `json:"maxDuration"`
}
type Context struct {
	ts int `json:"ts"`
}

// 风尚广告接口参数
type AdPostParams struct {
	deviceInfo  DeviceInfo    `json:"deviceInfo"`
	userInfo    UserInfo      `json:"userInfo"`
	appInfo     AppInfo       `json:"appInfo"`
	impRequests []ImpRequests `json:"impRequests"`
	context     Context       `json:"context"`
}




