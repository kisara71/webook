package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/go-playground/assert/v2"
	"github.com/kisara71/WeBook/webook/internal/web"
	"github.com/kisara71/WeBook/webook/ioc"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestJWT_SMS_SEND_CODE(t *testing.T) {
	rds := ioc.InitRedis()
	testCases := []struct {
		name   string
		before func(t *testing.T)
		after  func(t *testing.T)

		reqBody string

		wantCode int
		wantBody web.Result
	}{
		{
			name: "发送成功",
			reqBody: `
			{
				"phone": "1518668"
			}
			`,
			wantBody: web.Result{
				Code: 0,
				Msg:  "发送成功",
			},
			before: func(t *testing.T) {

			},
			after: func(t *testing.T) {
				rds.Del(context.Background(), "phone_code:user:1518668")
			},
			wantCode: 200,
		},
		{
			name: "发送验证码太频繁",
			before: func(t *testing.T) {
				rds.Set(context.Background(), "phone_code:user:1518668", 1234, time.Minute*15)
			},
			after: func(t *testing.T) {
				rds.Del(context.Background(), "phone_code:user:1518668")
			},
			reqBody: `
			{
				"phone": "1518668"
			}
			`,
			wantBody: web.Result{
				Code: 4,
				Msg:  "验证码发送太频繁",
			},
			wantCode: 200,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			server := ioc.InitWebServer()

			req, err := http.NewRequest(http.MethodPost,
				"/users/login_sms/code/send",
				bytes.NewBuffer([]byte(tc.reqBody)),
			)
			req.Header.Set("Content-Type", "application/json")
			require.NoError(t, err)

			resp := httptest.NewRecorder()

			server.ServeHTTP(resp, req)

			var rspBody web.Result
			err = json.NewDecoder(resp.Body).Decode(&rspBody)
			assert.Equal(t, err, nil)
			assert.Equal(t, resp.Code, tc.wantCode)
			assert.Equal(t, rspBody, tc.wantBody)
			tc.after(t)
		})
	}
}
