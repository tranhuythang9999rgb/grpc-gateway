package api

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"google.golang.org/grpc/grpclog"
)

type logResponseWriter struct {
	http.ResponseWriter
	statusCode int
	body       bytes.Buffer
}

func (rsp *logResponseWriter) WriteHeader(code int) {
	rsp.statusCode = code
	rsp.ResponseWriter.WriteHeader(code)
}

func (rsp *logResponseWriter) Write(b []byte) (int, error) {
	rsp.body.Write(b) // Ghi lại response body
	return rsp.ResponseWriter.Write(b)
}

func (rsp *logResponseWriter) Unwrap() http.ResponseWriter {
	return rsp.ResponseWriter
}

func newLogResponseWriter(w http.ResponseWriter) *logResponseWriter {
	return &logResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
}

// LogRequestBody logs the request and response details including method, URI, status code, and duration.
func LogRequestBody(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		lw := newLogResponseWriter(w)

		// Đọc và lưu lại request body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to read body: %v", err), http.StatusBadRequest)
			return
		}
		clonedR := r.Clone(r.Context())
		clonedR.Body = io.NopCloser(bytes.NewReader(body))

		// Gọi handler tiếp theo
		h.ServeHTTP(lw, clonedR)

		duration := time.Since(startTime)
		requestSize := len(body)
		responseSize := lw.body.Len()

		// Lấy địa chỉ IP của client
		ip := r.RemoteAddr
		if ip == "" {
			if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
				ip = forwarded
			} else {
				ip = "unknown"
			}
		}

		// Log thông tin chi tiết
		logMessage := fmt.Sprintf("Method: %s, URI: %s, Status: %d, Duration: %v, RequestSize: %dB, ResponseSize: %dB, IP: %s",
			r.Method, r.RequestURI, lw.statusCode, duration, requestSize, responseSize, ip)

		// Log thông tin lỗi nếu status code không phải 2xx
		if lw.statusCode != http.StatusOK {
			grpclog.Errorf("%s, Request Body: %s, Response Body: %s", logMessage, string(body), lw.body.String())
		} else if lw.statusCode == http.StatusOK {
			grpclog.Infof("%s, Request Body: %s, Response Body: %s", logMessage, string(body), lw.body.String())
		}
	})
}
