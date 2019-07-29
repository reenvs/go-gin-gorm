package model

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Order struct {
	ID             uint32      `gorm:"primary_key" json:"id"`
	AreaCode       string      `gorm:"size:64" json:"area_code"`      // 投放地域，多个逗号隔开
	OrderName      string      `gorm:"size:64" json:"order_name"`     // 投放名称，对应广告活动名称
	TotalCount     uint32      `json:"total_count"`                   // 投放总量
	DspId          string      `gorm:"size:64" json:"dsp_id"`         // 代理商编号
	DspName        string      `gorm:"size:64" json:"dsp_name"`       // 代理商名称
	Sign           string      `gorm:"-" json:"-"`                    // 签名结果
	TimeStamp      uint32      `gorm:"-" json:"time_stamp"`           // 时间戳
	DspOrderId     string      `gorm:"size:64" json:"dsp_order_id"`   // 由dsp生成的订单
	AdPositionId   string      `gorm:"size:64" json:"ad_position_id"` // 该需求需要投放的广告位id
	OrderDesc      string      `gorm:"size:64" json:"order_desc"`     // 说明
	BeginDate      string      `gorm:"size:64" json:"begin_date"`     // 订单开始投放时间 yyyy-mm-dd
	EndDate        string      `gorm:"size:64" json:"end_date"`       // 订单结束投放时间 yyyy-mm-dd
	BeginTime      string      `gorm:"size:16" json:"begin_time"`     // 开始时间段，格式 8:00
	EndTime        string      `gorm:"size:16" json:"end_time"`       // 结束时间段，格式 8:00
	UnitPrice      string      `json:"unit_price"`                    // 单价
	DayControlType uint32      `json:"day_control_type"`              // 单日数量控制类型：1 单日平均，默认；2 自定义单日数量
	DayCounts      string      `json:"day_counts"`                    // 如果单日平均该字段不传，不平均需传入每天的量
	Contact        string      `json:"contact"`                       // 联系人
	Mobile         string      `json:"mobile"`                        // 联系电话
	Email          string      `json:"email"`                         // 企业邮箱
	OrderID        string      `gorm:"size:64" json:"order_id"`       // 我方生成的订单id
	Status         uint32      `json:"status"`                        // 订单状态
	Disable        uint32      `json:"disable"`                       // 1 上架（默认）, 2 下架
	Ads            []*Material `json:"ads"`
	CreatedAt      time.Time   `json:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at"`
}

func (Order) TableName() string {
	return "order"
}

func initOrder(db *gorm.DB) error {
	var err error
	if db.HasTable(&Order{}) {
		err = db.AutoMigrate(&Order{}).Error
	} else {
		err = db.CreateTable(&Order{}).Error
	}
	return err
}

func dropOrder(db *gorm.DB) {
	db.DropTableIfExists(&Order{})
}

// 订单状态
const (
	AdStatusSubmiting = 1 // 待提交创意
	AdStatusReviewing = 2 // 审核中
	AdStatusReady     = 3 // 就绪（审核通过）
	AdStatusRunning   = 4 // 投放
	AdStatusOver      = 5 // 结束

	AdDisableAll      = 0 // 全部
	AdDisableOnShelf  = 1 // 上架
	AdDisableOffShelf = 2 // 下架
)

/*
	广告投放相关
*/
type TopBoxRequest struct {
	UID          string `json:"uId"`          // 用户业务账号
	AreaInfo     string `json:"areaInfo"`     // 区域信息
	AreaName     string `json:"areaName"`     // 区域名称
	UGroupID     string `json:"ugroupId"`     // 用户分组信息
	ManuFacturer string `json:"manufacturer"` // 机顶盒厂商
	Moudel       string `json:"moudel"`       // 机顶盒型号
	Version      string `json:"version"`      // 系统版本
	ADId         string `json:"adId"`         // 广告id
	SpotCode     string `json:"spotCode"`     // 广告位id
	SpotName     string `json:"spotName"`     // 广告位名称
	ResourceID   string `json:"resourceId"`   // 资源id
	//ResourceName     string `json:"resourceName"`     // 资源名称
	//ResourceType     string `json:"resourceType"`     // 资源类型
	//ResourceTypeName string `json:"resourceTypeName"` // 资源类型名称
	Mac        string `json:"mac"`        // 机顶盒MAC值
	SN         string `json:"sn"`         // 机顶盒SN值
	CommitTime string `json:"commitTime"` // 提交时间戳
	IP         string `json:"ip"`         // 用户ip
}

type TopBoxReturn struct {
	Code     string   `json:"code"`
	Monitors Monitors `json:"monitors"`
}

// dsp请求体
type DspRequest struct {
	DeviceInfo DeviceInfo `json:"deviceInfo"` // 设备相关信息
	AppInfo    AppInfo    `json:"appInfo"`
	ImpRequest []ImpReq   `json:"impRequest"`
	Context    Context    `json:"context"`
}

