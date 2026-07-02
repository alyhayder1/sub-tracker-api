package subscriptions

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"

	db "github.com/sub-tracker-hq/sub-tracker-api/db/generated"
	"github.com/sub-tracker-hq/sub-tracker-api/internal/auth"
	"github.com/sub-tracker-hq/sub-tracker-api/internal/common"
)

type Handler struct {
	Queries *db.Queries
}

func NewHandler(q *db.Queries) *Handler {
	return &Handler{Queries: q}
}

func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/", h.Create)
	r.Get("/", h.List)
	r.Get("/{id}", h.Get)
	r.Patch("/{id}", h.Update)
	r.Delete("/{id}", h.Delete)
	return r
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	authUser, ok := auth.UserFromContext(r.Context())
	if !ok {
		common.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req createRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		common.WriteError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if req.Name == "" {
		common.WriteError(w, http.StatusBadRequest, "name is required")
		return
	}
	if req.BillingCycle == "" {
		common.WriteError(w, http.StatusBadRequest, "billing_cycle is required")
		return
	}
	if req.Amount == "" {
		common.WriteError(w, http.StatusBadRequest, "amount is required")
		return
	}
	amount, err := common.NewMoney(req.Amount)
	if err != nil {
		common.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	currency := req.Currency
	if currency == "" {
		currency = "EUR"
	}
	if len(currency) != 3 {
		common.WriteError(w, http.StatusBadRequest, "currency must be a 3-letter ISO 4217 code")
		return
	}
	category := req.Category
	if category == "" {
		category = "other"
	}

	sub, err := h.Queries.CreateSubscription(r.Context(), db.CreateSubscriptionParams{
		UserID:          authUser.ID,
		ProviderKey:     req.ProviderKey,
		Name:            req.Name,
		Category:        category,
		Amount:          amount,
		Currency:        currency,
		BillingCycle:    req.BillingCycle,
		NextBillingDate: req.NextBillingDate,
		Status:          "active",
		Source:          "manual",
	})
	if err != nil {
		common.WriteError(w, http.StatusInternalServerError, "failed to create subscription")
		return
	}

	common.WriteJSON(w, http.StatusCreated, toResponse(sub))
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	authUser, ok := auth.UserFromContext(r.Context())
	if !ok {
		common.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	subs, err := h.Queries.ListSubscriptionsByUser(r.Context(), authUser.ID)
	if err != nil {
		common.WriteError(w, http.StatusInternalServerError, "failed to list subscriptions")
		return
	}

	responses := make([]response, len(subs))
	for i, s := range subs {
		responses[i] = toResponse(s)
	}
	common.WriteJSON(w, http.StatusOK, responses)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	authUser, ok := auth.UserFromContext(r.Context())
	if !ok {
		common.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	id, err := common.ParseUUID(chi.URLParam(r, "id"))
	if err != nil {
		common.WriteError(w, http.StatusBadRequest, "invalid id")
		return
	}

	sub, err := h.Queries.GetSubscription(r.Context(), db.GetSubscriptionParams{ID: id, UserID: authUser.ID})
	if errors.Is(err, pgx.ErrNoRows) {
		common.WriteError(w, http.StatusNotFound, "subscription not found")
		return
	}
	if err != nil {
		common.WriteError(w, http.StatusInternalServerError, "failed to get subscription")
		return
	}

	common.WriteJSON(w, http.StatusOK, toResponse(sub))
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	authUser, ok := auth.UserFromContext(r.Context())
	if !ok {
		common.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	id, err := common.ParseUUID(chi.URLParam(r, "id"))
	if err != nil {
		common.WriteError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var req updateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		common.WriteError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if req.Name == "" {
		common.WriteError(w, http.StatusBadRequest, "name is required")
		return
	}
	if req.BillingCycle == "" {
		common.WriteError(w, http.StatusBadRequest, "billing_cycle is required")
		return
	}
	if req.Amount == "" {
		common.WriteError(w, http.StatusBadRequest, "amount is required")
		return
	}
	amount, err := common.NewMoney(req.Amount)
	if err != nil {
		common.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	currency := req.Currency
	if currency == "" {
		currency = "EUR"
	}
	if len(currency) != 3 {
		common.WriteError(w, http.StatusBadRequest, "currency must be a 3-letter ISO 4217 code")
		return
	}
	category := req.Category
	if category == "" {
		category = "other"
	}

	sub, err := h.Queries.UpdateSubscription(r.Context(), db.UpdateSubscriptionParams{
		ID:              id,
		UserID:          authUser.ID,
		Name:            req.Name,
		Category:        category,
		Amount:          amount,
		Currency:        currency,
		BillingCycle:    req.BillingCycle,
		NextBillingDate: req.NextBillingDate,
	})
	if errors.Is(err, pgx.ErrNoRows) {
		common.WriteError(w, http.StatusNotFound, "subscription not found")
		return
	}
	if err != nil {
		common.WriteError(w, http.StatusInternalServerError, "failed to update subscription")
		return
	}

	common.WriteJSON(w, http.StatusOK, toResponse(sub))
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	authUser, ok := auth.UserFromContext(r.Context())
	if !ok {
		common.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	id, err := common.ParseUUID(chi.URLParam(r, "id"))
	if err != nil {
		common.WriteError(w, http.StatusBadRequest, "invalid id")
		return
	}

	rowsAffected, err := h.Queries.DeleteSubscription(r.Context(), db.DeleteSubscriptionParams{ID: id, UserID: authUser.ID})
	if err != nil {
		common.WriteError(w, http.StatusInternalServerError, "failed to delete subscription")
		return
	}
	if rowsAffected == 0 {
		common.WriteError(w, http.StatusNotFound, "subscription not found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
