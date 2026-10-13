package httpapi

import (
	"log/slog"
	"net/http"

	"github.com/Max-65/libraryInterface/internal/repository"
)

type Handler struct {
	repo          *repository.Repository
	allowedOrigin string
	log           *slog.Logger
}

func NewRouter(repo *repository.Repository, allowedOrigin string, logger *slog.Logger) http.Handler {
	h := &Handler{repo: repo, allowedOrigin: allowedOrigin, log: logger}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health", h.health)

	mux.HandleFunc("GET /api/groups/{code}", h.getGroup)
	mux.HandleFunc("GET /api/groups/{code}/students", h.getGroupStudents)
	mux.HandleFunc("GET /api/groups/{code}/recommendations", h.getGroupRecommendations)

	mux.HandleFunc("GET /api/students/{id}/loans", h.getStudentLoans)

	mux.HandleFunc("GET /api/books/{id}", h.getBook)
	mux.HandleFunc("GET /api/books/{id}/copies", h.getBookCopies)
	mux.HandleFunc("GET /api/books/search", h.searchBooks)

	mux.HandleFunc("POST /api/issuing", h.createIssuing)
	mux.HandleFunc("PATCH /api/issuing/{id}/return", h.returnIssuing)

	return cors(allowedOrigin, mux)
}
