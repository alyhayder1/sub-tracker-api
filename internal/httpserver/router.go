package httpserver

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"

	db "github.com/sub-tracker-hq/sub-tracker-api/db/generated"
	"github.com/sub-tracker-hq/sub-tracker-api/internal/auth"
	"github.com/sub-tracker-hq/sub-tracker-api/internal/common"
	"github.com/sub-tracker-hq/sub-tracker-api/internal/subscriptions"
)

func NewRouter(queries *db.Queries, pool *pgxpool.Pool, log *slog.Logger) chi.Router {
	r := chi.NewRouter()
	r.Use(chimiddleware.RequestID)
	r.Use(requestLogger(log))
	r.Use(chimiddleware.Recoverer)

	r.Get("/healthz", healthCheck(pool))

	r.Route("/", func(r chi.Router) {
		r.Use(auth.Middleware(queries))
		r.Mount("/subscriptions", subscriptions.NewHandler(queries).Routes())
	})

	return r
}

func healthCheck(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := pool.Ping(r.Context()); err != nil {
			common.WriteJSON(w, http.StatusServiceUnavailable, map[string]string{"status": "unavailable"})
			return
		}
		common.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	}
}

func requestLogger(log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := chimiddleware.NewWrapResponseWriter(w, r.ProtoMajor)
			next.ServeHTTP(ww, r)
			log.Info("request",
				"method", r.Method,
				"path", r.URL.Path,
				"status", ww.Status(),
				"duration_ms", time.Since(start).Milliseconds(),
				"request_id", chimiddleware.GetReqID(r.Context()),
			)
		})
	}
}
