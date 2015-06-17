package application

import (
	"encoding/json"
	"net/http"

	"github.com/mholt/binding"
	"github.com/thoas/muxer"
)

func NotFoundHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "404 Not found", http.StatusNotFound)
	})
}

type Handler func(muxer.Response, *Request, *Application)

var ImageHandler Handler = func(res muxer.Response, req *Request, app *Application) {
	debug, err := app.Jq.Bool("debug")

	if err != nil {
		debug = false
	}

	file, err := app.ImageFileFromRequest(req, true, true)

	if err != nil {
		if debug {
			panic(err)
		} else {
			res.Abort(404, "404 Not found")
			return
		}
	}

	res.SetHeaders(file.Headers, true)
	res.ResponseWriter.Write(file.Content())
}

var UploadHandler = func(res muxer.Response, req *http.Request, app *Application) {
	if !app.EnableUpload {
		res.Forbidden()
		return
	}

	if app.SourceStorage == nil {
		res.Abort(500, "Your application doesn't have a source storage")
		return
	}

	var err error

	multipartForm := new(MultipartForm)
	errs := binding.Bind(req, multipartForm)
	if errs.Handle(res) {
		return
	}

	file, err := multipartForm.Upload(app.SourceStorage)

	if err != nil {
		panic(err)
	}

	content, err := json.Marshal(map[string]string{
		"filename": file.Filename(),
		"path":     file.Path(),
		"url":      file.URL(),
	})

	if err != nil {
		panic(err)
	}

	res.ContentType("application/json")
	res.ResponseWriter.Write(content)
}

var GetHandler Handler = func(res muxer.Response, req *Request, app *Application) {
	file, err := app.ImageFileFromRequest(req, false, false)

	if err != nil {
		panic(err)
	}

	content, err := json.Marshal(map[string]string{
		"filename": file.Filename(),
		"path":     file.Path(),
		"url":      file.URL(),
	})

	if err != nil {
		panic(err)
	}

	res.ContentType("application/json")
	res.ResponseWriter.Write(content)
}

var RedirectHandler Handler = func(res muxer.Response, req *Request, app *Application) {
	file, err := app.ImageFileFromRequest(req, false, false)

	if err != nil {
		panic(err)
	}

	res.PermanentRedirect(file.URL())
}
