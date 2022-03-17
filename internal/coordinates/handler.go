package coordinates

import (
	"encoding/json"
	"math"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/zfirdavs/devim-lab/internal/currency"
	"github.com/zfirdavs/devim-lab/pkg/cast"
	"github.com/zfirdavs/devim-lab/pkg/errors"
)

type handler struct {
}

func NewHandler() handler {
	return handler{}
}

func (h handler) GetCoordinates() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		values := r.URL.Query()

		w.Header().Set("Content-Type", "application/json")

		// parse x from query parameters
		x, err := cast.ToFloat64("x", values.Get("x"))
		if err != nil {
			errors.Response(w, err, http.StatusBadRequest)
			return
		}

		// parse y from query parameters
		y, err := cast.ToFloat64("y", values.Get("y"))
		if err != nil {
			// fmt.Println(err)
			errors.Response(w, err, http.StatusBadRequest)
			return
		}

		// load env parameters from file
		err = godotenv.Load()
		if err != nil {
			errors.Response(w, err, http.StatusInternalServerError)
			return
		}

		// load x0 parameter from env file
		x0, err := cast.ToFloat64("X0", os.Getenv("X0"))
		if err != nil {
			errors.Response(w, err, http.StatusBadRequest)
			return
		}

		// load y0 parameter from env file
		y0, err := cast.ToFloat64("Y0", os.Getenv("Y0"))
		if err != nil {
			errors.Response(w, err, http.StatusBadRequest)
			return
		}

		// load radius parameter from env file
		radius, err := cast.ToFloat64("R", os.Getenv("R"))
		if err != nil {
			errors.Response(w, err, http.StatusBadRequest)
			return
		}

		// calculate formula
		formula := math.Pow(x-x0, 2)+math.Pow(y-y0, 2) <= math.Pow(radius, 2)
		if formula {
			response, err := currency.GetResponse("USD")
			if err != nil {
				errors.Response(w, err, http.StatusInternalServerError)
				return
			}

			err = json.NewEncoder(w).Encode(response)
			if err != nil {
				errors.Response(w, err, http.StatusInternalServerError)
				return
			}
			return
		}

		response, err := currency.GetResponse("EUR")
		if err != nil {
			errors.Response(w, err, http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			errors.Response(w, err, http.StatusInternalServerError)
			return
		}
	}
}
