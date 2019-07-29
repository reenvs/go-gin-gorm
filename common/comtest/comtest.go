package comtest

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	jsoniter "github.com/json-iterator/go"
	"io"
	"iptv/common/constant"
	"iptv/common/logger"
	"iptv/common/middleware"
	"iptv/common/util"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

const (
	DBSourceAms = "root:000000@(127.0.0.1:3306)/ams?charset=utf8&parseTime=True&loc=Local"
	DBSourceCad = "root:000000@(127.0.0.1:3306)/cad?charset=utf8&parseTime=True&loc=Local"
)

func CreateTestDB(t *testing.T, dbSource string) *gorm.DB {

	logger.Debugf("open %v :%v", "mysql", dbSource)
	db, err := gorm.Open("mysql", dbSource)
	if err != nil {
		t.Fatal("Open db Failed!!!!", err)
		return nil
	}
	db.LogMode(true)
	db.DB().SetMaxIdleConns(5)
	db.DB().SetMaxOpenConns(15)
	db.DB().SetConnMaxLifetime(time.Second * 600)
	return db
}

func initTest() {
	logger.SetLevel(logger.DEBUG)
}

//TestApiHandler  测试客户端api
func TestApiHandler(db *gorm.DB, method string, uriPath string,
	handler func(*gin.Context),
	moduleHandler gin.HandlerFunc,
	moduleHandler2 gin.HandlerFunc,
	signatureMiddleware gin.HandlerFunc,
	params map[string]interface{}, t *testing.T) map[string]interface{} {

	dbMiddleware := middleware.GetDbHandler(db, constant.ContextDb)
	//signatureMiddleware := middleware.SignatureVerifyHandler(false)

	pa := strings.Split(uriPath, "/")
	groupName := pa[1]
	r := gin.New()
	baseAPI := r.Group(groupName)
	var apiGroup *gin.RouterGroup
	var apiVer string
	apiVer = pa[2]
	baseAPI.Use(middleware.CIBNAPiVersion())
	apiGroup = baseAPI

	apiv10 := apiGroup.Group(apiVer)
	apiv10.Use(dbMiddleware)

	if moduleHandler2 != nil {
		apiv10.Use(moduleHandler)
	}
	apiv10.Use(middleware.BaseParamentersVerifyHandler)

	if signatureMiddleware != nil {
		apiv10.Use(signatureMiddleware)
	}

	if moduleHandler2 != nil {
		apiv10.Use(moduleHandler2)
	}

	switch method {
	case http.MethodPost:
		apiv10.POST(uriPath[len(fmt.Sprintf("/%s/%s", groupName, apiVer)):], handler)

	case http.MethodGet:
		apiv10.GET(uriPath[len(fmt.Sprintf("/%s/%s", groupName, apiVer)):], handler)
	}

	middleware.SignParams(params)

	jsonInBytes, err := jsoniter.Marshal(params)
	if err != nil {
		t.Errorf("%v %v %v", method, uriPath, err)
		return nil
	}

	var reqBody io.Reader

	switch method {
	case http.MethodPost:
		reqBody = bytes.NewReader(jsonInBytes)

	case http.MethodGet:
		getParams := make(url.Values)
		for k, v := range params {
			switch t := v.(type) {
			case []uint32:
				for _, vv := range t {
					getParams.Add(k, fmt.Sprintf("%v", vv))
				}
			default:
				getParams.Set(k, fmt.Sprintf("%v", v))
			}
		}
		uriPath = fmt.Sprintf("%s?%s", uriPath, getParams.Encode())
		reqBody = nil
	}

	logger.Debugf("%v:%v", method, string(jsonInBytes))

	req, _ := http.NewRequest(method, uriPath, reqBody)

	req.Header.Set("accept", "application/json")
	if method == http.MethodPost {
		req.Header.Set("content-type", "application/json; charset=UTF-8")
	}

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	logger.Debugf("RES:%v  %v", uriPath, w.Body.String())
	resObj := make(map[string]interface{})
	err = jsoniter.Unmarshal(w.Body.Bytes(), &resObj)
	if err != nil {
		t.Error(err)
		return nil
	}

	logger.Debugf("RES JSON:%v  %v", uriPath, util.Obj2JsonIndent(resObj))
	return resObj
}

//TestAdminScopeHandler 用于测试后台管理接口，并且接口需要校验是否有APP管理权限
func TestAdminScopeHandler(db *gorm.DB, method, uriPath string,
	//loginToken, amsAccessKey string,
	loginToken string,
	handler func(*gin.Context),
	//moduleHandler gin.HandlerFunc,
	params map[string]interface{}, t *testing.T) map[string]interface{} {
	dbMiddleware := middleware.GetDbHandler(db, constant.ContextDb)
	roDbMiddleware := middleware.GetDbHandler(db, constant.ContextReadOnlyDb)
	adminMiddleware := middleware.AdminVerifyHandler("http://127.0.0.1:26880")
	//logMiddleware := middleware.OperationLogHandler("http://127.0.0.1:26880", amsAccessKey)

	pa := strings.Split(uriPath, "/")
	groupName := pa[1]

	r := gin.New()
	baseApi := r.Group(groupName)
	//baseApi.Use(dbMiddleware, roDbMiddleware, logMiddleware, adminMiddleware)
	baseApi.Use(dbMiddleware, roDbMiddleware, adminMiddleware)
	//baseApi.Use(moduleHandler, middleware.CorsAllowHandler)
	baseApi.Use(middleware.CorsAllowHandler)
	switch method {
	case http.MethodPost:
		baseApi.POST(uriPath[len(groupName)+1:], handler)

	case http.MethodGet:
		baseApi.GET(uriPath[len(groupName)+1:], handler)
	}

	middleware.SignParams(params)

	jsonInBytes, err := jsoniter.Marshal(params)
	if err != nil {
		t.Errorf("%v %v %v", method, uriPath, err)
		return nil
	}

	var reqBody io.Reader

	switch method {
	case http.MethodPost:
		reqBody = bytes.NewReader(jsonInBytes)

	case http.MethodGet:
		getParams := make(url.Values)
		for k, v := range params {
			getParams.Set(k, fmt.Sprintf("%v", v))
		}
		uriPath = fmt.Sprintf("%s?%s", uriPath, getParams.Encode())
		reqBody = nil
	}

	logger.Debugf("%v:%v", method, string(jsonInBytes))

	req, _ := http.NewRequest(method, uriPath, reqBody)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", loginToken)
	req.Header.Set("Cookie", fmt.Sprintf("token=%v", loginToken))
	if method == http.MethodPost {
		req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	}

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	logger.Debugf("RES:%v  \n%v", uriPath, w.Body.String())
	resObj := make(map[string]interface{})
	err = jsoniter.Unmarshal(w.Body.Bytes(), &resObj)
	if err != nil {
		t.Error(err)
		return nil
	}
	logger.Debugf("RES JSON:%v  \n%v", uriPath, util.Obj2JsonIndent(resObj))
	return resObj
}