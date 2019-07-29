package constant

const (
	ContextStorageRoot    = "contextstorageroot"
	ContextDb             = "contextdb"
	ContextReadOnlyDb     = "contextreadonlydb"
	ContextStatDb         = "contextstatdb"
	ContextUser           = "contextuser"
	ContextOsType         = "contextostype"
	ContextApp            = "contextapp"
	ContextOsVersion      = "contextosversion"
	ContextAppVersion     = "contextappversion"
	ContextAppVersionObj  = "contextappversionobj"
	ContextAppKey         = "contextappkey"
	ContextInstallationId = "contextinstallationid"
	ContextAppId          = "contextappid"
	ContextVersionId      = "contextversionid"
	ContextChannel        = "contextchannel"
	ContextAdmin          = "contextadmin"
	ContextAdminId        = "contextadminid"
	ContextAdminName      = "contextadminname"
	ContextToken          = "contexttoken"
	ContextModuleAccess   = "contextmoduleaccess"
	ContextRequestBody    = "contextrequestbody"
	ContextTableName      = "contexttablename"
	ContextOperationType  = "contextoperationtype"
	ContextOperationAppId = "contextoperationappid"
	ContextOldValue       = "contextoldvalue"
	ContextNewValue       = "contextnewvalue"
	ContextError          = "contexterror"
	ContextScopeBody      = "contextscopebody"
	ContextScope          = "contextscope"
	ContextSolr           = "contextsolr"
	ContextUserGroup      = "contextusergroup"
	ContextJsonParams     = "context_jsonparams"
	ContextJsonBody       = "context_jsonbody"
	ContextParams         = "context_params"

	ContextTimestamp = "context_timestamp"
	ContextSign      = "context_sign"
	ContextDeviceID  = "context_device_id"
	ContextBSSID     = "context_bssid"
	ContextGCID      = "context_gcid"
	ContextIMEI      = "context_imei"
	ContextMac       = "context_mac"

	ContextApiVersion = "contextaapiver"

	ContextAllCRVideoIds     = "context_all_cr_video_ids"
	ContextUserInBlacklist   = "context_user_in_blacklist"
	ContextEnabledCRVideoIds = "context_enabled_cr_video_ids"

	ContextAllCRStreamIds     = "context_all_cr_stream_ids"
	ContextEnabledCRStreamIds = "context_enabled_cr_stream_ids"

	ContextAppChannelConfig = "context_app_channel_config"

	ContextClientAreaCode  = "context_client_aream_code"
	ContextClientCountry   = "context_client_country"
	ContextForceReloadUser = "context_force_reload_user"
)

const (
	ConfigEnableVideoLevelControl    = "enable_video_level_control"
	ConfigEnableCRGroupPloy          = "enable_crgroup_ploy"
	ConfigEnableLoadUserInfo         = "enable_load_user_info"
	ConfigEnableBlockBlacklistStream = "enable_block_blacklist_stream" //启用黑名单用户直接屏蔽频道信息
	ConfigEnableBlockBlacklistVideo  = "enable_block_blacklist_video"  //启用黑名单用户直接屏蔽点播视频信息
	VideoLevelMinUserLevel           = "video_level_min_user_level"
	VideoLevelMinDeviceScore         = "video_level_min_device_score"
	ConfigResourceStatusNotifyUrl    = "resource_status_notify_url"

	CmsChannelCleanCacheData = "CMS_CleanCacheData"
	ImsChannelCleanCacheData = "IMS_CleanCacheData"
	MmsChannelCleanCacheData = "MMS_CleanCacheData"
)
