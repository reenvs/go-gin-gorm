package middleware

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"iptv/common/constant"
	"iptv/common/logger"
	"iptv/common/util"
	"math"
	"net/http"
	"net/url"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func ParamentersPrepareHandler(c *gin.Context) {
	defer func() {
		c.Next()
	}()
	ParseParam(c)
}

func BaseParamentersVerifyHandler(c *gin.Context) {
	//timeStart := time.Now()
	defer func() {
		//logger.Errorf("GIN %s BaseParamentersVerifyHandler:%v", c.Request.RequestURI, time.Now().Sub(timeStart))
		c.Next()
	}()

	params, err := ParseParam(c)
	if err != nil {
		logger.Error(err)
		return
	}

	clientIP := c.ClientIP()
	if clientIP == "127.0.0.1" || clientIP == "localhost" {
		return
	}

	paramCheck := func(missList []string, pn string) []string {
		val, ok := params[pn]
		if !ok {
			return append(missList, pn)
		} else {
			c.Set(fmt.Sprintf("context_%s", pn), val)
		}
		return missList
	}

	missParams := paramCheck(nil, "timestamp")
	missParams = paramCheck(missParams, "sign")
	missParams = paramCheck(missParams, "device_id")
	missParams = paramCheck(missParams, "bssid")
	missParams = paramCheck(missParams, "gcid")

	if len(missParams) > 0 {
		logger.Debugf("MISS Param:%v not found!%v %v %v %v", missParams,
			c.ClientIP(), c.Request.RemoteAddr, c.Request.Method, c.Request.RequestURI)
	}
}

func ModuleSignatureVerifyHandler(c *gin.Context) {
	errCode := constant.Failure
	var timestamp int64
	var sign string
	params := map[string]interface{}{}

	defer func() {
		if errCode != constant.Success {
			logger.Warn("[ModuleSignatureVerifyFailed] URL: ", c.Request.URL.String(), " | ", timestamp, " | ", sign, " | ", params)
			c.AbortWithStatusJSON(http.StatusInternalServerError, util.Err2JsonObj(
				fmt.Errorf("[ModuleSignatureVerifyFailed] URL: %v | %v | %v | %v ", c.Request.URL.String(), timestamp, sign, params)))
			return
		}
		c.Next()
	}()

	// 读取json参数, 获取timestamp以及表单
	params, err := ParseParam(c)
	if err != nil {
		logger.Error(err)
		return
	}

	//logger.Debug("verify module signature timestamp：", params["timestamp"])

	if t, ok := params["timestamp"]; ok {
		if ts, k := t.(float64); k {
			timestamp = int64(ts)
		} else if ts, k := t.(string); k {
			timestamp, _ = strconv.ParseInt(ts, 10, 64)
		}
	}
	//logger.Debug(c.Request.Header)
	sign = c.Request.Header.Get(constant.ModuleSignature)
	accessKey := c.Request.Header.Get(constant.AccessKey)
	token := c.Request.Header.Get("Authorization")
	if len(token) == 0 {
		tokenc, _ := c.Request.Cookie("token")
		if tokenc != nil {
			token = tokenc.Value
		}
	}

	//模块之间远程调用
	if len(sign) > 0 || len(accessKey) > 0 {
		//比较timestamp是否过期
		timeNow := time.Now().Unix()
		if int(math.Abs(float64(timeNow-timestamp))) > int(10*60) {
			logger.Errorf("Signature verify failed! Timestamp is expired:timestamp=%v timeNow=%v timestamp=%v %v %v",
				params["timestamp"], timeNow, timestamp, timeNow-timestamp, params)
			errCode = constant.Failure
			return
		}

		// 根据参数（除了sign参数）生成签名字符串
		signCalc, data1 := MakeModuleSignature(params, constant.ModuleSalt, timestamp)
		signCalcV2, data2 := MakeSignatureWithKey(params, constant.ModuleSaltV2, timestamp)

		// 对比sign参数
		if sign == "" || (strings.ToLower(sign) != strings.ToLower(signCalc) &&
			strings.ToLower(sign) != strings.ToLower(signCalcV2)) {
			errCode = constant.Failure
			logger.Errorf("Module Signature verify failed! sign：%v\nsign1: %v %v\nsign2: %v %v\n %v",
				sign, signCalc, data1, signCalcV2, data2, c.Request.URL.String())
			return
		} else {
			errCode = constant.Success
		}
	} else if len(token) > 0 {
		//前端页面身份验证，adminMiddleware 已经验证过身份
		errCode = constant.Success
		return
	}

	//logger.Debug("Module Signature verify successed! ", signCalc, " | ", c.Request.URL.String())
}

