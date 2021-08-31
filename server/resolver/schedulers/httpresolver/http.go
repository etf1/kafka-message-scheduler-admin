package httpresolver

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/etf1/kafka-message-scheduler-admin/server/helper"
	"github.com/etf1/kafka-message-scheduler-admin/server/resolver/schedulers"
	log "github.com/sirupsen/logrus"
)

const (
	SchedulerDefaultPort = "8000"
	DefaultTimeout       = 2 * time.Second
)

var (
	ErrPartialResults = errors.New("resolver: partial results")
	ErrNoResults      = errors.New("resolver: no results")
)

type Scheduler struct {
	HostName  string     `json:"name"`
	HTTPPort  string     `json:"http_port"`
	Instances []Instance `json:"instances"`
}

type Instance struct {
	IP               net.IP   `json:"ip"`
	HostNames        []string `json:"hostname"`
	Topics           []string `json:"topics"`
	HistoryTopic     string   `json:"history_topic"`
	BootstrapServers string   `json:"bootstrap_servers"`
}

func (s Scheduler) Name() string {
	return s.HostName
}

func (s Scheduler) BootstrapServers() string {
	if len(s.Instances) > 0 {
		return s.Instances[0].BootstrapServers
	}
	return ""
}

func (s Scheduler) Topics() []string {
	if len(s.Instances) > 0 {
		return s.Instances[0].Topics
	}
	return nil
}

func (s Scheduler) History() string {
	if len(s.Instances) > 0 {
		return s.Instances[0].HistoryTopic
	}
	return ""
}

func (i Instance) Name() string {
	if len(i.HostNames) > 0 {
		return i.HostNames[0]
	}
	return i.IP.String()
}

type Resolver struct {
	Hosts []string
}

func NewResolver(hosts []string) Resolver {
	return Resolver{
		Hosts: hosts,
	}
}

func GetInfo(host string, timeout time.Duration) (resp *http.Response, err error) {
	return helper.Get(host, "/info", timeout)
}

func (r Resolver) List() ([]schedulers.Scheduler, error) {
	if len(r.Hosts) == 0 {
		return nil, fmt.Errorf("hosts undefined")
	}

	result := []schedulers.Scheduler{}

	for _, shost := range r.Hosts {
		host := shost

		log.Printf("scheduler host=>%v<", shost)

		port := SchedulerDefaultPort
		if strings.Contains(shost, ":") {
			host = strings.Split(shost, ":")[0]
			port = strings.Split(shost, ":")[1]
		}

		// renvoie la liste des ips pour un host donné: ex: google.com => [216.58.215.46 2a00:1450:4007:808::200e]
		log.Printf(">>>>>>> host: >%v<", host)
		ips, err := net.LookupIP(host)
		if err != nil {
			log.Errorf("unable to lookup ip for host %v: %v", host, err)
		}

		sch := Scheduler{
			HostName: host,
			HTTPPort: port,
		}

		log.Infof("ips: %+v", ips)

		for _, ip := range ips {
			log.Printf("ip: %v", ip)

			// keep only v4 ips
			if ip.To4() != nil {
				log.Printf("ip is v4: %v", ip)
				names, err := net.LookupAddr(ip.String())
				if err != nil {
					log.Errorf("unable to lookup addr for ip %v: %v", ip, err)
					continue
				}

				instance := Instance{
					IP:        ip,
					HostNames: names,
				}

				info, err := getKafkaInfo(instance.Name() + ":" + port)
				if err != nil {
					log.Errorf("unable to get kafka info for instance %v: %v", instance, err)
					continue
				}

				log.Printf("received info: %+v", info)

				instance.BootstrapServers = info.BootstrapServers
				instance.Topics = info.Topics
				instance.HistoryTopic = info.HistoryTopic
				sch.Instances = append(sch.Instances, instance)

				log.Printf("instances: %+v", sch.Instances)
			}
		}
		if len(sch.Instances) > 0 {
			result = append(result, sch)
		}
	}

	if len(result) == 0 && len(r.Hosts) > 0 {
		return result, ErrNoResults
	}

	if len(result) != len(r.Hosts) {
		return result, ErrPartialResults
	}

	return result, nil
}

type kafka struct {
	BootstrapServers string   `json:"bootstrap_servers"`
	Topics           []string `json:"topics"`
	HistoryTopic     string   `json:"history_topic"`
}
type info struct {
	Host          string   `json:"hostname"`
	Address       []net.IP `json:"address"`
	ServerAddress string   `json:"server_address"`
	kafka         `json:"kafka"`
}

func getKafkaInfo(host string) (info, error) {
	result := info{}

	err := helper.DecodeJSON(host, "/info", DefaultTimeout, &result)
	if err != nil {
		return result, fmt.Errorf("cannot get info from host %v: %v", host, err)
	}
	return result, nil
}
