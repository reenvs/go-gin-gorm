package util

import (
	"encoding/json"
	"io/ioutil"
	"iptv/common/logger"

	"github.com/gin-gonic/gin"
)

func GetRequestInfo(c *gin.Context, p interface{}) error {
	var err error
	b, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		logger.Error(err)
		return err
	}
	logger.Debug(string(b))
	if err = json.Unmarshal(b, p); err != nil {
		logger.Error(err)
		return err
	}
	return nil
}
