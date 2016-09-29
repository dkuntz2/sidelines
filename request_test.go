package sidelines

import (
	"github.com/jarcoal/httpmock"
	. "github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
)

func repeatResponse(req *http.Request) (*http.Response, error) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	req.Body.Close()

	return httpmock.NewBytesResponse(200, body), nil
}

func TestRequest(t *testing.T) {
	Convey("NewRequest", t, func() {
		Convey("with GET endpoint", func() {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()

			httpmock.RegisterResponder("GET", "example.com",
				httpmock.NewStringResponder(200, `<h1>hello world!</h1>`))

			endpointUrl, _ := url.Parse("example.com")
			endpoint := &Endpoint{
				Method: "GET",
				Url:    endpointUrl,
			}

			response, err := endpoint.Request()
			So(err, ShouldBeNil)
			So(response.StatusCode, ShouldEqual, 200)
			So(string(response.Body), ShouldEqual, "<h1>hello world!</h1>")
		})

		Convey("with POST endpoint", func() {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()

			httpmock.RegisterResponder("POST", "example.com", repeatResponse)

			body := url.Values{}
			body.Add("foo", "bar")
			body.Add("what", "yolo")

			endpointUrl, _ := url.Parse("example.com")
			endpoint := &Endpoint{
				Method:   "POST",
				Url:      endpointUrl,
				PostData: body,
			}

			response, err := endpoint.Request()
			So(err, ShouldBeNil)
			So(response.StatusCode, ShouldEqual, 200)
			So(string(response.Body), ShouldEqual, body.Encode())
		})

		Convey("with a Proxy", func() {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()

			httpmock.RegisterResponder(
				"GET", "https://example.com",
				httpmock.NewStringResponder(200, `<h1>through a proxy!!</h1>`))
			httpmock.RegisterResponder("CONNECT", "proxy.example.com:1080", func(req *http.Request) (*http.Response, error) {
				println(req.Header)
				body, _ := ioutil.ReadAll(req.Body)
				println(string(body))
				return nil, nil
			})

			endpointUrl, _ := url.Parse("https://example.com")
			proxyUrl, _ := url.Parse("proxy.example.com:1080")
			endpoint := &Endpoint{
				Method:   "GET",
				Url:      endpointUrl,
				ProxyUrl: proxyUrl,
			}

			response, err := endpoint.Request()
			So(err, ShouldBeNil)
			So(response.StatusCode, ShouldEqual, 200)
			So(string(response.Body), ShouldEqual, "<h1>hello world!</h1>")
		})
	})
}
