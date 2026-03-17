package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/blueprintautomation/blueprint-roi/internal/analytics"
)

type AnalyticsHandler struct {
	tracker *analytics.Tracker
	log     *zap.Logger
}

func NewAnalyticsHandler(tracker *analytics.Tracker, log *zap.Logger) *AnalyticsHandler {
	return &AnalyticsHandler{tracker: tracker, log: log}
}

func (h *AnalyticsHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, `{"error":"invalid tenant id"}`, http.StatusBadRequest)
		return
	}
	days := 30
	if d := r.URL.Query().Get("days"); d != "" {
		if parsed, err := strconv.Atoi(d); err == nil {
			days = parsed
		}
	}
	stats, err := h.tracker.GetFunnelStats(r.Context(), id, days)
	if err != nil {
		h.log.Error("Failed to get stats", zap.Error(err))
		http.Error(w, `{"error":"stats unavailable"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func (h *AnalyticsHandler) TrackEvent(w http.ResponseWriter, r *http.Request) {
	var event analytics.Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
		return
	}
	h.tracker.Track(r.Context(), &event)
	w.WriteHeader(http.StatusNoContent)
}
