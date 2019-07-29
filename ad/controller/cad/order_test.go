package cad

import (
	"iptv/common/comtest"
	"iptv/common/logger"
	"net/http"
	"testing"
)

// 保存、发布 订单 （下单）
func TestOrderSaveHandler(t *testing.T) {
	db := comtest.CreateTestDB(t, comtest.DBSourceCad)
	logger.SetLevel(logger.DEBUG)
	comtest.TestAdminScopeHandler(db, http.MethodPost, "/cad/order/save",
		"6fac2b13b3d3d5292569f92eea6d103d", //登陆token
		//comtest.CmsAccessKey,               //CMS accesskey
		AdOrderSaveHandler,
		map[string]interface{}{
			"area_code":"1156440000",
			"order_name":"222",
			"total_count":10000,
			"dsp_id":"4d14",
			"dsp_order_id":"123",
			"ad_position_id":"开始",
			"order_desc":"新款首发",
			"begin_date":"2019-06-30",
			"end_date":"2019-07-30",
			"day_control_type":1,
		},
		t)
}

// 订单列表
func TestOrderListHandler(t *testing.T) {
	db := comtest.CreateTestDB(t, comtest.DBSourceCad)
	logger.SetLevel(logger.DEBUG)
	comtest.TestAdminScopeHandler(db, http.MethodGet, "/cad/order/list",
		"6fac2b13b3d3d5292569f92eea6d103d", //登陆token
		//comtest.CmsAccessKey,               //CMS accesskey
		AdOrderListHandler,
		map[string]interface{}{
			"limit":10,
		},
		t)
}

// 订单详情
func TestOrderDetailHandler(t *testing.T) {
	db := comtest.CreateTestDB(t, comtest.DBSourceCad)
	logger.SetLevel(logger.DEBUG)
	comtest.TestAdminScopeHandler(db, http.MethodGet, "/cad/order/detail",
		"6fac2b13b3d3d5292569f92eea6d103d", //登陆token
		//comtest.CmsAccessKey,               //CMS accesskey
		AdOrderDetailHandler,
		map[string]interface{}{
			"id":1,
		},
		t)
}

// 订单状态修改
func TestOrderStatusHandler(t *testing.T) {
	db := comtest.CreateTestDB(t, comtest.DBSourceCad)
	logger.SetLevel(logger.DEBUG)
	comtest.TestAdminScopeHandler(db, http.MethodPost, "/cad/order/status",
		"6fac2b13b3d3d5292569f92eea6d103d", //登陆token
		//comtest.CmsAccessKey,               //CMS accesskey
		AdOrderStatusHandler,
		map[string]interface{}{
			"id":4,
			"status":4,
		},
		t)
}

// 订单上下架
func TestOrderDisableHandler(t *testing.T) {
	db := comtest.CreateTestDB(t, comtest.DBSourceCad)
	logger.SetLevel(logger.DEBUG)
	comtest.TestAdminScopeHandler(db, http.MethodPost, "/cad/order/disable",
		"6fac2b13b3d3d5292569f92eea6d103d", //登陆token
		//comtest.CmsAccessKey,               //CMS accesskey
		AdOrderDisableHandler,
		map[string]interface{}{
			"id":4,
			"disable":3,
		},
		t)
}

// 删除订单
func TestOrderDeleteHandler(t *testing.T) {
	db := comtest.CreateTestDB(t, comtest.DBSourceCad)
	logger.SetLevel(logger.DEBUG)
	comtest.TestAdminScopeHandler(db, http.MethodPost, "/cad/order/delete",
		"6fac2b13b3d3d5292569f92eea6d103d", //登陆token
		//comtest.CmsAccessKey,               //CMS accesskey
		AdOrderDeleteHandler,
		map[string]interface{}{
			"id":[]uint32{4},
		},
		t)
}