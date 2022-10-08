package httpserver

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/illidaris/rest/sender"
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
			mux.HandleFunc("/student", StudentGet)
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, r)
			// get response
			response := w.Result()
			convey.So(response.StatusCode, convey.ShouldEqual, http.StatusOK)
			reader := response.Body
			defer reader.Close()
			_, err := ioutil.ReadAll(reader)
			convey.So(err, convey.ShouldBeError)
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
	req := &StudentGetRequest{
		StudentReq: d,
		Response:   &StudentResponse{},
	}
	o := sender.NewSender(
		sender.WithHeader(sender.HeaderKeyContentType, req.GetContentType()),
	)
	f2 := gomonkey.ApplyMethodFunc(reflect.TypeOf(o), "Invoke", func(ctx context.Context, request sender.IRequest) (interface{}, error) {
		// build newSenderContext
		sc, err := o.NewSenderContext(ctx, request)
		if err != nil {
			return nil, err
		}
		// mock server
		mux := http.NewServeMux()
		mux.HandleFunc("/student", StudentGet)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, sc.Request)
		// get response
		response := w.Result()
		reader := response.Body
		defer reader.Close()
		bs, err := ioutil.ReadAll(reader)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(bs, request.GetResponse())
		return request.GetResponse(), err
	})
	defer f2.Reset()

	convey.Convey("TestSignatureMiddlewareWithGET", t, func() {
		convey.Convey("success test", func() {
			// use sender
			ctx := context.Background()
			r, err := StudentGetInvoke(ctx, "", d)
			convey.So(err, convey.ShouldBeNil)
			println(r.Name)
			// convey.So(responseData.Data.Title, convey.ShouldEqual, "xxx")
			// convey.So(responseData.Data.Body, convey.ShouldEqual, "ccc")
		})
	})
}
