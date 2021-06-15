package restapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/etf1/kafka-message-scheduler-admin/server/db"
	"github.com/etf1/kafka-message-scheduler-admin/server/resolver/schedulers"
	"github.com/etf1/kafka-message-scheduler-admin/server/sort"
	"github.com/etf1/kafka-message-scheduler/schedule"
	"github.com/gorilla/mux"
	"github.com/prometheus/common/log"
)

// coldDB represents schedules stored in a persistent database
// liveDB represents schedules live in the schedulers' instances
func initRouter(coldDB db.DB, liveDB db.DB, resv schedulers.Resolver) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/schedulers", listSchedulers(resv)).Methods(http.MethodGet)
	router.HandleFunc("/scheduler/{name}/schedules", searchSchedules(coldDB)).Methods(http.MethodGet)
	router.HandleFunc("/scheduler/{name}/schedule/{id}", getSchedule(coldDB)).Methods(http.MethodGet)
	router.HandleFunc("/live/scheduler/{name}/schedules", searchSchedules(liveDB)).Methods(http.MethodGet)
	router.HandleFunc("/live/scheduler/{name}/schedule/{id}", getSchedule(liveDB)).Methods(http.MethodGet)
	return router
}

func searchSchedules(d db.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		schedulerName := vars["name"]
		scheduleID := r.URL.Query().Get("schedule-id")
		epochFrom := r.URL.Query().Get("epoch-from")
		epochTo := r.URL.Query().Get("epoch-to")
		sortBy := r.URL.Query().Get("sort-by")
		max := nvl(r.URL.Query().Get("max"), 100)

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

		type result struct {
			Found     int                 `json:"found"`
			Schedules []schedule.Schedule `json:"schedules"`
		}

		found, list, err := d.Search(query)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		w.Write([]byte(fmt.Sprintf("{%q: %d, %q: [", "found", found, "schedules")))

		encoder := json.NewEncoder(w)
		first := true

		for s := range list {
			if !first {
				w.Write([]byte(","))
			}
			if first {
				first = false
			}
			err := encoder.Encode(s)
			if err != nil {
				log.Errorf("unable to encode json: %v", err)
				return
			}

		}

		w.Write([]byte("]}"))
		//w.Write([]byte(fmt.Sprintf("%q: %d}", "found", found)))

		/*
			found, list, err := d.Search(query)
			if err != nil {
				respondWithError(w, http.StatusInternalServerError, err.Error())
				return
			}

			if found == 0 || len(list) == 0 {
				respondWithJSON(w, http.StatusOK, result{
					Found:     0,
					Schedules: []schedule.Schedule{},
				})
				return
			}

			if true {
				encoder := json.NewEncoder(w)

				for i := 0; i < len(list); i++ {
					err := encoder.Encode(list[i])
					if err != nil {
						log.Errorf("unable to encode json: %v", err)
						return
					}
				}
			} else {
				respondWithJSON(w, http.StatusOK, result{
					Found:     found,
					Schedules: list,
				})
			}
		*/
	}
}

func listSchedulers(resv schedulers.Resolver) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		sch, err := resv.List()
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		if len(sch) == 0 {
			respondWithJSON(w, http.StatusOK, []schedulers.Scheduler{})
			return
		}

		respondWithJSON(w, http.StatusOK, sch)
	}
}

func getSchedule(d db.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		sch, err := d.Get(vars["name"], vars["id"])
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		if len(sch) == 0 {
			respondWithJSON(w, http.StatusNotFound, nil)
			return
		}

		respondWithJSON(w, http.StatusOK, sch)
	}
}

func nvl(s string, defaultValue int) int {
	i, err := strconv.ParseInt(s, 10, 0)
	if err != nil {
		return defaultValue
	}
	return int(i)
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

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if payload != nil {
		response, _ := json.Marshal(payload)
		w.Write(response)
	}
}
