package application

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strconv"

	"github.com/facebookgo/grace/gracehttp"
	"github.com/honeybadger-io/honeybadger-go"
)

func Run(path string) error {
	runtime.GOMAXPROCS(runtime.NumCPU())

	app, err := NewFromConfigPath(path)

	if err != nil {
		return err
	}

	n := app.InitRouter()

	server := &http.Server{Addr: fmt.Sprintf(":%s", strconv.Itoa(app.Port())), Handler: n}

	log.Fatal(gracehttp.Serve(server), honeybadger.Handler(n))

	return nil
}
