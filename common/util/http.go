package util

import (
	"bytes"
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"iptv/common/logger"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"

	jsoniter "github.com/json-iterator/go"
	"github.com/json-iterator/go/extra"
)

func init() {
	extra.RegisterFuzzyDecoders()
}

/*
	this method allows to send a request and bind the response to the provided json model.
*/
func GetWithError(aurl string, outObj interface{}, errModel interface{}, headers map[string]string) error {
	req, err := http.NewRequest("GET", aurl, nil)
	req.Close = true
	for k, v := range headers {
		req.Header.Add(k, v)
	}

	transport := &http.Transport{DisableKeepAlives: true, //短链接
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	client := &http.Client{Transport: transport}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error(aurl + "," + err.Error())
		return err
	}

	defer resp.Body.Close()

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(aurl + "," + err.Error())
		return err
	}

	if len(respBytes) > 0 {
		// CAN BE REMOVED
		//logger.Debug("",string(respBytes))
		err = jsoniter.Unmarshal(respBytes, outObj)
		if err != nil {
			err2 := jsoniter.Unmarshal(respBytes, errModel)
			if err2 != nil {
				logger.Errorf("GetWithError:%v  return:%v format error:%v", aurl, string(respBytes), err)
				return err
			}
		} else {
			err2 := jsoniter.Unmarshal(respBytes, errModel)
			if err2 != nil {
				logger.Errorf("GetWithError:%v  return:%v format error:%v", aurl, string(respBytes), err)
				return err
			}
		}
	}

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad response status [%s] data:%s!", resp.Status, showErrorData(respBytes))
		logger.Error(aurl + "," + err.Error())
		return err
	}

	return nil
}

func Get(aurl string, outObj interface{}, headers map[string]string) error {
	req, err := http.NewRequest("GET", aurl, nil)
	//短链接
	req.Close = true
	for k, v := range headers {
		req.Header.Add(k, v)
	}

	logger.Debugf("GET:%v headers:%v ", aurl, headers)
	transport := &http.Transport{
		DisableKeepAlives: true, //短链接
		TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: transport}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error(aurl + "," + err.Error())
		return err
	}

	defer resp.Body.Close()

	logger.Debugf("%+v", resp.Header)
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(aurl + "," + err.Error())
		return err
	}

	if len(respBytes) > 0 {
		// CAN BE REMOVED
		if len(respBytes) > 128 {
			logger.Debugf("GET:%v  from :%v %v ...", aurl, resp.Header["Service-Provider"], resp.Header["Server-Version"])
			logger.Debugf("GET:%v  return:%v ...", aurl, string(respBytes[:128]))
		} else {
			logger.Debugf("GET:%v  return:%v ...", aurl, string(respBytes))
		}
		if err = jsoniter.Unmarshal(respBytes, outObj); err != nil {
			logger.Errorf("GET:%v  return:%v format error:%v", aurl, string(respBytes), err)
			return err
		}
	}

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad response status [%s] data:%s!", resp.Status, showErrorData(respBytes))
		logger.Error(aurl + "," + err.Error())
		return err
	}

	return nil
}

func showErrorData(respBytes []byte) string {
	if len(respBytes) > 128 {
		return string(respBytes[:128])
	}

	return string(respBytes)
}

func GetWithProxy(aurl string, outObj interface{}, headers map[string]string) error {
	var client *http.Client
	httpProxy := os.Getenv("THIRD_HTTP_PROXY")
	if len(httpProxy) > 0 {
		proxy := func(_ *http.Request) (*url.URL, error) {
			return url.Parse(httpProxy) //根据定义Proxy func(*Request) (*url.URL, error)这里要返回url.URL
		}
		transport := &http.Transport{DisableKeepAlives: true, //短链接
			Proxy: proxy}
		client = &http.Client{Transport: transport}
		logger.Debug("Proxy:" + httpProxy)
	} else {
		transport := &http.Transport{DisableKeepAlives: true, //短链接
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
		client = &http.Client{Transport: transport}
	}
	req, err := http.NewRequest("GET", aurl, nil)
	req.Close = true
	for k, v := range headers {
		req.Header.Add(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		logger.Error(aurl + "," + err.Error())
		return err
	}

	defer resp.Body.Close()
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(aurl + "," + err.Error())
		return err
	}
	// CAN BE REMOVED
	// logger.Debug(string(respBytes))
	if err = jsoniter.Unmarshal(respBytes, outObj); err != nil {
		logger.Errorf("GetWithProxy:%v  return:%v format error:%v", aurl, string(respBytes), err)
		return err
	}

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad response status [%s]  data:%s!", resp.Status, showErrorData(respBytes))
		logger.Error(err)
		return err
	}

	return nil
}

/*
	this method allows to send a request and bind the response to the provided xml model.
*/
func GetXml(aurl string, outObj interface{}, headers map[string]string) error {
	req, err := http.NewRequest("GET", aurl, nil)
	req.Close = true
	for k, v := range headers {
		req.Header.Add(k, v)
	}

	transport := &http.Transport{DisableKeepAlives: true, //短链接
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	client := &http.Client{Transport: transport}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error(aurl + "," + err.Error())
		return err
	}

	defer resp.Body.Close()
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(aurl + "," + err.Error())
		return err
	}
	// CAN BE REMOVED
	// logger.Debug(string(respBytes))
	if err = xml.Unmarshal(respBytes, outObj); err != nil {
		logger.Errorf("GETXml:%v  return:%v format error:%v", aurl, string(respBytes), err)
		return err
	}

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad response status [%s]  data:%s!", resp.Status, showErrorData(respBytes))
		logger.Error(aurl + "," + err.Error())
		return err
	}

	return nil
}

