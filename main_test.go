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

func Snatch(t *testing.T, uid int) (int, map[string]interface{}) {
	log.Printf("发抢红包请求：%d", uid)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = &http.Request{
		Header: make(http.Header),
	}
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
		return code, nil
	}
	data := datamap["data"].(map[string]interface{})
	log.Printf("%d抢红包：%+v, data为%+v", uid, code, data)
	return code, data
}

func Open(t *testing.T, uid int, envelope_id int) (int, map[string]interface{}) {
	log.Printf("发拆红包请求：%d，%d", uid, envelope_id)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = &http.Request{
		Header: make(http.Header),
	}
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
		return code, nil
	}
	data := datamap["data"].(map[string]interface{})
	log.Printf("%d拆红包%+v：%+v, data为%+v", uid, envelope_id, code, data)
	return code, data
}

func List(t *testing.T, uid int) (int, map[string]interface{}) {
	log.Printf("发红包列表请求：%d", uid)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = &http.Request{
		Header: make(http.Header),
	}
	MockJsonPost(ctx, map[string]interface{}{"uid": uid})
	redenvelope.GetWalletList(ctx)
	assert.Equal(t, w.Code, http.StatusOK)
	body, err := ioutil.ReadAll(w.Body)
	HandleERR(t, err)
	var datamap map[string]interface{}
	err = json.Unmarshal(body, &datamap)
	HandleERR(t, err)
	code := int(datamap["code"].(float64))
	if datamap["data"] == nil {
		return code, nil
	}
	data := datamap["data"].(map[string]interface{})
	log.Printf("%d获取红包列表：%+v, data为%+v", uid, code, data)
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

	code, _ := Snatch(t, 19)
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

	code, _ := Open(t, 20, 163)
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

	code, data := List(t, 19)
	if code != 0 {
		t.Log("err, not expected code, the code is", code)
		t.FailNow()
	}

	t.Log("总金额：", data["amount"])
	t.Log("红包列表：", data["envelope_list"])
}

func TestSetRedEnvelopeConfig(t *testing.T) {
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
	max_count := 999999
	probability := 0.5
	budget := 100000000
	total_number := 1000000
	min_value := 1
	max_value := 10000
	MockJsonPost(ctx, map[string]interface{}{"max_count": max_count, "probability": probability,
		"budget": budget, "total_number": total_number, "min_value": min_value, "max_value": max_value})
	redenvelope.SetRedEnvelopeConfig(ctx)
	assert.Equal(t, w.Code, http.StatusOK)
	body, err := ioutil.ReadAll(w.Body)
	HandleERR(t, err)
	var datamap map[string]interface{}
	err = json.Unmarshal(body, &datamap)
	HandleERR(t, err)
	code := int(datamap["code"].(float64))

	if code != 0 {
		t.Log("err, not expected code, the code is", code)
		t.FailNow()
	}
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

	uid := 25
	red_envelope_cnt := 0
	value := 0
	// 期望抢30次红包可以走完这个流程
	i := 0
	max_count := -1
	for ; i < 30; i += 1 {
		code, sdata := Snatch(t, uid)
		if sdata != nil {
			max_count = int(sdata["max_count"].(float64))
			envelope_id := int(sdata["envelope_id"].(float64))

			// 抢到红包了，尝试拆开
			red_envelope_cnt++
			code, odata := Open(t, uid, envelope_id)
			if code == 0 {
				value += int(odata["value"].(float64))
				log.Printf("第%d次抢到红包并且成功拆开，目前获得金额%d", red_envelope_cnt, value)
			} else {
				HandleERR(t, fmt.Errorf("拆刚得到的红包%d时出错，code为%d", envelope_id, code))
			}
			if red_envelope_cnt > max_count {
				HandleERR(t, fmt.Errorf("已经抢了%d次红包，超过正确限额%d", red_envelope_cnt, max_count))
			}
		} else if code == 2 {
			// 用户抢到限额了
			if red_envelope_cnt == max_count {
				break
			} else {
				HandleERR(t, fmt.Errorf("到达限额时red_envelope_cnt: %d, max_count: %d 不相等，-1可能是根本没抢到红包就到限额了", red_envelope_cnt, max_count))
			}
		} else if code != 1 {
			HandleERR(t, fmt.Errorf("抢红包时出错，code为%d", code))
		}
	}
	if i >= 30 {
		HandleERR(t, fmt.Errorf("已尝试抢了30次红包，但仍未达到限额，请检查"))
	}
	code, data := List(t, uid)
	if code != 0 {
		HandleERR(t, fmt.Errorf("获取红包列表时出错，code为%d", code))
	}
	// 确保红包金额相等
	assert.Equal(t, int(data["amount"].(float64)), value)
	// 确保红包列表数量相等
	assert.Equal(t, len(data["envelope_list"].([]interface{})), red_envelope_cnt)
}
