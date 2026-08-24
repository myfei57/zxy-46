package console

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"bldghvac/internal/ahu"
	"bldghvac/internal/audit"
	"bldghvac/internal/demand"
	"bldghvac/internal/ns"
	"bldghvac/internal/plant"
	"bldghvac/internal/quota"
	"bldghvac/internal/schedule"
	"bldghvac/internal/store"
	"bldghvac/internal/vav"
	"bldghvac/internal/weather"
	"bldghvac/internal/zone"
)

type Deps struct {
	Store    *store.Store
	Ns       *ns.Service
	Zone     *zone.Service
	Schedule *schedule.Service
	Demand   *demand.Service
	Weather  *weather.Service
	Audit    *audit.Service
	Quota    *quota.Service
	Plant    *plant.Service
	Ahu      *ahu.Service
	Vav      *vav.Service
}

type Server struct {
	deps   Deps
	router chi.Router
}

func NewServer(deps Deps) *Server {
	s := &Server{deps: deps, router: chi.NewRouter()}
	s.routes()
	return s
}

func (s *Server) routes() {
	r := s.router
	r.Use(s.logRequests)
	r.Get("/", s.handleIndex)
	r.Get("/zones", s.handleZonesPage)
	r.Get("/plants", s.handlePlantsPage)
	r.Get("/vav", s.handleVavPage)
	r.Get("/alarms", s.handleAlarmsPage)
	r.Get("/api/health", s.handleHealth)
	r.Get("/api/zones", s.handleZonesAPI)
	r.Get("/api/plants", s.handlePlantsAPI)
	r.Get("/api/vav", s.handleVavAPI)
	r.Get("/api/alarms", s.handleAlarmsAPI)
	r.Get("/api/ahu", s.handleAhuAPI)
	r.Get("/api/weather", s.handleWeatherAPI)
	r.Get("/api/schedules", s.handleSchedulesAPI)
	r.Get("/api/trend", s.handleTrendAPI)
	r.Get("/api/demand", s.handleDemandAPI)
	r.Get("/api/audit", s.handleAuditAPI)
	r.Get("/api/files", s.handleFilesAPI)
	r.Post("/api/plants/start/{id}", s.handlePlantStart)
	r.Post("/api/plants/stop/{id}", s.handlePlantStop)
	r.Post("/api/zones/{id}/setpoint", s.handleZoneSetpoint)
	r.Post("/api/zones/{id}/setback", s.handleZoneSetback)
	r.Post("/api/vav/{id}/cycle", s.handleVavCycle)
	r.Post("/api/vav/{id}/damper", s.handleVavDamper)
	r.Post("/api/ahu/start/{id}", s.handleAhuStart)
	r.Post("/api/ahu/stop/{id}", s.handleAhuStop)
	r.Post("/api/ahu/prestart/{zone}", s.handleAhuPrestart)
	r.Post("/api/weather/compensate", s.handleWeatherCompensate)
	r.Post("/api/weather/apply-comp", s.handleWeatherApply)
	r.Post("/api/weather/refresh", s.handleWeatherRefresh)
	r.Post("/api/demand/evaluate", s.handleDemandEvaluate)
	r.Post("/api/schedules/{id}/offset", s.handleScheduleOffset)
	r.Post("/api/trend/close", s.handleTrendClose)
	r.Post("/api/alarms/clear/{id}", s.handleAlarmClear)
	r.Post("/api/audit/clear", s.handleAuditClear)
}

func (s *Server) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, s.router)
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	renderPage(w, "BldgHVAC", indexPageHTML)
}
