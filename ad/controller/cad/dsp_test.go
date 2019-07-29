package cad

import (
	"iptv/common/comtest"
	"iptv/common/logger"
	"net/http"
	"testing"
)

// 保存、发布 广告商
func TestDspSaveHandler(t *testing.T) {
	db := comtest.CreateTestDB(t, comtest.DBSourceCad)
	logger.SetLevel(logger.DEBUG)
	comtest.TestAdminScopeHandler(db, http.MethodPost, "/cad/dsp/save",
		"6fac2b13b3d3d5292569f92eea6d103d", //登陆token
		//comtest.CmsAccessKey,               //CMS accesskey
		DspSaveHandler,
		map[string]interface{}{
			"name":"中兴",
		},
		t)
}

// Dsp列表
func TestDspListHandler(t *testing.T) {
	db := comtest.CreateTestDB(t, comtest.DBSourceCad)
	logger.SetLevel(logger.DEBUG)
	comtest.TestAdminScopeHandler(db, http.MethodGet, "/cad/dsp/list",
		"6fac2b13b3d3d5292569f92eea6d103d", //登陆token
		//comtest.CmsAccessKey,               //CMS accesskey
		DspListHandler,
		map[string]interface{}{
			"limit":10,
		},
		t)
}

// Dsp详情
func TestDspDetailHandler(t *testing.T) {
	db := comtest.CreateTestDB(t, comtest.DBSourceCad)
	logger.SetLevel(logger.DEBUG)
	comtest.TestAdminScopeHandler(db, http.MethodGet, "/cad/dsp/detail",
		"6fac2b13b3d3d5292569f92eea6d103d", //登陆token
		//comtest.CmsAccessKey,               //CMS accesskey
		DspDetailHandler,
		map[string]interface{}{
			"id":1,
		},
		t)
}

// Dsp删除
func TestDspDeleteHandler(t *testing.T) {
	db := comtest.CreateTestDB(t, comtest.DBSourceCad)
	logger.SetLevel(logger.DEBUG)
	comtest.TestAdminScopeHandler(db, http.MethodPost, "/cad/dsp/delete",
		"6fac2b13b3d3d5292569f92eea6d103d", //登陆token
		//comtest.CmsAccessKey,               //CMS accesskey
		DspDeleteHandler,
		map[string]interface{}{
			"id":[]uint32{12},
		},
		t)
}