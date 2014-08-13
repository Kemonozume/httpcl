package httpcl

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"
)

type UserAgent struct {
	Name string `json:"user-agent"`
}

func transform_to_useragent(resp *http.Response, c interface{}) (err error) {
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return errors.New("status not 200")
	}
	user := c.(*UserAgent)
	by, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(by, user)
	return
}

func Test_DoTransform(t *testing.T) {
	var agent UserAgent
	err := Get("http://httpbin.org/user-agent").
		SetUserAgent("httpcl").
		DoTransform(transform_to_useragent, &agent)
	if err != nil {
		t.Error(err.Error())
	}

	if agent.Name != "httpcl" {
		t.Errorf("Name should be \"%v\" is \"%v\"", "httpcl", agent.Name)
	}

	err2 := Get("http://httpbin.org/user-agen").SetUserAgent("httpcl").DoTransform(transform_to_useragent, &agent)
	if err2.Error() != "status not 200" {
		t.Error(err2.Error())
	}

	cl := &Client{}
	cl.SetUserAgent("httpcl")
	cl.request, _ = http.NewRequest("GET", "http://httpbin.org/user-agent", nil)
	err3 := cl.DoTransform(transform_to_useragent, &agent)
	if err3 == nil {
		t.Error("should fail because of no request")
	}
}

func Test_DoTransformJson(t *testing.T) {
	var json map[string]interface{}
	err := Get("http://httpbin.org/user-agent").
		SetUserAgent("httpcl").
		DoTransform(TransformToJson, &json)

	if err != nil {
		t.Error(err.Error())
	}
	if json["user-agent"] != "httpcl" {
		t.Error("user-agent should be \"httpcl\" is \"%s\"", json["user-agent"])
	}
}

func Test_DoTransformString(t *testing.T) {
	var str string
	err := Get("http://httpbin.org/user-agent").
		SetUserAgent("httpcl").
		DoTransform(TransformToString, &str)

	if err != nil {
		t.Error(err.Error())
	}
	if str == "" {
		t.Error("string transform failed")
	}
}

func Test_DoTransformStringFail(t *testing.T) {
	var c interface{}
	err := Get("http://httpbin.org/user-agent").
		SetUserAgent("httpcl").
		DoTransform(TransformToString, c)

	if err == nil {
		t.Error("should fail c is not a string")
	}
}
