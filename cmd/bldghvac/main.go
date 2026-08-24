package main

import (
	"flag"
	"log"
	"time"

	"bldghvac/internal/ahu"
	"bldghvac/internal/audit"
	"bldghvac/internal/chiller"
	"bldghvac/internal/console"
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

func main() {
	addr := flag.String("addr", ":8080", "listen address")
	data := flag.String("data", "./data", "data directory")
	flag.Parse()

	st := store.New(*data)
	if err := st.EnsureDir("ns", "zone", "plant", "trend", "quota", "audit"); err != nil {
		log.Fatal(err)
	}
	nsSvc := ns.NewService(st)
	zoneSvc := zone.NewService(st)
	schSvc := schedule.NewService(st)
	demandSvc := demand.NewService(st, zoneSvc)
	weatherSvc := weather.NewService(st)
	auditSvc := audit.NewService(st)
	quotaSvc := quota.NewService(st)
	plantSvc := plant.NewService(st, weatherSvc, auditSvc, quotaSvc)
	ahuSvc := ahu.NewService(st, zoneSvc, schSvc, weatherSvc)
	vavSvc := vav.NewService(st, plantSvc.Setpoint())
	svr := console.NewServer(console.Deps{
		Store:    st,
		Ns:       nsSvc,
		Zone:     zoneSvc,
		Schedule: schSvc,
		Demand:   demandSvc,
		Weather:  weatherSvc,
		Audit:    auditSvc,
		Quota:    quotaSvc,
		Plant:    plantSvc,
		Ahu:      ahuSvc,
		Vav:      vavSvc,
	})
	seed(nsSvc, zoneSvc, schSvc, demandSvc, weatherSvc, auditSvc, quotaSvc, plantSvc, ahuSvc, vavSvc)
	log.Printf("bldghvac listening on %s with data dir %s", *addr, *data)
	if err := svr.ListenAndServe(*addr); err != nil {
		log.Fatal(err)
	}
}

func seed(nsSvc *ns.Service, zoneSvc *zone.Service, schSvc *schedule.Service, demandSvc *demand.Service, weatherSvc *weather.Service, auditSvc *audit.Service, quotaSvc *quota.Service, plantSvc *plant.Service, ahuSvc *ahu.Service, vavSvc *vav.Service) {
	now := time.Now()
	_ = nsSvc.AddZone(ns.Zone{ID: "east-1", Name: "东区一层", Building: "east", Floor: 1})
	_ = nsSvc.AddZone(ns.Zone{ID: "east-2", Name: "东区二层", Building: "east", Floor: 2})
	_ = nsSvc.AddZone(ns.Zone{ID: "west-1", Name: "西区一层", Building: "west", Floor: 1})
	lobbyID := ns.NewZoneID()
	_ = nsSvc.AddZone(ns.Zone{ID: lobbyID, Name: "大厅", Building: "east", Floor: 1})
	_ = zoneSvc.UpsertProfile(zone.Profile{ZoneID: "east-1", Setpoint: 22, Deadband: 1, MinDamper: 30})
	_ = zoneSvc.UpsertProfile(zone.Profile{ZoneID: "east-2", Setpoint: 22, Deadband: 1, MinDamper: 30})
	_ = zoneSvc.UpsertProfile(zone.Profile{ZoneID: "west-1", Setpoint: 21, Deadband: 1, MinDamper: 30})
	_ = zoneSvc.UpsertProfile(zone.Profile{ZoneID: lobbyID, Setpoint: 20, Deadband: 1, MinDamper: 30})
	start := schedule.NextWorkdayStart(now, 8, 0)
	_ = schSvc.Add(schedule.Schedule{ID: "s1", ZoneID: "east-1", Start: start})
	_ = schSvc.Add(schedule.Schedule{ID: "s2", ZoneID: "east-2", Start: start})
	_ = schSvc.Add(schedule.Schedule{ID: "s3", ZoneID: "west-1", Start: start})
	_ = demandSvc.SwitchTariff(demand.Tariff{Name: "peak", Baseline: 480})
	_ = weatherSvc.PushReading(weather.Reading{At: now, Temp: 24, Humidity: 55})
	_ = auditSvc.Record(audit.Entry{Source: "system", ZoneID: "east-1", Action: "boot", Detail: "controller started"})
	_ = quotaSvc.RecordMeter("meter-demo", 128.5)
	_ = plantSvc.RegisterChiller(chiller.NewUnit("c1", "1号冷水机组"))
	_ = plantSvc.RegisterChiller(chiller.NewUnit("c2", "2号冷水机组"))
	_ = plantSvc.RegisterPump(chiller.Pump{ID: "p1", ChillerID: "c1"})
	_ = plantSvc.RegisterPump(chiller.Pump{ID: "p2", ChillerID: "c2"})
	_ = ahuSvc.Register(ahu.Unit{ID: "ahu-e1", ZoneID: "east-1", Name: "东一新风机组"})
	_ = ahuSvc.Register(ahu.Unit{ID: "ahu-w1", ZoneID: "west-1", Name: "西一新风机组"})
	_ = vavSvc.RegisterTerminal(vav.Terminal{ID: "vav-e1", ZoneID: "east-1", Name: "东一VAV末端"})
	_ = vavSvc.RegisterTerminal(vav.Terminal{ID: "vav-w1", ZoneID: "west-1", Name: "西一VAV末端"})
	trendStart := time.Date(now.Year(), now.Month(), now.Day(), 9, 0, 0, 0, now.Location())
	_ = plantSvc.RecordTrend(store.TrendPoint{At: trendStart, Kw: 180, Kwh: 2})
	_ = plantSvc.RecordTrend(store.TrendPoint{At: trendStart.Add(30 * time.Minute), Kw: 210, Kwh: 4})
}
