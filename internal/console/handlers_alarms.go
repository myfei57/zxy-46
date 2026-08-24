package console

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

func (s *Server) handleAlarmsPage(w http.ResponseWriter, r *http.Request) {
	renderPage(w, "告警台", alarmsPageHTML)
}

func (s *Server) handleAlarmsAPI(w http.ResponseWriter, r *http.Request) {
	type alarmView struct {
		Zone    string
		Message string
		At      time.Time
	}
	alarms := make([]alarmView, 0)
	for _, c := range s.deps.Plant.AllChillers() {
		if c.Tripped {
			alarms = append(alarms, alarmView{Zone: c.Name, Message: c.TripReason, At: time.Now()})
		}
	}
	supply := s.deps.Plant.Setpoint().Current()
	for _, z := range s.deps.Ns.AllZones() {
		comfort, err := s.deps.Zone.ComfortScore(z.ID, supply)
		if err == nil && comfort.Label != "ok" {
			alarms = append(alarms, alarmView{Zone: z.Name, Message: "分区" + comfort.Label, At: time.Now()})
		}
	}
	writeJSON(w, http.StatusOK, map[string]any{"alarms": alarms})
}

func (s *Server) handleAlarmClear(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := s.deps.Plant.ClearTrip(id); err != nil {
		writeJSON(w, http.StatusNotFound, map[string]any{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"cleared": id})
}
