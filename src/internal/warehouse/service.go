package warehouse

import (
	"encoding/json"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
)

type Store interface {
	Get(key string) string
	Set(key, value string) error
	Delete(key string) error
	Join(nodeID, addr string) error
}

type Service struct {
	addr string
	lis  net.Listener

	store Store
}

func New(addr string, store Store) *Service {
	return &Service{
		addr:  addr,
		store: store,
	}
}

func (s *Service) Start() error {
	server := http.Server{
		Handler: s,
	}

	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	s.lis = lis

	http.Handle("/", s)

	go func() {
		if err := server.Serve(s.lis); err != nil {
			log.Fatal(err)
		}
	}()

	return nil
}

func (s *Service) Close() {
	s.lis.Close()
}

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/key") {
		s.handleKeyRequest(w, r)
	} else if r.URL.Path == "/join" {
		s.handleJoinRequest(w, r)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func (s *Service) handleKeyRequest(w http.ResponseWriter, r *http.Request) {
	getKey := func() string {
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) != 3 {
			return ""
		}
		return parts[2]
	}

	switch r.Method {
	case "GET":
		key := getKey()
		if key == "" {
			w.WriteHeader(http.StatusBadRequest)
		}
		value := s.store.Get(key)
		b, err := json.Marshal(map[string]string{key: value})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		io.WriteString(w, string(b))
	case "POST":
		m := map[string]string{}
		if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		for key, value := range m {
			if err := s.store.Set(key, value); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
	case "DELETE":
		key := getKey()
		if key == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if err := s.store.Delete(key); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *Service) handleJoinRequest(w http.ResponseWriter, r *http.Request) {
	m := map[string]string{}
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if len(m) != 2 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	remoteAddr, ok := m["addr"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	nodeID, ok := m["id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := s.store.Join(nodeID, remoteAddr); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *Service) Addr() net.Addr {
	return s.lis.Addr()
}
