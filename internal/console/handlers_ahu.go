package console

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

func (s *Server) handleAhuAPI(w http.ResponseWriter, r *http.Request) {
	type unitView struct {
		ID          string
		Zone        string
		Name        string
		Running     bool
		Prestart    bool
		FreeCooling bool
	}
	views := make([]unitView, 0)
	for _, unit := range s.deps.Ahu.AllUnits() {
		views = append(views, unitView{
			ID: unit.ID, Zone: unit.ZoneID, Name: unit.Name,
			Running: unit.Running, Prestart: unit.Prestart, FreeCooling: unit.FreeCooling,
		})
	}
	writeJSON(w, http.StatusOK, map[string]any{"units": views})
}

func (s *Server) handleAhuStart(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := s.deps.Ahu.StartUnit(id); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"started": id})
}

func (s *Server) handleAhuStop(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := s.deps.Ahu.StopUnit(id); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"stopped": id})
}

func (s *Server) handleAhuPrestart(w http.ResponseWriter, r *http.Request) {
	zoneID := chi.URLParam(r, "zone")
	now := time.Now()
	sch, err := s.deps.Schedule.ForZone(zoneID)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]any{"error": err.Error()})
		return
	}
	decision, err := s.deps.Ahu.EvaluatePrestart(zoneID, sch.Start, now)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}
	if err := s.deps.Ahu.MarkPrestart(zoneID, decision.ShouldStart); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"zone": zoneID, "prestart": decision})
}
