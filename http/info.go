package http

import (
	"encoding/json"
	"fmt"
	iiifconfig "github.com/go-iiif/go-iiif/v4/config"
	iiifdriver "github.com/go-iiif/go-iiif/v4/driver"
	iiiflevel "github.com/go-iiif/go-iiif/v4/level"
	iiifprofile "github.com/go-iiif/go-iiif/v4/profile"
	iiifservice "github.com/go-iiif/go-iiif/v4/service"
	gohttp "net/http"
)

func InfoHandler(config *iiifconfig.Config, driver iiifdriver.Driver) (gohttp.HandlerFunc, error) {

	fn := func(w gohttp.ResponseWriter, r *gohttp.Request) {

		ctx := r.Context()

		parser, err := NewIIIFQueryParser(r)

		if err != nil {
			gohttp.Error(w, err.Error(), gohttp.StatusBadRequest)
			return
		}

		id, err := parser.GetIIIFParameter("identifier")

		if err != nil {
			gohttp.Error(w, err.Error(), gohttp.StatusBadRequest)
			return
		}

		image, err := driver.NewImageFromConfig(config, id)

		if err != nil {
			gohttp.Error(w, err.Error(), gohttp.StatusInternalServerError)
			return
		}

		endpoint := EndpointFromRequest(r)

		level, err := iiiflevel.NewLevelFromConfig(config, endpoint)

		if err != nil {
			gohttp.Error(w, err.Error(), gohttp.StatusInternalServerError)
			return
		}

		profile, err := iiifprofile.NewProfile(endpoint, image, level)

		if err != nil {
			gohttp.Error(w, err.Error(), gohttp.StatusInternalServerError)
			return
		}

		for _, service_name := range config.Profile.Services.Enable {

			service_uri := fmt.Sprintf("%s://", service_name)
			service, err := iiifservice.NewService(ctx, service_uri, config, image)

			if err != nil {
				gohttp.Error(w, err.Error(), gohttp.StatusInternalServerError)
				return
			}

			profile.AddService(service)
		}

		b, err := json.Marshal(profile)

		if err != nil {
			gohttp.Error(w, err.Error(), gohttp.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Write(b)

	}

	h := gohttp.HandlerFunc(fn)
	return h, nil
}
