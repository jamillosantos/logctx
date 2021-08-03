package logctx

import (
	"io"
	"io/ioutil"
	"net/http"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type HttpResponseOptions struct {
	LimitBody int64
	SkipBody  bool
}

type httpResponseMarshaler struct {
	HttpResponseOptions
	response *http.Response
}

func (h httpResponseMarshaler) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	if h.response.Request != nil {
		encoder.AddString("request.uri", h.response.Request.RequestURI)
		encoder.AddString("request.method", h.response.Request.Method)
	} else {
		encoder.AddBool("request", false)
	}
	encoder.AddString("status", h.response.Status)
	encoder.AddInt("status_code", h.response.StatusCode)
	encoder.AddInt64("content_length", h.response.ContentLength)
	if h.SkipBody {
		return nil
	}
	var reader io.Reader = h.response.Body
	if h.LimitBody > 0 {
		reader = &io.LimitedReader{R: h.response.Body, N: h.LimitBody}
	}
	bodyBytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}
	encoder.AddByteString("body", bodyBytes)
	return nil
}

// HttpResponse creates a Object fields that will output the response Status, StatusCode, ContentLength and Body.
func HttpResponse(key string, response *http.Response) zap.Field {
	return zap.Object(key, httpResponseMarshaler{response: response})
}

// HttpResponseWithOpts creates a Object fields that will output the response Status, StatusCode, ContentLength. The body will
// not be part of the output.
func HttpResponseWithOpts(key string, response *http.Response, opts HttpResponseOptions) zap.Field {
	return zap.Object(key, httpResponseMarshaler{response: response, HttpResponseOptions: opts})
}
