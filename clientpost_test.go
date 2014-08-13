package httpcl

import (
	"encoding/json"
	"io/ioutil"
	"net/url"
	"os"
	"strings"
	"testing"
)

func Test_PostSimple(t *testing.T) {
	cl := Post("http://httpbin.org/post", "test", "value", "test1", 1)
	if cl.Error != nil {
		t.Error(cl.Error.Error())
	}
	resp, err := cl.Do()
	if err != nil {
		t.Error(err.Error())
	}

	defer resp.Body.Close()
	by, _ := ioutil.ReadAll(resp.Body)

	var i interface{}
	err = json.Unmarshal(by, &i)
	if err != nil {
		t.Error(err.Error())
	}

	m := i.(map[string]interface{})
	data := m["data"].(string)
	verifyPost(t, data, 2, "test=value&test1=1")
}

func Test_PostIoReader(t *testing.T) {
	f, err := os.Open("postioreader")
	if err != nil {
		t.Error(err.Error())
	}
	defer f.Close()

	cl := Post("http://httpbin.org/post", f)
	if cl.Error != nil {
		t.Error(cl.Error.Error())
	}

	resp, err1 := cl.Do()
	if err1 != nil {
		t.Error(err1.Error())
	}

	defer resp.Body.Close()
	by, _ := ioutil.ReadAll(resp.Body)

	var i interface{}
	err = json.Unmarshal(by, &i)
	if err != nil {
		t.Error(err.Error())
	}

	m := i.(map[string]interface{})
	data := m["data"].(string)
	verifyPost(t, data, 2, "test=value&test1=1")
}

func Test_PostValues(t *testing.T) {
	values := url.Values{}
	values.Add("test", "value")
	values.Add("test1", "1")
	resp, err := Post("http://httpbin.org/post", values).Do()
	if err != nil {
		t.Error(err.Error())
	}

	defer resp.Body.Close()
	by, _ := ioutil.ReadAll(resp.Body)

	var i interface{}
	err = json.Unmarshal(by, &i)
	if err != nil {
		t.Error(err.Error())
	}

	m := i.(map[string]interface{})
	data := m["data"].(string)
	verifyPost(t, data, 2, "test=value&test1=1")
}

func Test_PostFail(t *testing.T) {
	err := Post("http://httpbin.org/post", "test", "test", "test").Error
	if err.Error() != "parameters not correct, expected 4 parameter got 3" {
		t.Error(err.Error())
	}

	err1 := Post("http://httpbin.org/post", "test").Error
	if err1.Error() != "parameters not correct string" {
		t.Error(err1.Error())
	}

	var i int16
	i = 1
	vals := map[string]interface{}{
		"test": i,
	}
	err2 := Post("http://httpbin.org/post", vals).Error
	if err2.Error() != "unsupported type for post int16" {
		t.Error(err2.Error())
	}

	err3 := Post("http://httpbin.org/post", "test", int16(1)).Error
	if err3.Error() != "unsupported type for post int16" {
		t.Error(err3.Error())
	}
}

func Test_PostEmpty(t *testing.T) {
	resp, err := Post("http://httpbin.org/post").Do()
	if err != nil {
		t.Error(err.Error())
	}

	defer resp.Body.Close()
	by, _ := ioutil.ReadAll(resp.Body)

	var i interface{}
	err = json.Unmarshal(by, &i)
	if err != nil {
		t.Error(err.Error())
	}

	m := i.(map[string]interface{})
	data := m["data"].(string)
	if data != "" {
		t.Error("post params should be empty, is %v", data)
	}

}

func Test_PostMap(t *testing.T) {
	var test6 float64
	test6 = 2.344
	vals := map[string]interface{}{
		"test":  "value2",
		"test1": 2,
		"test2": true,
		"test3": int64(2),
		"test4": rune('a'),
		"test5": uint64(20),
		"test6": test6,
	}
	resp, err := Post("http://httpbin.org/post", vals).Do()
	if err != nil {
		t.Error(err.Error())
	}

	defer resp.Body.Close()
	by, _ := ioutil.ReadAll(resp.Body)

	var i interface{}
	err = json.Unmarshal(by, &i)
	if err != nil {
		t.Error(err.Error())
	}

	m := i.(map[string]interface{})
	data := m["data"].(string)
	verifyPost(t, data, 7, "test=value2&test1=2&test2=true&test3=2&test4=a&test5=20&test6=2.344000")
}

func verifyPost(t *testing.T, data string, length int, actual string) {
	actlength := len(strings.Split(data, "&"))
	if actlength != length {
		t.Errorf("length should be %v is %v", length, actlength)
	}

	if data != actual {
		t.Errorf("data not correct: %v", data)
	}
}
