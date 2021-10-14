package app

import (
	"log"
	"net/http"

	"github.com/flatcar-linux/flog/app/handler"
	"github.com/flatcar-linux/flog/app/model"
	"github.com/flatcar-linux/flog/config"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type App struct {
	Router *mux.Router
	DB     *gorm.DB
}

func (a *App) Initialize(config *config.Config) {
	/* 	dbURI := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True",
		config.DB.Username,
		config.DB.Password,
		config.DB.Host,
		config.DB.Port,
		config.DB.Name,
		config.DB.Charset)

	db, err := gorm.Open(config.DB.Dialect, dbURI) */

	// TODO: Migrate to using the ORM layer in way that it intelligently understands DB Dialect
	// and constructs the dbURI
	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Could not connect database")
	}

	a.DB = model.DBMigrate(db)
	a.Router = mux.NewRouter()
	a.setRouters()
}

// setRouters sets the all required routers
func (a *App) setRouters() {
	// Routing for handling the projects
	a.Get("/projects", a.handleRequest(handler.GetAllProjects))

	//a.Post("/projects", a.handleRequest(handler.CreateProject))
	//a.Get("/projects/{title}", a.handleRequest(handler.GetProject))
	//a.Put("/projects/{title}", a.handleRequest(handler.UpdateProject))
	//a.Delete("/projects/{title}", a.handleRequest(handler.DeleteProject))
	//a.Put("/projects/{title}/archive", a.handleRequest(handler.ArchiveProject))
	//a.Delete("/projects/{title}/archive", a.handleRequest(handler.RestoreProject))
}

// Get wraps the router for GET method
func (a *App) Get(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("GET")
}

// Post wraps the router for POST method
func (a *App) Post(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("POST")
}

// Put wraps the router for PUT method
func (a *App) Put(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("PUT")
}

// Delete wraps the router for DELETE method
func (a *App) Delete(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("DELETE")
}

// Run the app on it's router
func (a *App) Run(host string) {
	log.Fatal(http.ListenAndServe(host, a.Router))
}

type RequestHandlerFunction func(db *gorm.DB, w http.ResponseWriter, r *http.Request)

func (a *App) handleRequest(handler RequestHandlerFunction) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(a.DB, w, r)
	}
}