/*
	This method allows to post a json input and retrieve the response.
*/
func Post(aurl string, input interface{}, outObj interface{}, headers map[string]string) error {
	jsonInBytes, err := jsoniter.Marshal(input)
	if err != nil {
		logger.Error(aurl + "," + err.Error())
		return err
	}

	logger.Debug("POST:", string(jsonInBytes))
	req, err := http.NewRequest("POST", aurl, bytes.NewReader(jsonInBytes))
	req.Close = true
	req.Header.Set("accept", "application/json")
	req.Header.Set("content-type", "application/json; charset=UTF-8")

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	transport := &http.Transport{DisableKeepAlives: true, //短链接
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	client := &http.Client{Transport: transport}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error(aurl + "," + err.Error())
		return err
	}

	defer resp.Body.Close()
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(aurl + "," + err.Error())
		return err
	}

	if resp.StatusCode != http.StatusNotFound {
		if len(respBytes) > 0 {
			if len(respBytes) > 128 {
				logger.Debugf("POST:%v  return:%v ...", aurl, string(respBytes[:128]))
			} else {
				logger.Debugf("POST:%v  return:%v ...", aurl, string(respBytes))
			}
			if err = jsoniter.Unmarshal(respBytes, outObj); err != nil {
				logger.Errorf("POST:%v  return:%v format error:%v", aurl, string(respBytes), err)
				return err
			}
		}
	}

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad response status [%s] data:%s!", resp.Status, showErrorData(respBytes))
		logger.Error(aurl + "," + err.Error())
		return err
	}
	//logger.Debug("POST:", string(respBytes), "return:", resp.StatusCode, string(respBytes))
	return nil
}

/*
	This method allows to post a form input and retrieve the response.
*/
func PostFile(aurl string, formFile string, filename string, input []byte, outObj interface{}, headers map[string]string) error {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	fileWriter, err := bodyWriter.CreateFormFile(formFile, filename)
	if err != nil {
		logger.Error("error writing to buffer", aurl+","+err.Error())
		return err
	}
	fileWriter.Write(input)
	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	resp, err := http.Post(aurl, contentType, bodyBuf)
	if err != nil {
		logger.Error(aurl + "," + err.Error())
		return err
	}

	defer resp.Body.Close()
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(aurl + "," + err.Error())
		return err
	}
	if err = jsoniter.Unmarshal(respBytes, outObj); err != nil {
		logger.Errorf("PostFile:%v  return:%v format error:%v", aurl, string(respBytes), err)
		return err
	}

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad response status [%s] ! data:%s", resp.Status, showErrorData(respBytes))
		logger.Error(aurl + "," + err.Error())
		return err
	}

	return nil
}

/*表单POST请求*/
func PostForm(aurl string, paramStr string, outObj interface{}, headers map[string]string) error {
	req, _ := http.NewRequest("POST", aurl, strings.NewReader(paramStr))
	req.Header.Set("accept", "application/json")
	req.Header.Set("content-type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Close = true
	for k, v := range headers {
		req.Header.Add(k, v)
	}

	transport := &http.Transport{DisableKeepAlives: true, //短链接
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	client := &http.Client{Transport: transport}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error(aurl + "," + err.Error())
		return err
	}

	defer resp.Body.Close()
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(aurl + "," + err.Error())
		return err
	}
	if len(respBytes) > 0 {
		//logger.Error(showErrorData(respBytes))
		if err = jsoniter.Unmarshal(respBytes, outObj); err != nil {
			logger.Errorf("PostForm:%v  return:%v format error:%v", aurl, string(respBytes), err)
			return err
		}
	}

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad response status [%s] data:%s!", resp.Status, showErrorData(respBytes))
		logger.Error(aurl + "," + err.Error())
		return err
	}

	return nil
}

func Delete(aurl string, outObj interface{}, headers map[string]string) error {
	req, _ := http.NewRequest("Delete", aurl, nil)
	req.Header.Set("accept", "application/json")
	req.Header.Set("content-type", "application/json; charset=UTF-8")
	req.Close = true
	for k, v := range headers {
		req.Header.Add(k, v)
	}

	transport := &http.Transport{DisableKeepAlives: true, //短链接
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	client := &http.Client{Transport: transport}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error(aurl + "," + err.Error())
		return err
	}

	defer resp.Body.Close()
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(aurl + "," + err.Error())
		return err
	}

	if len(respBytes) > 0 {
		if err = jsoniter.Unmarshal(respBytes, outObj); err != nil {
			logger.Errorf("DELETE:%v  return:%v format error:%v", aurl, string(respBytes), err)
			return err
		}
	}

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad response status [%s] data:%s!", resp.Status, showErrorData(respBytes))
		logger.Error(aurl + "," + err.Error())
		return err
	}

	return nil
}
