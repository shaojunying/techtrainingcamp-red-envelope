package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"red_envelope/api/redenvelope"
	"red_envelope/config"
	"red_envelope/database"
	"red_envelope/middleware"
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

func Snatch(t *testing.T, w *httptest.ResponseRecorder, ctx *gin.Context, uid int) (int, redenvelope.SuccessSnatch) {
	MockJsonPost(ctx, map[string]interface{}{"uid": uid})
	redenvelope.SnatchRedEnvelope(ctx)
	assert.Equal(t, w.Code, http.StatusOK)
	body, err := ioutil.ReadAll(w.Body)
	HandleERR(t, err)
	var datamap map[string]interface{}
	err = json.Unmarshal(body, &datamap)
	HandleERR(t, err)
	code := int(datamap["code"].(float64))
	if datamap["data"] == nil {
		return code, redenvelope.SuccessSnatch{}
	}
	data := datamap["data"].(redenvelope.SuccessSnatch)
	log.Printf("%d抢红包：%+v, data为%+v", uid, code, data)
	return code, data
}

func Open(t *testing.T, w *httptest.ResponseRecorder, ctx *gin.Context, uid int, envelope_id int) (int, redenvelope.SuccessOpen) {
	MockJsonPost(ctx, map[string]interface{}{"uid": uid, "envelope_id": envelope_id})
	redenvelope.OpenRedEnvelope(ctx)
	assert.Equal(t, w.Code, http.StatusOK)
	body, err := ioutil.ReadAll(w.Body)
	HandleERR(t, err)
	var datamap map[string]interface{}
	err = json.Unmarshal(body, &datamap)
	HandleERR(t, err)
	code := int(datamap["code"].(float64))
	if datamap["data"] == nil {
		return code, redenvelope.SuccessOpen{}
	}
	data := datamap["data"].(redenvelope.SuccessOpen)
	log.Printf("%d拆红包%+v：%+v, data为%+v", uid, envelope_id, code, data)
	return code, data
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
	code, _ := Snatch(t, w, ctx, 15)
	if code > 1 {
		t.Log("err, not expected code 0 or 1, the code is", code)
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
	code, _ := Open(t, w, ctx, 15, 800)
	if code != 0 && code != 4 {
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

	uid := 15
	red_envelope_cnt := 0
	value := 0
	// 期望抢30次红包可以走完这个流程
	i := 0
	for ; i < 30; i += 1 {
		code, sdata := Snatch(t, w, ctx, uid)
		if code == 0 {
			// 抢到红包了，尝试拆开
			red_envelope_cnt++
			code, odata := Open(t, w, ctx, uid, sdata.EnvelopeID)
			if code == 0 {
				value += odata.Value
				log.Printf("第%d次抢到红包并且成功拆开，目前获得金额%d", red_envelope_cnt, value)
			} else {
				HandleERR(t, fmt.Errorf("拆刚得到的红包%d时出错，code为%d", sdata.EnvelopeID, code))
			}
			if red_envelope_cnt > ctx.GetInt(redenvelope.MaxCountField) {
				HandleERR(t, fmt.Errorf("已经抢了%d次红包，超过正确限额%d", red_envelope_cnt, ctx.GetInt(redenvelope.MaxCountField)))
			}
		}
		if code == 2 {
			// 用户抢到限额了
			if red_envelope_cnt == sdata.CurCount && sdata.CurCount == sdata.MaxCount && sdata.CurCount == ctx.GetInt(redenvelope.MaxCountField) {
				break
			} else {
				HandleERR(t, fmt.Errorf("到达限额时red_envelope_cnt: %d, sdata.CurCount: %d, sdata.MaxCount: %d, MaxCountField: %d 不相等", red_envelope_cnt, sdata.CurCount, sdata.MaxCount, ctx.GetInt(redenvelope.MaxCountField)))
			}
		}
	}
	if i >= 30 {
		HandleERR(t, fmt.Errorf("已尝试抢了30次红包，但仍未达到限额，请检查"))
	}
	t.FailNow()
}
