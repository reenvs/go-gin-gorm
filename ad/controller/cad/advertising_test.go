package cad

import (
	"iptv/common/comtest"
	"iptv/common/logger"
	"net/http"
	"testing"
)

// 提交广告创意(保存)
func TestAdvertisingSaveHandler(t *testing.T) {
	db := comtest.CreateTestDB(t, comtest.DBSourceCad)
	logger.SetLevel(logger.DEBUG)
	comtest.TestAdminScopeHandler(db, http.MethodPost, "/cad/advertising/save",
		"6fac2b13b3d3d5292569f92eea6d103d", //登陆token
		//comtest.CmsAccessKey,               //CMS accesskey
		AdvertisingSaveHandler,
		map[string]interface{}{
			"order_id":"20190622994c00000001",
			"dsp_id":"994c",
			"type":6,
			"action_type":1,
			"landing_url":"www.111.mp4",
			"package_name":"阿迪",
			"is_origin":0,
			"ad_origin_content_type":1,
			"end_time":"2019-11-11",
			"title":"阿迪达斯",
			"src":"www.111.mp4",
			"duration":100,
		},
		t)
}

// 创意code码
func TestAdvertisingCodeHandler(t *testing.T) {
	db := comtest.CreateTestDB(t, comtest.DBSourceCad)
	logger.SetLevel(logger.DEBUG)
	comtest.TestAdminScopeHandler(db, http.MethodPost, "/cad/advertising/code",
		"6fac2b13b3d3d5292569f92eea6d103d", //登陆token
		//comtest.CmsAccessKey,               //CMS accesskey
		AdvertisingCodeHandler,
		map[string]interface{}{
			"id":5,
			"code":"11111111111111111111",
		},
		t)
}

// 创意列表
func TestAdvertisingListHandler(t *testing.T) {
	db := comtest.CreateTestDB(t, comtest.DBSourceCad)
	logger.SetLevel(logger.DEBUG)
	comtest.TestAdminScopeHandler(db, http.MethodGet, "/cad/advertising/list",
		"6fac2b13b3d3d5292569f92eea6d103d", //登陆token
		//comtest.CmsAccessKey,               //CMS accesskey
		AdvertisingListHandler,
		map[string]interface{}{
			"limit":10,
		},
		t)
}

// 创意详情
func TestAdvertisingDetailHandler(t *testing.T) {
	db := comtest.CreateTestDB(t, comtest.DBSourceCad)
	logger.SetLevel(logger.DEBUG)
	comtest.TestAdminScopeHandler(db, http.MethodGet, "/cad/advertising/detail",
		"6fac2b13b3d3d5292569f92eea6d103d", //登陆token
		//comtest.CmsAccessKey,               //CMS accesskey
		AdvertisingDetailHandler,
		map[string]interface{}{
			"id":1,
		},
		t)
}

// 创意删除
func TestAdvertisingDeleteHandler(t *testing.T) {
	db := comtest.CreateTestDB(t, comtest.DBSourceCad)
	logger.SetLevel(logger.DEBUG)
	comtest.TestAdminScopeHandler(db, http.MethodPost, "/cad/advertising/delete",
		"6fac2b13b3d3d5292569f92eea6d103d", //登陆token
		//comtest.CmsAccessKey,               //CMS accesskey
		AdvertisingDeleteHandler,
		map[string]interface{}{
			"id":[]uint32{4,5},
		},
		t)
}