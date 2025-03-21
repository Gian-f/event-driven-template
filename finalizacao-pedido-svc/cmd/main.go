package main

import (
	"context"
	"encoding/json"
	"finalizacao-pedido-svc/infra/db"
	"finalizacao-pedido-svc/infra/logger"
	"fmt"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func main() {
	// Inicializar componentes
	logger.InitLogger()
	defer logger.Sync()

	// Inicializar OpenTelemetry
	shutdownTracer := logger.InitTracer()
	defer shutdownTracer()

	// Contexto principal com tracing
	ctx, span := logger.Tracer.Start(context.Background(), "main")
	defer span.End()

	// Conectar bancos de dados
	db.ConnectOracle()
	db.ConnectSQLServer()

	// Configurar o handler com otelhttp
	handler := otelhttp.NewHandler(
		logger.Logger(setupRoutes()),
		"",
		otelhttp.WithSpanNameFormatter(func(operation string, r *http.Request) string {
			return fmt.Sprintf("%s %s", r.Method, r.URL.Path)
		}),
	)

	// Configurar servidor HTTP
	server := &http.Server{
		Addr:         ":" + getEnv("PORT", "8080"),
		Handler:      handler,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 90 * time.Second,
		IdleTimeout:  300 * time.Second,
	}

	// Configurar graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Iniciar servidor em goroutine
	go func() {
		logger.Log.Info("Servidor iniciado",
			zap.String("endereço", "http://localhost:"+getEnv("PORT", "8080")),
			zap.String("ambiente", getEnv("APP_ENV", "development")),
		)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.Fatal("Erro ao iniciar servidor", zap.Error(err))
		}
	}()

	// Aguardar sinal de desligamento
	<-stop
	logger.Log.Info("Recebido sinal de desligamento")

	// Desligamento gracioso
	shutdownCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Log.Error("Erro no desligamento do servidor", zap.Error(err))
	} else {
		defer db.CloseConnections()
		logger.Log.Info("Servidor desligado corretamente")
	}
}

func setupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", helloHandler)
	mux.HandleFunc("/health", healthHandler)

	return mux
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	response := map[string]string{
		"status":      "up",
		"version":     getEnv("APP_VERSION", "unversioned"),
		"environment": getEnv("APP_ENV", "development"),
	}
	json.NewEncoder(w).Encode(response)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	response := map[string]string{
		"status": "hello world!",
	}
	json.NewEncoder(w).Encode(response)
}
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
