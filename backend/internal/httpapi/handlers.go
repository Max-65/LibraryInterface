package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/Max-65/libraryInterface/internal/models"
	"github.com/Max-65/libraryInterface/internal/repository"

	"github.com/jackc/pgx/v5"
)

func (h *Handler) health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) getGroup(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	g, err := h.repo.GroupByCode(r.Context(), code)
	if err != nil {
		h.log.Error(err.Error())
		writeRepoError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, g)
}

func (h *Handler) getGroupStudents(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	items, err := h.repo.GroupStudents(r.Context(), code)
	if err != nil {
		h.log.Error(err.Error())
		writeRepoError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (h *Handler) getGroupRecommendations(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	items, err := h.repo.GroupRecommendations(r.Context(), code)
	if err != nil {
		h.log.Error(err.Error())
		writeRepoError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (h *Handler) getStudentLoans(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		h.log.Warn("invalid student id")
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid student id"})
		return
	}
	items, err := h.repo.StudentLoans(r.Context(), id)
	if err != nil {
		h.log.Error(err.Error())
		writeRepoError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (h *Handler) searchBooks(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	if q == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "q is required"})
		return
	}
	items, err := h.repo.SearchBooks(r.Context(), q)
	if err != nil {
		h.log.Error(err.Error())
		writeRepoError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (h *Handler) getBook(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		h.log.Warn("invalid book id")
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid book id"})
		return
	}
	item, err := h.repo.BookByID(r.Context(), id)
	if err != nil {
		h.log.Error(err.Error())
		writeRepoError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (h *Handler) getBookCopies(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		h.log.Warn("invalid book id")
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid book id"})
		return
	}
	items, err := h.repo.BookCopies(r.Context(), id)
	if err != nil {
		h.log.Warn("invalid book id")
		writeRepoError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"copy_ids": items})
}

func (h *Handler) createIssuing(w http.ResponseWriter, r *http.Request) {
	var req models.IssueRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}
	if req.StudentID <= 0 || req.CopyID <= 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "student_id and copy_id are required"})
		return
	}

	item, err := h.repo.CreateIssuing(r.Context(), req)
	if err != nil {
		writeRepoError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, item)
}

func (h *Handler) returnIssuing(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid issuing id"})
		return
	}
	item, err := h.repo.ReturnIssuing(r.Context(), id)
	if err != nil {
		writeRepoError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func writeRepoError(w http.ResponseWriter, err error) {
	if errors.Is(err, repository.ErrNotFound) || errors.Is(err, pgx.ErrNoRows) {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
		return
	}
	writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
