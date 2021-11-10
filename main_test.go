package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"red_envelope/api/redenvelope"
	"red_envelope/config"
	"red_envelope/database"
	"red_envelope/middleware"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/magiconair/properties/assert"
)

func MockJsonPost(c *gin.Context /* the test context */, jsonMap interface{}) {
	middleware.LoadConfig(c)
	c.Request.Method = "POST" // or PUT
	c.Request.Header.Set("Content-Type", "application/json")

	jsonBytes, err := json.Marshal(jsonMap)
	if err != nil {
		panic(err)
	}
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(jsonBytes))
}

func HandleERR(t *testing.T, err error) {
	if err == nil {
		return
	}
	t.Log(err)
	t.FailNow()
}

func TestSnatchRedEnvelope(t *testing.T) {
	//读取配置
	config.InitConf()

	//启动数据库
	db := database.InitDB()
	defer db.Close()
	database.InitRedis()

	err := database.InitMQ()
	HandleERR(t, err)
	defer database.CloseMQ()

	w := httptest.NewRecorder()
	ctx, r := gin.CreateTestContext(w)
	r.Use(middleware.ConfigLoadingMiddleware())

	ctx.Request = &http.Request{
		Header: make(http.Header),
	}

	MockJsonPost(ctx, map[string]interface{}{"uid": 15})

	redenvelope.SnatchRedEnvelope(ctx)
	assert.Equal(t, w.Code, http.StatusOK)
	body, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Log(err)
	}
	var data map[string]interface{}
	if err = json.Unmarshal(body, &data); err != nil {
		t.Log(err)
	}
	code := int(data["code"].(float64))
	if !(reflect.DeepEqual(0, code)) {
		t.Log("err, not expected code, the code is", code)
		t.FailNow()
	}
}

func TestOpenRedEnvelope(t *testing.T) {
	//读取配置
	config.InitConf()

	//启动数据库
	db := database.InitDB()
	defer db.Close()
	database.InitRedis()

	err := database.InitMQ()
	HandleERR(t, err)
	defer database.CloseMQ()

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	ctx.Request = &http.Request{
		Header: make(http.Header),
	}

	MockJsonPost(ctx, map[string]interface{}{"uid": 15, "envelope_id": 800})

	redenvelope.OpenRedEnvelope(ctx)
	assert.Equal(t, w.Code, http.StatusOK)
	body, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Log(err)
	}
	var data map[string]interface{}
	if err = json.Unmarshal(body, &data); err != nil {
		t.Log(err)
	}
	code := int(data["code"].(float64))
	if !(reflect.DeepEqual(0, code)) {
		t.Log("err, not expected code, the code is", code)
		t.FailNow()
	}
}

func TestGetWalletList(t *testing.T) {
	//读取配置
	config.InitConf()

	//启动数据库
	db := database.InitDB()
	defer db.Close()
	database.InitRedis()

	err := database.InitMQ()
	HandleERR(t, err)
	defer database.CloseMQ()

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	ctx.Request = &http.Request{
		Header: make(http.Header),
	}
	MockJsonPost(ctx, map[string]interface{}{"uid": 15})

	redenvelope.GetWalletList(ctx)
	var data map[string]interface{}

	body, err := ioutil.ReadAll(w.Body)
	HandleERR(t, err)
	json.Unmarshal(body, &data)
	assert.Equal(t, w.Code, http.StatusOK)
}

func TestWorkflow(t *testing.T) {

}
