package restcache

import (
	"time"
)

//type omit *struct{}

type Service struct {
	Id        int       `json:"id"`
	Name      string    `json:"Name"`
	Completed bool      `json:"completed"`
	Due       time.Time `json:"due"`
}

type ServiceCreate struct {
	Name string `json:"Name"`
}

type ServiceCreated struct {
	Service string `json:"Service"`
	Status  string `json:"Status"`
	Desc    string `json:"Desc"`
}

type Services []Service

var services Services
var currentId = 0

// Give us some seed data¬
func init() {
	addService(Service{Name: "go-rnd"})
	addService(Service{Name: "go-keygen"})
}

func addService(s Service) Service {
	currentId += 1
	s.Id = currentId
	services = append(services, s)
	logServiceRegistered(s)
	return s
}

func findService(name string) Service {
	for _, s := range services {
		if s.Name == name {
			return s
		}
	}
	// return empty Todo if not found
	return Service{}
}
