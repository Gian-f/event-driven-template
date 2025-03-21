package logger

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

type responseRecorder struct {
	http.ResponseWriter
	body       *bytes.Buffer
	statusCode int
}

func (rec *responseRecorder) Write(b []byte) (int, error) {
	rec.body.Write(b)
	return rec.ResponseWriter.Write(b)
}

func (rec *responseRecorder) WriteHeader(statusCode int) {
	rec.statusCode = statusCode
	rec.ResponseWriter.WriteHeader(statusCode)
}

func InitLogger() {
	logDir := "logs"
	logFilePath := filepath.Join(logDir, "server.log")

	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		if err := os.MkdirAll(logDir, 0755); err != nil {
			panic("Falha ao criar diretório de logs: " + err.Error())
		}
	}

	config := zapcore.EncoderConfig{
		TimeKey:      "timestamp",
		LevelKey:     "level",
		MessageKey:   "message",
		CallerKey:    "caller",
		EncodeTime:   zapcore.ISO8601TimeEncoder,
		EncodeLevel:  zapcore.CapitalLevelEncoder,
		EncodeCaller: zapcore.ShortCallerEncoder,
	}

	fileEncoder := zapcore.NewJSONEncoder(config)
	consoleEncoder := zapcore.NewConsoleEncoder(config)

	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic("Não foi possível criar o arquivo de log: " + err.Error())
	}

	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, zapcore.AddSync(logFile), zapcore.InfoLevel),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel),
	)

	Log = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Capturar headers do request
		requestHeaders := make(map[string]string)
		for k, v := range r.Header {
			requestHeaders[k] = strings.Join(v, ", ")
		}

		// Capturar corpo do request
		var requestBody map[string]interface{}
		if r.Body != nil {
			bodyBytes, _ := io.ReadAll(r.Body)
			if err := json.Unmarshal(bodyBytes, &requestBody); err == nil {
				r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			}
		}

		// Criar recorder para response
		rec := &responseRecorder{
			ResponseWriter: w,
			body:           &bytes.Buffer{},
			statusCode:     http.StatusOK, // Default status
		}

		start := time.Now()

		// Executar o handler principal com contexto atualizado
		next.ServeHTTP(rec, r.WithContext(ctx))

		// Capturar headers do response
		responseHeaders := make(map[string]string)
		for k, v := range rec.Header() {
			responseHeaders[k] = strings.Join(v, ", ")
		}

		// Obter span do contexto
		span := trace.SpanFromContext(ctx)
		traceID := span.SpanContext().TraceID().String()
		spanID := span.SpanContext().SpanID().String()

		// Adicionar atributos ao span
		span.SetAttributes(
			attribute.String("http.method", r.Method),
			attribute.String("http.url", r.URL.String()),
			attribute.String("http.client_ip", r.RemoteAddr),
			attribute.Int("http.status_code", rec.statusCode),
			attribute.String("http.user_agent", r.UserAgent()),
			attribute.String("http.request_headers", serializeHeaders(requestHeaders)),
			attribute.String("http.request_body", serializeBody(requestBody)),
			attribute.String("http.response_headers", serializeHeaders(responseHeaders)),
			attribute.String("http.response_body", rec.body.String()),
		)

		// Logar com todos os detalhes
		Log.Info("Requisição processada",
			zap.String("trace_id", traceID),
			zap.String("span_id", spanID),
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.String("query", r.URL.RawQuery),
			zap.String("ip", r.RemoteAddr),
			zap.Duration("duration", time.Since(start)),
			zap.String("user_agent", r.UserAgent()),
			zap.Int("status_code", rec.statusCode),
			zap.Any("request_headers", requestHeaders),
			zap.Any("request_body", requestBody),
			zap.Any("response_headers", responseHeaders),
			zap.String("response_body", rec.body.String()),
		)
	})
}

func serializeHeaders(headers map[string]string) string {
	jsonData, _ := json.Marshal(headers)
	return string(jsonData)
}

func serializeBody(body map[string]interface{}) string {
	if body == nil {
		return ""
	}
	jsonData, _ := json.Marshal(body)
	return string(jsonData)
}

func Sync() {
	if Log != nil {
		_ = Log.Sync()
	}
}
