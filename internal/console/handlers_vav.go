package console

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"bldghvac/internal/ahu"
)

func (s *Server) handleVavPage(w http.ResponseWriter, r *http.Request) {
	renderPage(w, "VAV 末端", vavPageHTML)
}

func (s *Server) handleVavAPI(w http.ResponseWriter, r *http.Request) {
	type vavView struct {
		ID     string
		Zone   string
		Name   string
		Damper float64
		Reheat bool
		Flow   float64
		Mode   string
	}
	views := make([]vavView, 0)
	for _, term := range s.deps.Vav.AllTerminals() {
		verdict := s.deps.Ahu.FreeCoolingVerdict(term.ZoneID)
		views = append(views, vavView{
			ID: term.ID, Zone: term.ZoneID, Name: term.Name,
			Damper: term.Damper, Reheat: term.ReheatOpen, Flow: term.Flow,
			Mode: verdictReason(verdict, s.deps.Ahu.VentMode(term.ZoneID)),
		})
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"vav":      views,
		"baseline": s.deps.Vav.Baseline(),
	})
}

func verdictReason(verdict ahu.FreeCoolingVerdict, mode string) string {
	if verdict.Allowed {
		return "free-cooling"
	}
	return mode
}

func (s *Server) handleVavCycle(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	term, err := s.deps.Vav.Terminal(id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]any{"error": err.Error()})
		return
	}
	supply := s.deps.Plant.Setpoint().Current()
	decision, err := s.deps.Zone.ComputeSetpoint(term.ZoneID, supply)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}
	sequence, err := s.deps.Vav.ExecuteSequence(term.ZoneID, decision)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"sequence": sequence})
}

func (s *Server) handleVavDamper(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req struct {
		Damper float64
	}
	if !decodeBody(w, r, &req) {
		return
	}
	term, err := s.deps.Vav.Terminal(id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]any{"error": err.Error()})
		return
	}
	if err := s.deps.Vav.SetDamper(term.ZoneID, req.Damper); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"terminal": id, "damper": req.Damper})
}
