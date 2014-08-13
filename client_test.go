package httpcl

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
)

func Test_Do(t *testing.T) {
	cl := &Client{}
	_, err := cl.Do()
	if err == nil {
		t.Error("do without request should fail")
	}
	_, err = Get("").Do()
	if err == nil {
		t.Error("do without url should fail")
	}
	_, err = Get("http://www.google.d").Do()
	if err == nil {
		t.Error("do with false url should fail")
	}

	cl = &Client{}
	cl.FollowRedirect(false)
	cl.request, _ = http.NewRequest("GET", "http://www.google.com", nil)
	_, err = cl.Do()
	if err == nil {
		t.Error("follow redirect before request should fail")
	}

	_, err = Get("127.0.0.1").FollowRedirect(false).Do()
	if err == nil {
		t.Error("redirect ??")
	}

}

func Test_Get200(t *testing.T) {
	cl := Get("http://httpbin.org/status/200")
	if cl.Error != nil {
		t.Error(cl.Error.Error())
	}
	resp, err := cl.Do()
	if err != nil {
		t.Error(err.Error())
	}
	if resp == nil {
		t.Error("response is nil")
	}
	if cl.StatusCode != 200 {
		t.Errorf("statuscode: %v error: %v", cl.StatusCode, cl.Error)
	}
}

func Test_Get404(t *testing.T) {
	cl := Get("http://httpbin.org/status/404")
	if cl.Error != nil {
		t.Error(cl.Error.Error())
	}
	resp, err := cl.Do()
	if err != nil {
		t.Error(err.Error())
	}
	if resp == nil {
		t.Error("response is nil")
	}
	if cl.StatusCode != 404 {
		t.Errorf("statuscode: %v error: %v", cl.StatusCode, cl.Error)
	}
}

func Test_BasicAuth(t *testing.T) {
	cl := Get("http://httpbin.org/basic-auth/user/passwd").SetBasicAuth("user", "passwd")
	cl.Do()
	if cl.StatusCode != 200 {
		t.Error("statuscode should be 200 is %v", cl.StatusCode)
	}
}

func Test_Cookie(t *testing.T) {
	cl := Get("http://httpbin.org/cookies")
	if cl.Error != nil {
		t.Error(cl.Error.Error())
	}

	cookie := &http.Cookie{
		Name:       "test",
		Value:      "test",
		Domain:     "www.google.com",
		Path:       "/",
		RawExpires: "Mon, 18 Aug 2018 14:47:29 GMT",
	}

	err := cl.AddCookie(cookie).Error
	if err != nil {
		t.Error(err.Error())
	}

	req := cl.GetRequest()
	if req == nil {
		t.Error("request shouldn't be empty")
	}

	length := len(req.Cookies())
	if length != 1 {
		t.Errorf("cookies length should be 1 is %v", length)
	}

	resp, err := cl.Do()
	if err != nil {
		t.Error(err.Error())
	}

	defer resp.Body.Close()
	by, err := ioutil.ReadAll(resp.Body)

	var respjson interface{}
	err = json.Unmarshal(by, &respjson)
	if err != nil {
		t.Error(err.Error())
	}

	m := respjson.(map[string]interface{})
	mcookies := m["cookies"].(map[string]interface{})

	if mcookies["test"] != "test" {
		t.Errorf("test should be \"test\" is \"%v\"", mcookies["test"])
	}

}

func Test_Cookies(t *testing.T) {
	cl := Get("http://httpbin.org/cookies")
	if cl.Error != nil {
		t.Error(cl.Error.Error())
	}

	cookie1 := &http.Cookie{
		Name:       "test",
		Value:      "test",
		Domain:     "www.google.com",
		Path:       "/",
		RawExpires: "Mon, 18 Aug 2018 14:47:29 GMT",
	}

	cookie2 := &http.Cookie{
		Name:       "test1",
		Value:      "test1",
		Domain:     "www.google.com",
		Path:       "/",
		RawExpires: "Mon, 18 Aug 2018 14:47:29 GMT",
	}

	cookies := []*http.Cookie{
		cookie1,
		cookie2,
	}

	err := cl.AddCookies(cookies).Error
	if err != nil {
		t.Error(err.Error())
	}

	req := cl.GetRequest()
	if req == nil {
		t.Error("request shouldn't be empty")
	}

	length := len(req.Cookies())
	if length != 2 {
		t.Errorf("cookies length should be 1 is %v", length)
	}

	resp, err := cl.Do()
	if err != nil {
		t.Error(err.Error())
	}

	defer resp.Body.Close()
	by, err := ioutil.ReadAll(resp.Body)

	var respjson interface{}
	err = json.Unmarshal(by, &respjson)
	if err != nil {
		t.Error(err.Error())
	}

	m := respjson.(map[string]interface{})
	mcookies := m["cookies"].(map[string]interface{})

	if mcookies["test"] != "test" {
		t.Errorf("test should be \"test\" is \"%v\"", mcookies["test"])
	}

	if mcookies["test1"] != "test1" {
		t.Errorf("test1 should be \"test1\" is \"%v\"", mcookies["test1"])
	}
}

