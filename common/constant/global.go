package constant

const (
	MaxPageSize     = 256
	DefaultPageSize = 10

	SortAsc       = "asc"
	SortDesc      = "desc"
	DefaultColumn = "id"

	CdnSettingKey = "cdn_setting_key"

	// Kv store keys
	MergeAdSyncedAt         = "merge_ad_synced_at"
	ResourceGroupSyncedAt   = "resource_group_synced_at"
	IptvVideoGroupSyncedAt  = "iptv_video_group_synced_at"
	IptvSeriesGroupSyncedAt = "iptv_series_group_synced_at"
	PersonSyncedAt          = "person_synced_at"
	CpSyncedAt              = "cp_synced_at"
	PersonMediaSyncedAt     = "personmedia_synced_at"
	LanguageSyncedAt        = "language_synced_at"
	TranslationSyncedAt     = "translation_synced_at"
	DongFangMediaSyncedAt   = "dongfang_media_synced_at"
	ChinaCloudMediaSyncedAt = "chinacloud_media_synced_at"

	WhitelistSettingKey    = "whitelist_setting_key"
	CommentSettingKey      = "comment_setting_key"
	MmsProviderSettingKey  = "mmsprovider_setting_key"
	ContentParamSettingKey = "contentparam_setting_key"
	JumpSettingKey         = "jump_setting_key"
	FilterSettingKey       = "filter_setting_key"
	ConfigSettingKey       = "config_setting_key"
	AppConfigSettingKey    = "app_config_setting_key"
	AppProviderSettingKey  = "app_provider_setting_key"
	LibrarySettingKey      = "library_setting_key"
	RedisSettingKey        = "redis_setting_key"
	SmsSettingKey          = "sms_setting_key"
	ThirdPartyInfoKey      = "third_party_key"
	PlayStatKey            = "play_stat_key"
	ChannelCacheKey        = "channel_cache_key"
	NotificationSettingKey = "notification_setting_key"
	FunctionSettingKey     = "function_setting_key"
	AdSettingKey           = "ad_setting_key"
	PartnerSettingKey      = "partner_setting_key"
	ChargeSettingKey       = "charge_setting_key"
	YiPlusSettingKey       = "yi_plus_setting_key"     //Yi+接口相关配置
	ReportSettingKey       = "report_setting_key"      //chinatv 刷量配置
	ProviderSettingKey     = "provider_setting_key"    //real相关配置
	TranslationSettingKey  = "translation_setting_key" //多语言翻译开关配置
	TvmSettingKey          = "tvm_setting_key"         //天脉相关配置
	CdnTtlSettingKey       = "cdnttl_setting_key"      //cdn生效时间
	AdUnitTimeoutKey       = "ad_unit_timeout_key"     //广告请求超时时间

	TideSyncTime         = "tide_sync_time"
	TideIncreaseSyncTime = "tide_increase_sync_time"
	CRGroupSyncedAt      = "cr_group_synced_at"
	CRGroupPloySyncedAt  = "cr_groupploy_synced_at"

	TmpStorage      = "tmp"
	VideoTmpStorage = "video_tmp"
	LogStorage      = "log"

	DefaultLanguageCode = "zh"

	IptvShenZhenCategory = "深圳IPTV新闻630"

	InternalTypeAll   = 0 //全部
	InternalTypeTrue  = 1 //内部商品
	InternalTypeFalse = 2 //非内部商品

	ActiveTypeAll   = 0 //全部
	ActiveTypeTrue  = 1 //已上架
	ActiveTypeFalse = 2 //已下架

	DisabledTypeAll   = 0 //全部
	DisabledTypeTrue  = 1 //禁用
	DisabledTypeFalse = 2 //可用

	OnlineTypeAll   = 0 //全部
	OnlineTypeTrue  = 1 //已上线
	OnlineTypeFalse = 2 //未上线

	DownloadableTypeAll   = 0 //全部
	DownloadableTypeTrue  = 1 //可下载
	DownloadableTypeFalse = 2 //不可下载

	SystemTypeAll   = 0 //全部分组
	SystemTypeTrue  = 1 //系统分组
	SystemTypeFalse = 2 //非系统分组

	UserGroupTypeAdd    = 1 // 用户组记录增加
	UserGroupTypeDelete = 2 // 用户组记录删除

	RedeemCodeStatusAll     = 0 //全部
	RedeemCodeStatusUsed    = 1 //已使用
	RedeemCodeStatusNotUsed = 2 //未使用

	OsTypeAll = 0 //全部系统类型

	UserCacheTTL     = 300     // 用户结构缓存时间，单位秒
	CommentCacheTTL  = 300     // 评论缓存时间
	AreaFeedCacheTTL = 60 * 30 // 地域消息流缓存时间

	SearchTypeUnknown    = 0 //未知类型
	SearchTypeProduct    = 1 //商品搜索, search_value为商品id
	SearchTypePackage    = 2 //套餐搜索, search_value为套餐id
	SearchTypeVipProduct = 3 //vip商品, search_value留空
	SearchTypeVipPackage = 4 //vip套餐, search_value留空
	SearchTypePrivilege  = 5 //特权商品搜索, search_value为商品的特权

	ReachedTypeAll   = 0
	ReachedTypeTrue  = 1
	ReachedTypeFalse = 2

	ExpiredTypeAll   = 0
	ExpiredTypeTrue  = 1
	ExpiredTypeFalse = 2

	AccessKey       = "AccessKey"
	ModuleSignature = "Signature"
	ModuleSalt      = "module_salt"
	ModuleSaltV2    = "2833D0F2E9E9739AF92FDA1CAFE4505CF4C654EF"

	PropertyArtist    = "艺术家"
	PropertyWriter    = "作者"
	PropetyCategoryId = 5

	SmsProviderYunPian    = 1
	SmsProviderYunTongXin = 2

	/*---------- global app_id definition ----------*/
	MobileTVAppId        = 1
	ChinaTVAppId         = 11
	ShenZhenIptv630AppId = 4
	YouShiTVAppId        = 111
	HebeiAppId           = 3 // 河北IPTV app id

	YOUSHITV_WOPAY_KEY = "3d555dd9267b23059b590297252d6f5b"

	CRPloyDisabelControlID     = 1
	CRPloyMaxInternalControlID = 10000

	EpgBlockMinStartDeviation = 5 * 60
	EpgBlockStartDeviation    = 30 * 60 //30 分钟，EpgBlock 时间误差

	TimeBlockDeviation = 5 * 60 //30 分钟，EpgBlock 时间误差

	DeviceScoreSuspectBlacklistMin      = 5
	BSSIDSuspectBlacklistCountMin       = 5
	DeviceResetSuspectBlacklistCountMin = 6 //疑似重置
	DeviceResetCountMin                 = 3 // 基本确认重置

	UserInfoCacheLoadInterval    = 5 * 60
	InstallationIPsCacheInterval = 20

	BlacklistLoadInterval = 3 * 60

	//开启广告的最小安装量
	EnableAdMinInstallationCount = 50

	ParamsDecodeKey = "CiZa4fDec10jzgHC"
	ParamsKeyPrefix = "c042286720dw3c5k"
	ParamsAesIv     = "q@7m*d3e4pk5c3wd"
)
