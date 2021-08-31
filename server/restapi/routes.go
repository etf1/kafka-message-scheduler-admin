package restapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/etf1/kafka-message-scheduler-admin/server/db"
	"github.com/etf1/kafka-message-scheduler-admin/server/resolver/schedulers"
	"github.com/etf1/kafka-message-scheduler-admin/server/sort"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
)

func NewRouter(coldDB, liveDB, historyDB db.DB, resv schedulers.Resolver) http.Handler {
	return cors.AllowAll().Handler(initRouter(coldDB, liveDB, historyDB, resv))
}

// coldDB represents schedules stored in a persistent database
// liveDB represents schedules live in the schedulers' instances
func initRouter(coldDB, liveDB, historyDB db.DB, resv schedulers.Resolver) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/stats", stats(liveDB, coldDB, historyDB, resv)).Methods(http.MethodGet)
	router.HandleFunc("/schedulers", listSchedulers(resv)).Methods(http.MethodGet)
	router.HandleFunc("/scheduler/{name}/schedules", searchSchedules(coldDB)).Methods(http.MethodGet)
	router.HandleFunc("/scheduler/{name}/schedule/{id}", getSchedule(coldDB)).Methods(http.MethodGet)
	router.HandleFunc("/live/scheduler/{name}/schedules", searchSchedules(liveDB)).Methods(http.MethodGet)
	router.HandleFunc("/live/scheduler/{name}/schedule/{id}", getSchedule(liveDB)).Methods(http.MethodGet)
	router.HandleFunc("/history/scheduler/{name}/schedules", searchSchedules(historyDB)).Methods(http.MethodGet)
	router.HandleFunc("/history/scheduler/{name}/schedule/{id}", getSchedule(historyDB)).Methods(http.MethodGet)
	return router
}

// type ResponseSchedule struct {
// 	ID          string `json:"id"`
// 	Epoch       int64  `json:"epoch"`
// 	Timestamp   int64  `json:"timestamp"`
// 	Topic       string `json:"topic"`
// 	TargetTopic string `json:"target-topic"`
// 	TargetKey   string `json:"target-key"`
// }
// type ResponseList struct {
// 	SchedulerName string           `json:"scheduler"`
// 	Schedule      ResponseSchedule `json:"schedule"`
// }

func searchSchedules(d db.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		globalStart := time.Now()

		vars := mux.Vars(r)

		schedulerName := vars["name"]
		scheduleID := r.URL.Query().Get("schedule-id")
		epochFrom := r.URL.Query().Get("epoch-from")
		epochTo := r.URL.Query().Get("epoch-to")
		sortBy := r.URL.Query().Get("sort-by")
		max := max(r.URL.Query().Get("max"))

		query := db.SearchQuery{
			Limit: db.Limit{
				Max: max,
			},
			Filter: db.Filter{
				SchedulerName: schedulerName,
				ScheduleID:    scheduleID,
				EpochRange: db.EpochRange{
					From: epoch(epochFrom),
					To:   epoch(epochTo),
				},
			},
			SortBy: sort.ToSortBy(sortBy),
		}

		found, list, err := d.Search(query)
		if err != nil {
			respondWithError(w, err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err = w.Write([]byte(fmt.Sprintf("{%q: %d, %q: [", "found", found, "schedules")))
		if err != nil {
			log.Errorf("cannot write response start: %v", err)
		}

		encoder := json.NewEncoder(w)
		first := true

		for s := range list {
			start := time.Now()
			if !first {
				_, err = w.Write([]byte(","))
				if err != nil {
					log.Errorf("cannot write response comma separator: %v", err)
				}
			}
			if first {
				first = false
			}
			err = encoder.Encode(s)
			if err != nil {
				log.Errorf("unable to encode json: %v", err)
				return
			}
			log.Warnf("searchSchedules.encode done elapsed=%v", time.Since(start))
		}

		_, err = w.Write([]byte("]}"))
		if err != nil {
			log.Errorf("cannot write response end: %v", err)
		}

		log.Warnf("searchSchedules.all done elapsed=%v", time.Since(globalStart))
	}
}

func stats(liveDB, coldDB, historyDB db.DB, resv schedulers.Resolver) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		schs, err := resv.List()
		if err != nil {
			respondWithError(w, err.Error())
			return
		}

		type stat struct {
			SchedulerName string `json:"scheduler"`
			TotalLive     int    `json:"total_live"`
			History       int    `json:"total_history"`
			Total         int    `json:"total"`
		}

		result := []stat{}
		for _, sch := range schs {
			q := db.SearchQuery{Filter: db.Filter{
				SchedulerName: sch.Name(),
			}}
			totalLive, _, err := liveDB.Search(q)
			if err != nil {
				log.Errorf("stats on live DB failed: %v", err)
			}
			totalHistory, _, err := historyDB.Search(q)
			if err != nil {
				log.Errorf("stats on history DB failed: %v", err)
			}
			total, _, err := coldDB.Search(q)
			if err != nil {
				log.Errorf("stats on live DB failed: %v", err)
			}
			result = append(result, stat{
				SchedulerName: sch.Name(),
				TotalLive:     totalLive,
				History:       totalHistory,
				Total:         total,
			})
		}

		respondWithJSON(w, http.StatusOK, result)
	}
}

func listSchedulers(resv schedulers.Resolver) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		schs, err := resv.List()
		if err != nil {
			respondWithError(w, err.Error())
			return
		}
		respondWithJSON(w, http.StatusOK, schs)
	}
}

func getSchedule(d db.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		sch, err := d.Get(vars["name"], vars["id"])
		if err != nil {
			respondWithError(w, err.Error())
			return
		}

		if len(sch) == 0 {
			respondWithJSON(w, http.StatusNotFound, nil)
			return
		}

		respondWithJSON(w, http.StatusOK, sch)
	}
}

func epoch(s string) int64 {
	if s != "" {
		n, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return 0
		}
		return n
	}
	return 0
}

func max(s string) int {
	if s != "" {
		n, err := strconv.Atoi(s)
		if err != nil {
			return 0
		}
		return n
	}
	return 0
}

func respondWithError(w http.ResponseWriter, message string) {
	respondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if payload != nil {
		response, _ := json.Marshal(payload)
		_, err := w.Write(response)
		if err != nil {
			log.Errorf("cannot respond with JSON payload: %v", err)
		}
	}
}
