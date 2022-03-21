package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"newsaggregator/pkg/storage"
	"strconv"
)

type API struct {
	db storage.Interface
	r  *mux.Router
}

// New returns api object.
func New(db storage.Interface) *API {
	api := API{
		db: db,
	}
	api.r = mux.NewRouter()
	api.endpoints()
	return &api
}

// Registering API handlers.
func (api *API) endpoints() {
	api.r.Use(api.HeadersMiddleware)
	// get n latest news
	api.r.HandleFunc("/news/{n}", api.posts).Methods(http.MethodGet, http.MethodOptions)
	// webapp
	api.r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./webapp"))))
}

// Router returns requests router.
// Required for providing router to web server.
func (api *API) Router() *mux.Router {
	return api.r
}

func (api *API) HeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(w, r)
	})
}

// Getting all posts.
func (api *API) posts(w http.ResponseWriter, r *http.Request) {
	s := mux.Vars(r)["n"]
	n, err := strconv.Atoi(s)
	posts, err := api.db.Posts(int64(n))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	bytes, err := json.Marshal(posts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(bytes)
}