// timestamp + params with key + salt 进行sha1加密后hex编码
func MakeSignatureWithKey(params map[string]interface{}, salt string, timestamp int64) (string, string) {
	source := getSignatureData(params)

	// 计算出密钥并加密
	hash := sha1.New()
	hash.Write([]byte(strconv.FormatInt(timestamp, 10)))
	hash.Write([]byte(source))
	hash.Write([]byte(salt))

	return fmt.Sprintf("%X", hash.Sum(nil)), fmt.Sprintf("%v%v%v", timestamp, source, salt)
}

// timestamp + params + salt 进行sha1加密后hex编码
func MakeModuleSignature(params map[string]interface{}, salt string, timestamp int64) (string, string) {
	p := strconv.FormatInt(timestamp, 10)
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		if v, ok := params[k]; ok {
			switch t := v.(type) {
			case int, int32, int64, uint, uint32, uint64:
				p += url.QueryEscape(fmt.Sprintf("%d", t))
			case float32, float64:
				//p += url.QueryEscape(fmt.Sprintf("%.0f", t))
				p += url.QueryEscape(strconv.FormatFloat(v.(float64), 'f', -1, 64))
			case string:
				p += url.QueryEscape(fmt.Sprintf("%s", t))
			case bool:
				p += url.QueryEscape(fmt.Sprintf("%t", t))
			}
		}
	}
	p += salt
	//logger.Debug("make module sign str：", p)
	hash := sha1.New()
	//暂时不修正，避免其他模块没有升级导致校验失败，新的功能使用新的签名函数MakeSignatureWithKey
	//hash.Write([]byte(p))
	//return hex.EncodeToString(hash.Sum(nil))
	return hex.EncodeToString(hash.Sum([]byte(p))), p
}

func SignParams(params map[string]interface{}) map[string]interface{} {
	timeNow := time.Now().Unix()
	t, ok := params["timestamp"]
	if ok {
		timeNow = t.(int64)
	} else {
		params["timestamp"] = timeNow
	}
	sign, _, _ := MakeSignature(params, timeNow)
	params["sign"] = sign

	return params
}

func SignParamsWithSalt(params map[string]interface{}, salt string) {
	timeNow := time.Now().Unix()
	t, ok := params["timestamp"]
	if ok {
		timeNow = t.(int64)
	} else {
		params["timestamp"] = timeNow
	}
	sign, _ := MakeSignatureWithKey(params, salt, timeNow)
	params["sign"] = sign
}

func getSignatureData(params map[string]interface{}) string {
	var p []string
	for k, v := range params {
		if k == "sign" {
			continue
		}

		switch t := v.(type) {
		case float32, float64:
			// 2017-11-10：float64转成string类型去解析
			p = append(p, k+"="+url.QueryEscape(strconv.FormatFloat(v.(float64), 'f', -1, 64)))
			// p = append(p, k+"="+url.QueryEscape(fmt.Sprintf("%.0f", v)))
		case string:
			p = append(p, k+"="+url.QueryEscape(fmt.Sprintf("%v", t)))
		case bool:
			p = append(p, k+"="+url.QueryEscape(fmt.Sprintf("%v", t)))
		case int, int32, int64, uint, uint32, uint64:
			p = append(p, k+"="+url.QueryEscape(fmt.Sprintf("%d", t)))
		case []string:
			sort.Strings(t)
			for _, arrv := range t {
				p = append(p, k+"="+url.QueryEscape(fmt.Sprintf("%v", arrv)))
			}
		case []int: // TODO：整形数组的处理
			sort.Ints(t)
			for _, arrv := range t {
				p = append(p, k+"="+url.QueryEscape(fmt.Sprintf("%d", arrv)))
			}
		case []int64: // TODO：整形数组的处理
			sort.Sort(util.Int64Slice(t))
			for _, arrv := range t {
				p = append(p, k+"="+url.QueryEscape(fmt.Sprintf("%d", arrv)))
			}
		case []uint32: // TODO：整形数组的处理
			sort.Sort(util.Uint32Slice(t))
			for _, arrv := range t {
				p = append(p, k+"="+url.QueryEscape(fmt.Sprintf("%d", arrv)))
			}
		case []interface{}: // TODO：整形数组的处理
			// IOS 数组都解析为 []interface{}，需要特殊处理
			iarry := util.ToIntSlice(t)
			if iarry != nil {
				sort.Ints(iarry)
				for _, arrv := range iarry {
					p = append(p, k+"="+url.QueryEscape(fmt.Sprintf("%d", arrv)))
				}
				break
			}

			sarry := util.ToStringSlice(t)
			if sarry != nil {
				sort.Strings(sarry)
				for _, arrv := range sarry {
					p = append(p, k+"="+url.QueryEscape(fmt.Sprintf("%v", arrv)))
				}
				break
			}
			logger.Errorf("not in sign data:%+v %+v", reflect.TypeOf(t[0]), v)
		default:
			logger.Errorf("not in sign data:%+v %+v", reflect.TypeOf(v), v)
		}
	}
	par, _ := url.ParseQuery(strings.Join(p, "&"))
	//url.Values Encode 对key 做了排序
	source := par.Encode()
	source, _ = url.QueryUnescape(source)

	return source
}

