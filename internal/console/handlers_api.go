package console

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"bldghvac/internal/schedule"
)

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"status":         "ok",
		"data_dir":       s.deps.Store.Root(),
		"data_dir_ready": s.deps.Store.Exists("ns"),
	})
}

func (s *Server) handleWeatherAPI(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"latest":        s.deps.Weather.Latest(),
		"history":       s.deps.Weather.History(5),
		"forecast":      s.deps.Weather.Forecast(),
		"compensation":  s.deps.Weather.Compensation(),
		"boundary":      s.deps.Weather.Boundary(),
	})
}

func (s *Server) handleWeatherCompensate(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Curve float64
	}
	if !decodeBody(w, r, &req) {
		return
	}
	if err := s.deps.Weather.PushCompensation(req.Curve, time.Now().Add(5*time.Minute)); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"compensation": "pending"})
}

func (s *Server) handleWeatherApply(w http.ResponseWriter, r *http.Request) {
	if err := s.deps.Weather.ApplyCompensation(time.Now()); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}
	if err := s.deps.Plant.SyncWeatherComp(s.deps.Weather.Compensation()); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"state": s.deps.Weather.Compensation().State})
}

func (s *Server) handleWeatherRefresh(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"boundary": s.deps.Weather.RefreshBoundary()})
}

func (s *Server) handleSchedulesAPI(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"schedules":    s.deps.Schedule.List(),
		"workday":      schedule.IsWorkday(time.Now()),
		"day_sequence": s.deps.Schedule.DaySequence(time.Now()),
	})
}

func (s *Server) handleScheduleOffset(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req struct {
		OffsetMinutes int
	}
	if !decodeBody(w, r, &req) {
		return
	}
	sch, err := s.deps.Schedule.ForZone(id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]any{"error": err.Error()})
		return
	}
	if err := s.deps.Schedule.SetOffset(sch.ZoneID, req.OffsetMinutes); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"zone": sch.ZoneID, "offset": req.OffsetMinutes})
}

func (s *Server) handleTrendAPI(w http.ResponseWriter, r *http.Request) {
	day := r.URL.Query().Get("day")
	if day == "" {
		day = time.Now().Format("2006-01-02")
	}
	total, err := s.deps.Quota.ReportTotal(day)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}
	meters, err := s.deps.Quota.MeterReadings(day)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}
	overBudget, err := s.deps.Quota.OverBudget(day, 500)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}
	insight, err := s.deps.Quota.DailyInsight(day)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"day":         day,
		"total":       total,
		"meters":      meters,
		"over_budget": overBudget,
		"insight":     insight,
	})
}

func (s *Server) handleDemandAPI(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"tariff":        s.deps.Demand.Tariff(),
		"baseline":      s.deps.Demand.Baseline(),
		"last_baseline": s.deps.Demand.LastBaseline(),
		"shed":          s.deps.Demand.ShedSummary(),
	})
}

func (s *Server) handleDemandEvaluate(w http.ResponseWriter, r *http.Request) {
	var req struct {
		DemandKw float64
	}
	if !decodeBody(w, r, &req) {
		return
	}
	decision := s.deps.Demand.Evaluate(req.DemandKw)
	for _, z := range s.deps.Ns.AllZones() {
		if err := s.deps.Demand.ApplyLimit(z.ID, decision); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
			return
		}
	}
	writeJSON(w, http.StatusOK, map[string]any{"decision": decision, "shed": s.deps.Demand.ShedSummary()})
}

func (s *Server) handleAuditAPI(w http.ResponseWriter, r *http.Request) {
	entries, err := s.deps.Audit.List(50)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}
	plantEntries, err := s.deps.Audit.BySource("plant")
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"entries":      entries,
		"plant_events": len(plantEntries),
	})
}

func (s *Server) handleAuditClear(w http.ResponseWriter, r *http.Request) {
	if err := s.deps.Audit.Clear(); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"cleared": true})
}

func (s *Server) handleFilesAPI(w http.ResponseWriter, r *http.Request) {
	dir := r.URL.Query().Get("dir")
	if dir == "" {
		dir = "audit"
	}
	files, err := s.deps.Store.List(dir)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"dir": dir, "files": files})
}
