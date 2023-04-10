package httpserver

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/smartystreets/goconvey/convey"
)

func TestStudentGetHttpRequest(t *testing.T) {
	convey.Convey("TestStudentGetHttpRequest", t, func() {
		convey.Convey("normal test", func() {
			// use sender
			ctx := context.Background()
			req := StudentReq{ID: 1, Name: "xxx"}
			r, _ := StudentGetHttpRequest(ctx, "", req)
			// mock server
			mux := http.NewServeMux()
			mux.HandleFunc("/student", StudentNoSignGet)
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, r)
			// get response
			response := w.Result()
			convey.So(response.StatusCode, convey.ShouldEqual, http.StatusOK)
			reader := response.Body
			defer reader.Close()
			bs, err := ioutil.ReadAll(reader)
			convey.So(err, convey.ShouldBeNil)
			println(string(bs))
		})
	})
}

func TestStudentGetHttpInHead(t *testing.T) {
	convey.Convey("TestStudentGetHttpInHead", t, func() {
		convey.Convey("normal test", func() {
			// use sender
			ctx := context.Background()
			req := StudentReq{ID: 1, Name: "xxx"}
			r, _ := StudentGetHttpSignInHead(ctx, "", req)
			// mock server
			mux := http.NewServeMux()
			mux.HandleFunc("/student", StudentGetNoToken)
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, r)
			// get response
			response := w.Result()
			convey.So(response.StatusCode, convey.ShouldEqual, http.StatusOK)
			reader := response.Body
			defer reader.Close()
			bs, err := ioutil.ReadAll(reader)
			convey.So(err, convey.ShouldBeNil)
			println(string(bs))
		})
	})
}

func TestStudentGetHttpInURL(t *testing.T) {
	convey.Convey("TestStudentGetHttpInURL", t, func() {
		convey.Convey("normal test", func() {
			// use sender
			ctx := context.Background()
			req := StudentReq{ID: 1, Name: "xxx"}
			r, _ := StudentGetHttpSignInURL(ctx, "", req)
			// mock server
			mux := http.NewServeMux()
			mux.HandleFunc("/student", StudentGet)
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, r)
			// get response
			response := w.Result()
			convey.So(response.StatusCode, convey.ShouldEqual, http.StatusOK)
			reader := response.Body
			defer reader.Close()
			bs, err := ioutil.ReadAll(reader)
			convey.So(err, convey.ShouldBeNil)
			println(string(bs))
		})
	})
}

func TestStudentGetInvoke(t *testing.T) {
	d := StudentReq{ID: 1, Name: "xxx"}
	convey.Convey("TestSignatureMiddlewareWithGET", t, func() {
		f := gomonkey.ApplyMethodFunc(reflect.TypeOf(http.DefaultClient), "Do", func(req *http.Request) (*http.Response, error) {
			// mock server
			mux := http.NewServeMux()
			mux.HandleFunc("/student", StudentGet)
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, req)
			// get response
			response := w.Result()
			return response, nil
		})
		defer f.Reset()
		convey.Convey("success test", func() {
			// use sender
			ctx := context.Background()
			_, err := StudentGetInvoke(ctx, "http://localhost:8080", d)
			convey.So(err, convey.ShouldBeNil)
		})
	})
}

func TestStudentGetNotokenInvoke(t *testing.T) {
	d := StudentReq{ID: 1, Name: "xxx"}
	convey.Convey("TestSignatureMiddlewareWithGET", t, func() {
		f := gomonkey.ApplyMethodFunc(reflect.TypeOf(http.DefaultClient), "Do", func(req *http.Request) (*http.Response, error) {
			// mock server
			mux := http.NewServeMux()
			mux.HandleFunc("/student", StudentGetNoToken)
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, req)
			// get response
			response := w.Result()
			return response, nil
		})
		defer f.Reset()
		convey.Convey("success test", func() {
			// use sender
			ctx := context.Background()
			_, err := StudentGetNoTokenInvoke(ctx, "http://localhost:8080", d)
			convey.So(err, convey.ShouldBeNil)
		})
	})
}
