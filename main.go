package main

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/ilfey/webmd/internal/components"
	"github.com/ilfey/webmd/internal/fstree"
	"github.com/ilfey/webmd/internal/middlewares"
	"github.com/ilfey/webmd/internal/pages"
	"github.com/rotisserie/eris"
	"github.com/sirupsen/logrus"

	"github.com/kyoto-framework/kyoto/v2"
	"github.com/kyoto-framework/zen/v2"
)

// setupPages registers a
func setupPages(router *mux.Router, root *fstree.Dir, logger *logrus.Logger) {

	indexHandler, err := pages.PDir(root)
	if err != nil {
		panic(eris.Wrap(err, "failed to create index handler"))
	}

	router.HandleFunc("/", kyoto.HandlerPage(indexHandler))

	root.ExecuteOnAllDirs(func(dir *fstree.Dir) error {
		handler, err := pages.PDir(dir)
		if err != nil {
			return err
		}

		route := "/" + dir.Route()

		logger.Infof("register dir: %v", route)

		router.HandleFunc(route, kyoto.HandlerPage(handler))

		return nil
	})

	root.ExecuteOnAllFiles(func(file *fstree.File) error {
		handler, err := pages.PPage(file)
		if err != nil {
			return err
		}

		route := "/" + file.Route()

		logger.Infof("register page: %v", route)

		router.HandleFunc(route, kyoto.HandlerPage(handler))

		return nil
	})
}

func setupActions(router *mux.Router) {
	component := components.CDirMenu(nil)
	pattern := kyoto.ActionConf.Path + kyoto.ComponentName(component) + "/"
	router.PathPrefix(pattern).HandlerFunc(kyoto.HandlerAction(component))
}

func setupMiddlewares(router *mux.Router, logger *logrus.Logger) {
	router.Use(middlewares.Logging(logger))
}

// setupAssets registers a static files handler.
func setupAssets(mux *mux.Router) {
	mux.PathPrefix("/dist/").Handler(
		http.StripPrefix("/dist/", http.FileServer(http.Dir(".dist"))),
	)
}

func main() {
	logger := logrus.New()
	logger.Formatter = &logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "01/02 15:04:05",
	}

	dirEntry, err := os.ReadDir(".md")
	if err != nil {
		logger.Error(err)
	}

	root, err := fstree.NewFromEntries(dirEntry, ".md")
	if err != nil {
		logger.Error(err)
	}

	router := mux.NewRouter()

	// Setup kyoto
	kyoto.TemplateConf.ParseGlob = "templates/*.html"
	kyoto.TemplateConf.FuncMap = kyoto.ComposeFuncMap(
		kyoto.FuncMap, zen.FuncMap,
	)

	setupMiddlewares(router, logger)
	setupAssets(router)
	setupPages(router, root, logger)
	setupActions(router)

	http.ListenAndServe(":8080", router)
}
