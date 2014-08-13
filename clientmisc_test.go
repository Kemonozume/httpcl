package httpcl

import "testing"

func Test_Put(t *testing.T) {
	resp, err := Put("http://jsonplaceholder.typicode.com/posts/1", "id", 255, "title", "foo", "body", "bar", "userId", 255).Do()
	if err != nil {
		t.Error(err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("StatusCode should be 200 is %v", resp.StatusCode)
	}
}

func Test_Patch(t *testing.T) {
	resp, err := Patch("http://jsonplaceholder.typicode.com/posts/1", "title", "httpcl").Do()
	if err != nil {
		t.Error(err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("StatusCode should be 200 is %v", resp.StatusCode)
	}
}

func Test_Delete(t *testing.T) {
	resp, err := Delete("http://jsonplaceholder.typicode.com/posts/1").Do()
	if err != nil {
		t.Error(err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 && resp.StatusCode != 204 && resp.StatusCode != 202 {
		t.Errorf("StatusCode should be 200|204|202 is %v", resp.StatusCode)
	}
}

func Test_Head(t *testing.T) {
	var str string
	err := Head("http://www.google.com").DoTransform(TransformToString, &str)
	if err != nil {
		t.Error(err.Error())
	}

	if str != "" {
		t.Error("response should be empty")
	}
}