// 设备相关信息
type DeviceInfo struct {
	DeviceID     string `json:"deviceId"`      // 设备编码，mac
	IP           string `json:"ip"`            // ip地址
	Brand        string `json:"brand"`         // 设备品牌
	Model        string `json:"model"`         // 设备型号
	Is4K         bool   `json:"is4k"`          // 是否高清
	OsVersion    string `json:"os_version"`    // 电视操作系统的版本
	ScreenWidth  uint32 `json:"screen_width"`  // 屏幕宽度
	ScreenHeight uint32 `json:"screen_height"` // 屏幕高度
	Carrier      string `json:"carrier"`       // 运营商版本
	DeviceMac    string `json:"device_mac"`    // 设备的mac地址
	Ssid         string `json:"ssid"`          // 无线网络唯一id
	Bssid        string `json:"bssid"`         // 无线路由mac地址
	City         string `json:"city"`          // 广协标准城市码
}

// 用户信息用于进行大数据只能推送广告
type UserInfo struct {
	UID     string `json:"uid"` // 用户账号
	Profile struct {
		// 用户数据，媒体方存入用户表，广告请求时查询
		Gender      string `json:"gender"`
		Age         uint32 `json:"age"`
		ChildrenNum uint32 `json:"children_num"`
	}
	Tags []string `json:"tags"` // 媒体 已有的用户标签数据
}

// 媒体相关信息
type AppInfo struct {
	Platform     string `json:"platform"`     // 媒体情况
	PlatformType string `json:"platformType"` // 客户端类型 : 1 PC,2 mobile,3 OTT
	Version      string `json:"version"`      // 客户端版本
}

// 广告请求
type ImpReq struct {
	ImpID       string   `json:"impId"`       // 广告曝光id
	PlotID      string   `json:"plotId"`      // 广告位编号,风尚提供统一编码
	DealID      string   `json:"dealId"`      // 广告订单id
	PlotType    string   `json:"plotType"`    // 广告类型
	Path        string   `json:"path"`        // 广告位的路径
	Channel     string   `json:"channel"`     // 广告位所在节目的类型
	Tags        []string `json:"tags"`        // 广告位在节目标签
	Title       string   `json:"title"`       // 广告位所在节目名称
	Function    string   `json:"function"`    // 触发广告曝光的动作
	Width       uint32   `json:"width"`       // 广告位的宽度
	Height      uint32   `json:"height"`      // 广告位的高度
	MinDuration uint32   `json:"minDuration"` // 广告的最短时长
	MaxDuration uint32   `json:"maxDuration"` // 广告的最长时长
}

type Context struct {
	TS uint32 `json:"ts"` // 广告请求时间
}

// dsp响应结果
type DspResponse struct {
	Status uint32 `json:"status"`
	Cost   uint32 `json:"cost"`  // 耗时
	ReqID  string `json:"reqId"` // 请求id，PDB 生成的唯一 id
	Ads    Ads    `json:"ads"`
}

type Ads struct {
	AdID       string   `json:"adId"`       // 广告活动id
	ImpId      string   `json:"impId"`      //  同impRequest.impId
	DealId     string   `json:"dealId"`     //  订单id
	AdPlotId   string   `json:"adPlotId"`   // 广告位 id
	PlotType   string   `json:"plotType"`   // 广告位类型
	CreativeId string   `json:"creativeId"` // 风尚创意 id
	MediaType  string   `json:"mediaType"`  // 创意类型
	MediaId    string   `json:"mediaId"`    // 创意在媒体端的 id
	MediaUrl   string   `json:"mediaUrl"`   // 创意在媒体端具体的源文件
	Width      uint32   `json:"width"`      // 广告位宽度
	Height     uint32   `json:"height"`     // 广告高度
	Duration   uint32   `json:"duration"`   // 广告时长,单位为秒
	StartDelay uint32   `json:"startDelay"` // 开始播放的时间
	Price      uint32   `json:"price"`      // 竞价
	Monitors   Monitors `json:"monitors"`   // 监控
}

type Monitors struct {
	Click        []string `json:"click"`        // 广告点击监控， 监控支持多个 todo
	Impress      []string `json:"impress"`      // 广告曝光监控
	Start        []string `json:"start"`        // 开始播放的监控
	FirstQuarter []string `json:"firstQuarter"` // 播放到 1/4 的监控
	ThirdQuarter []string `json:"thirdQuarter"` // 播放到 1/2 的监控
	Midpoint     []string `json:"midpoint"`     // 播放到 3/4 的监控
	End          []string `json:"end"`          // 播放结束的监控
}

/*
	广告报告相关
*/
type AdReport struct {
	TopBoxRequest
}
