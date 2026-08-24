package console

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

func (s *Server) handlePlantsPage(w http.ResponseWriter, r *http.Request) {
	renderPage(w, "冷站控制台", plantsPageHTML)
}

func (s *Server) handlePlantsAPI(w http.ResponseWriter, r *http.Request) {
	type plantView struct {
		ID         string
		Name       string
		State      string
		Pump       bool
		Compressor bool
		LoadKw     float64
		Lead       bool
		Tripped    bool
		Reason     string
	}
	chillers := s.deps.Plant.AllChillers()
	lead := s.deps.Plant.LeadChiller()
	pumps := s.deps.Plant.AllPumps()
	views := make([]plantView, 0, len(chillers))
	for _, c := range chillers {
		view := plantView{
			ID: c.ID, Name: c.Name, State: string(c.State),
			Pump: c.PumpRunning, Compressor: c.CompressorRunning,
			LoadKw: c.LoadKw, Lead: c.ID == lead, Tripped: c.Tripped, Reason: c.TripReason,
		}
		for _, p := range pumps {
			if p.ChillerID == c.ID {
				view.Pump = p.Running
			}
		}
		views = append(views, view)
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"plants":   views,
		"baseline": s.deps.Plant.Baseline(),
		"lead":     lead,
		"load_forecast": s.deps.Plant.LoadForecast(),
		"peak_load":     s.deps.Plant.PeakLoad(),
	})
}

func (s *Server) handlePlantStart(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	report, err := s.deps.Plant.StartChiller(id)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, report)
}

func (s *Server) handlePlantStop(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := s.deps.Plant.StopChiller(id); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"stopped": id})
}

func (s *Server) handleTrendClose(w http.ResponseWriter, r *http.Request) {
	day := r.URL.Query().Get("day")
	if day == "" {
		day = time.Now().Format("2006-01-02")
	}
	total, err := s.deps.Plant.CloseDay(day)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, total)
}

func decodeBody(w http.ResponseWriter, r *http.Request, target any) bool {
	if err := json.NewDecoder(r.Body).Decode(target); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return false
	}
	return true
}
