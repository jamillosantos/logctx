//go:generate go run github.com/golang/mock/mockgen -package logctx -destination zap_mock_test.go go.uber.org/zap/zapcore ObjectEncoder

package logctx

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func Test_httpResponseMarshaler_MarshalLogObject(t *testing.T) {
	wantRequestMethod := http.MethodPost
	wantRequestURI := "/request-uri"
	wantHttpStatus := http.StatusInternalServerError
	wantHttpStatusText := http.StatusText(wantHttpStatus)
	wantContentLength := int64(10)
	wantBodyBytes := "1234567890"

	t.Run("success", func(t *testing.T) {
		field := HttpResponse("req1", &http.Response{
			Request: &http.Request{
				Method:     wantRequestMethod,
				RequestURI: wantRequestURI,
			},
			Status:        wantHttpStatusText,
			StatusCode:    wantHttpStatus,
			ContentLength: wantContentLength,
			Body:          io.NopCloser(strings.NewReader(wantBodyBytes)),
		})
		assert.Equal(t, zapcore.ObjectMarshalerType, field.Type)
		require.IsType(t, httpResponseMarshaler{}, field.Interface)
		m := field.Interface.(httpResponseMarshaler)

		ctrl := gomock.NewController(t)
		encoder := NewMockObjectEncoder(ctrl)

		encoder.EXPECT().AddString("request.method", wantRequestMethod)
		encoder.EXPECT().AddString("request.uri", wantRequestURI)
		encoder.EXPECT().AddString("status", wantHttpStatusText)
		encoder.EXPECT().AddInt("status_code", wantHttpStatus)
		encoder.EXPECT().AddInt64("content_length", wantContentLength)
		encoder.EXPECT().AddByteString("body", []byte(wantBodyBytes))

		err := m.MarshalLogObject(encoder)
		require.NoError(t, err)
	})

	t.Run("limit body", func(t *testing.T) {
		givenBodyLimit := int64(4)
		field := HttpResponse("req1", &http.Response{
			Request: &http.Request{
				Method:     wantRequestMethod,
				RequestURI: wantRequestURI,
			},
			Status:        wantHttpStatusText,
			StatusCode:    wantHttpStatus,
			ContentLength: wantContentLength,
			Body:          io.NopCloser(strings.NewReader(wantBodyBytes)),
		})
		assert.Equal(t, zapcore.ObjectMarshalerType, field.Type)
		require.IsType(t, httpResponseMarshaler{}, field.Interface)
		m := field.Interface.(httpResponseMarshaler)
		m.LimitBody = givenBodyLimit

		ctrl := gomock.NewController(t)
		encoder := NewMockObjectEncoder(ctrl)

		encoder.EXPECT().AddString("request.method", wantRequestMethod)
		encoder.EXPECT().AddString("request.uri", wantRequestURI)
		encoder.EXPECT().AddString("status", wantHttpStatusText)
		encoder.EXPECT().AddInt("status_code", wantHttpStatus)
		encoder.EXPECT().AddInt64("content_length", wantContentLength)
		encoder.EXPECT().AddByteString("body", []byte(wantBodyBytes)[:givenBodyLimit])

		err := m.MarshalLogObject(encoder)
		require.NoError(t, err)
	})

	t.Run("with no request", func(t *testing.T) {
		field := HttpResponse("req1", &http.Response{
			Status:        wantHttpStatusText,
			StatusCode:    wantHttpStatus,
			ContentLength: wantContentLength,
			Body:          io.NopCloser(strings.NewReader(wantBodyBytes)),
		})
		assert.Equal(t, zapcore.ObjectMarshalerType, field.Type)
		require.IsType(t, httpResponseMarshaler{}, field.Interface)
		m := field.Interface.(httpResponseMarshaler)

		ctrl := gomock.NewController(t)
		encoder := NewMockObjectEncoder(ctrl)

		encoder.EXPECT().
			AddBool("request", false)
		encoder.EXPECT().AddString("status", wantHttpStatusText)
		encoder.EXPECT().AddInt("status_code", wantHttpStatus)
		encoder.EXPECT().AddInt64("content_length", wantContentLength)
		encoder.EXPECT().AddByteString("body", []byte(wantBodyBytes))

		err := m.MarshalLogObject(encoder)
		require.NoError(t, err)
	})

	t.Run("with no body", func(t *testing.T) {
		field := HttpResponseWithOpts("req1", &http.Response{
			Status:        wantHttpStatusText,
			StatusCode:    wantHttpStatus,
			ContentLength: wantContentLength,
			Body:          io.NopCloser(strings.NewReader(wantBodyBytes)),
		}, HttpResponseOptions{
			SkipBody: true,
		})
		assert.Equal(t, zapcore.ObjectMarshalerType, field.Type)
		require.IsType(t, httpResponseMarshaler{}, field.Interface)
		m := field.Interface.(httpResponseMarshaler)

		ctrl := gomock.NewController(t)
		encoder := NewMockObjectEncoder(ctrl)

		encoder.EXPECT().
			AddBool("request", false)
		encoder.EXPECT().AddString("status", wantHttpStatusText)
		encoder.EXPECT().AddInt("status_code", wantHttpStatus)
		encoder.EXPECT().AddInt64("content_length", wantContentLength)

		err := m.MarshalLogObject(encoder)
		require.NoError(t, err)
	})
}
