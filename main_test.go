package main

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/magiconair/properties/assert"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"red_envelope/api/redenvelope"
	"red_envelope/config"
	"red_envelope/database"
	"reflect"
	"testing"
)

func MockJsonPost(c *gin.Context /* the test context */, jsonMap interface{}) {
	c.Request.Method = "POST" // or PUT
	c.Request.Header.Set("Content-Type", "application/json")

	jsonBytes, err := json.Marshal(jsonMap)
	if err != nil {
		panic(err)
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonBytes))
}

func TestSnatchRedEnvelope(t *testing.T) {
	//读取配置
	config.InitConf()

	//启动数据库
	db := database.InitDB()
	defer db.Close()

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	ctx.Request = &http.Request{
		Header: make(http.Header),
	}

	MockJsonPost(ctx, map[string]interface{}{"uid": 15})

	redenvelope.SnatchRedEnvelope(ctx)
	assert.Equal(t, http.StatusOK, w.Code)
	body, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Log(err)
	}
	var data map[string]interface{}
	if err = json.Unmarshal(body, &data); err != nil {
		t.Log(err)
	}
	code := int(data["code"].(float64))
	if !(reflect.DeepEqual(0, code) || reflect.DeepEqual(1, code)) {
		t.Log("err, code is not expected")
	}
}

func TestOpenRedEnvelope(t *testing.T) {
	//读取配置
	config.InitConf()

	//启动数据库
	db := database.InitDB()
	defer db.Close()

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	ctx.Request = &http.Request{
		Header: make(http.Header),
	}

	MockJsonPost(ctx, map[string]interface{}{"uid": 15, "envelope_id": 18})

	redenvelope.OpenRedEnvelope(ctx)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetWalletList(t *testing.T) {
	//读取配置
	config.InitConf()

	//启动数据库
	db := database.InitDB()
	defer db.Close()

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	ctx.Request = &http.Request{
		Header: make(http.Header),
	}

	MockJsonPost(ctx, map[string]interface{}{"uid": 15})

	redenvelope.GetWalletList(ctx)
	assert.Equal(t, http.StatusOK, w.Code)
}
