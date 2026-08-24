package console

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

func (s *Server) handleZonesPage(w http.ResponseWriter, r *http.Request) {
	renderPage(w, "分区控制台", zonesPageHTML)
}

func (s *Server) handleZonesAPI(w http.ResponseWriter, r *http.Request) {
	type zoneView struct {
		ID       string
		Name     string
		Building string
		Setpoint float64
		Supply   float64
		Reheat   bool
		Shed     bool
		Setback  bool
		Label    string
		Score    float64
	}
	supply := s.deps.Plant.Setpoint().Current()
	profiles := map[string]float64{}
	for _, p := range s.deps.Zone.AllProfiles() {
		profiles[p.ZoneID] = p.Setpoint
	}
	views := make([]zoneView, 0)
	for _, z := range s.deps.Ns.AllZones() {
		view := zoneView{
			ID: z.ID, Name: z.Name, Building: z.Building,
			Supply: supply, Setback: s.deps.Zone.InSetback(z.ID),
			Shed: s.deps.Zone.ShedState(z.ID),
		}
		if setpoint, ok := profiles[z.ID]; ok {
			view.Setpoint = setpoint
		}
		if decision, err := s.deps.Zone.ComputeSetpoint(z.ID, supply); err == nil {
			view.Reheat = decision.Reheat
		}
		if comfort, err := s.deps.Zone.ComfortScore(z.ID, supply); err == nil {
			view.Label = comfort.Label
			view.Score = comfort.Score
		}
		views = append(views, view)
	}
	writeJSON(w, http.StatusOK, map[string]any{"zones": views})
}

func (s *Server) handleZoneSetpoint(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req struct {
		Setpoint float64
	}
	if !decodeBody(w, r, &req) {
		return
	}
	if _, err := s.deps.Ns.GetZone(id); err != nil {
		writeJSON(w, http.StatusNotFound, map[string]any{"error": err.Error()})
		return
	}
	profile, err := s.deps.Zone.Profile(id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]any{"error": err.Error()})
		return
	}
	profile.Setpoint = req.Setpoint
	if err := s.deps.Zone.UpsertProfile(profile); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"zone": id, "setpoint": req.Setpoint})
}

func (s *Server) handleZoneSetback(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req struct {
		Active bool
	}
	if !decodeBody(w, r, &req) {
		return
	}
	if _, err := s.deps.Ns.GetZone(id); err != nil {
		writeJSON(w, http.StatusNotFound, map[string]any{"error": err.Error()})
		return
	}
	now := time.Now()
	if req.Active {
		if err := s.deps.Zone.EnterSetback(id, now, 6, 10); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
			return
		}
		s.deps.Plant.Setpoint().SetBaseline(6)
	} else {
		rec, err := s.deps.Zone.ExitSetback(id, now)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
			return
		}
		if err := s.deps.Plant.RecoverBaseline(rec, now); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
			return
		}
	}
	writeJSON(w, http.StatusOK, map[string]any{"zone": id, "setback": req.Active})
}
