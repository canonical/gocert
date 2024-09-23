package server

import (
	"net/http"

	"github.com/canonical/notary/internal/metrics"
)

// NewHandler takes in a config struct, passes it along to any handlers that will need
// access to it, and takes an http.Handler that will be used to handle metrics.
// then builds and returns it for a server to consume
func NewHandler(config *HandlerConfig) http.Handler {
	apiV1Router := http.NewServeMux()
	apiV1Router.HandleFunc("GET /certificate_requests", adminOrUser(config.JWTSecret, ListCertificateRequests(config)))
	apiV1Router.HandleFunc("POST /certificate_requests", adminOrUser(config.JWTSecret, CreateCertificateRequest(config)))
	apiV1Router.HandleFunc("GET /certificate_requests/{id}", adminOrUser(config.JWTSecret, GetCertificateRequest(config)))
	apiV1Router.HandleFunc("DELETE /certificate_requests/{id}", adminOrUser(config.JWTSecret, DeleteCertificateRequest(config)))
	apiV1Router.HandleFunc("POST /certificate_requests/{id}/certificate", adminOrUser(config.JWTSecret, CreateCertificate(config)))
	apiV1Router.HandleFunc("POST /certificate_requests/{id}/certificate/reject", adminOrUser(config.JWTSecret, RejectCertificate(config)))
	apiV1Router.HandleFunc("DELETE /certificate_requests/{id}/certificate", adminOrUser(config.JWTSecret, DeleteCertificate(config)))

	apiV1Router.HandleFunc("GET /accounts", adminOnly(config.JWTSecret, ListAccounts(config)))
	apiV1Router.HandleFunc("POST /accounts", adminOrFirstUser(config.JWTSecret, config.DB, CreateAccount(config)))
	apiV1Router.HandleFunc("GET /accounts/{id}", adminOrMe(config.JWTSecret, GetAccount(config)))
	apiV1Router.HandleFunc("DELETE /accounts/{id}", adminOnly(config.JWTSecret, DeleteAccount(config)))
	apiV1Router.HandleFunc("POST /accounts/{id}/change_password", adminOrMe(config.JWTSecret, ChangeAccountPassword(config)))

	m := metrics.NewMetricsSubsystem(config.DB)
	frontendHandler := newFrontendFileServer()
	ctx := middlewareContext{
		jwtSecret: config.JWTSecret,
	}
	apiMiddlewareStack := createMiddlewareStack(
		metricsMiddleware(m),
		loggingMiddleware(&ctx),
	)
	metricsMiddlewareStack := createMiddlewareStack(
		metricsMiddleware(m),
	)

	router := http.NewServeMux()
	router.HandleFunc("POST /login", Login(config))
	router.HandleFunc("GET /status", GetStatus(config))
	router.Handle("/metrics", m.Handler)
	router.Handle("/api/v1/", http.StripPrefix("/api/v1", apiMiddlewareStack(apiV1Router)))
	router.Handle("/", metricsMiddlewareStack(frontendHandler))

	return router
}