func Test_Headers(t *testing.T) {
	cl := Get("http://httpbin.org/headers")
	if cl.Error != nil {
		t.Error(cl.Error.Error())
	}

	headers := map[string]string{
		"Referer":    "httpcl",
		"User-Agent": "httpcl",
	}

	err := cl.AddHeader("Connection", "close").Error
	if err != nil {
		t.Error(err.Error())
	}
	err = cl.AddHeaderMap(headers).Error
	if err != nil {
		t.Error(err.Error())
	}

	resp, err := cl.Do()
	if err != nil {
		t.Error(err.Error())
	}
	defer resp.Body.Close()
	by, err := ioutil.ReadAll(resp.Body)

	var respjson interface{}
	err = json.Unmarshal(by, &respjson)
	if err != nil {
		t.Error(err.Error())
	}

	m := respjson.(map[string]interface{})
	headersact := m["headers"].(map[string]interface{})

	if headersact["Connection"] != "close" {
		t.Errorf("Connection should be \"close\" is \"%v\"", headersact["Connection"])
	}
	if headersact["Referer"] != "httpcl" {
		t.Errorf("Referer should be \"httpcl\" is \"%v\"", headersact["Referer"])
	}
	if headersact["User-Agent"] != "httpcl" {
		t.Errorf("User-Agent should be \"httpcl\" is \"%v\"", headersact["User-Agent"])
	}
}

func Test_UserAgent(t *testing.T) {
	resp, err := Get("http://httpbin.org/user-agent").SetUserAgent("httpcl").Do()
	if err != nil {
		t.Error(err.Error())
	}
	defer resp.Body.Close()
	by, err := ioutil.ReadAll(resp.Body)

	var respjson interface{}
	err = json.Unmarshal(by, &respjson)
	if err != nil {
		t.Error(err.Error())
	}

	agent := respjson.(map[string]interface{})["user-agent"]

	if agent != "httpcl" {
		t.Error("agent should be \"httpcl\" is \"%v\"", agent)
	}

}

func Test_GetClient(t *testing.T) {
	cl := Get("http://www.google.com")
	if cl.Error != nil {
		t.Error(cl.Error.Error())
	}
	client := cl.GetClient()
	if client != nil {
		t.Error("client should be nil")
	}

	cl.Do()
	client2 := cl.GetClient()
	if client2 == nil {
		t.Error("client shouldn't be nil")
	}
}

func Test_SetClientRequest(t *testing.T) {
	cl1 := &Client{}
	cl2 := Get("http://www.google.com")
	cl2.SetClient(&http.Client{})

	cl1.SetRequest(cl2.GetRequest())
	cl1.SetClient(cl2.GetClient())

	resp, err := cl1.Do()
	if err != nil {
		t.Error(err.Error())
	}

	if resp.StatusCode != 200 {
		t.Errorf("statuscode should be 200 is %v", resp.StatusCode)
	}
	resp.Body.Close()

}

func Test_FollowRedirect(t *testing.T) {
	resp, err := Get("http://httpbin.org/redirect-to?url=http://example.com/").
		FollowRedirect(false).
		Do()
	if err != nil {
		t.Error(err.Error())
	}

	url, err := resp.Location()
	if err != nil {
		t.Error(err.Error())
	}

	if url.String() != "http://example.com/" {
		t.Errorf("url should be \"%s\" is \"%s\"", "http://example.com/", url.String())
	}
}

func Test_RunWith(t *testing.T) {
	cl := &Client{}
	err := cl.hasRequest()
	if err == nil {
		t.Error("client should have no request")
	}

	cl = cl.runWithHasRequest(func() {
		cl.FollowRedirect(true)
	})

	if cl.Error == nil {
		t.Error("error should be set because of no request")
	}

	if cl.Error.Error() != "no request" {
		t.Error(err.Error())
	}

}