// 将输入参数组装成map, 只保留值为字符串和数字类型的键值。按key进行排序后组装成字符串进行md5加密
func MakeSignature(params map[string]interface{}, timestamp int64) (string, string, string) {
	source := getSignatureData(params)

	// 计算出密钥并加密
	t := uint64(timestamp)
	n := t % 64
	secret := fmt.Sprintf("%d", ((t<<n)&(0x7fffffffffffffff))|(t>>(64-n)))
	mac := hmac.New(sha1.New, []byte(secret))
	mac.Write([]byte(source))
	logger.Debugf("signature calculated: [%v] secret:%v", source, secret)

	return fmt.Sprintf("%X", mac.Sum(nil)), secret, source
}

func ParseParam(c *gin.Context) (map[string]interface{}, error) {
	var err error
	params := map[string]interface{}{}
	if ctype := c.Request.Header.Get("Content-Type"); len(ctype) > 0 &&
		strings.Contains(ctype, "application/json") {
		eparams, ok := c.Get("json_params")
		if !ok {
			resp, _ := ioutil.ReadAll(c.Request.Body)
			c.Request.Body.Close()
			c.Request.Body = util.NewBuffReadCloser(resp)
			if len(resp) > 0 {
				if err = json.Unmarshal(resp, &params); err == nil {
					c.Set(constant.ContextJsonParams, params)
					c.Set(constant.ContextJsonBody, string(resp))
					//return params, nil
					//不返回，继续解析 post url  参数
				} else {
					logger.Errorf("ParseParam:%s %s [%v],%v", c.Request.Method, c.Request.RequestURI, string(resp), err)
					//return nil, err
					//不返回，继续解析 post url  参数
				}
			}
		} else {
			params = eparams.(map[string]interface{})
			//return params, nil
			//不返回，继续解析 post url  参数
		}
	}

	postLen := 0
	if c.Request.Method == "POST" {
		err = c.Request.ParseForm()
		if err == nil {
			for k, v := range c.Request.PostForm {
				if len(v) == 1 {
					params[k] = v[0]
				} else {
					params[k] = v
				}
				postLen++
			}
		} else {
			logger.Error(err)
		}
	}
	if c.Request.Method == "GET" || postLen <= 0 {
		for k, v := range c.Request.URL.Query() {
			if len(v) == 1 {
				params[k] = v[0]
			} else {
				params[k] = v
			}
		}
	}
	//logger.Debug(params)
	c.Set(constant.ContextParams, params)
	return params, nil
}

func getInt(str interface{}) int {
	val, _ := strconv.Atoi(fmt.Sprint(str))
	return val
}

func getInt64(str interface{}, method string) int64 {
	var old string
	if method == "POST" {
		old = fmt.Sprintf("%v", str)
		var newVal float64
		newVal, _ = strconv.ParseFloat(old, 64)
		s := strconv.FormatFloat(newVal, 'f', -1, 64)
		i, _ := strconv.ParseInt(s, 10, 64)
		return i
	} else {
		old = fmt.Sprintf("%s", str)
		i, _ := strconv.ParseInt(old, 10, 64)
		return i
	}
}
