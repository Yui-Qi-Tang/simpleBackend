package httphandler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"regexp"
	"simpleBackend/log"
	"time"

	"simpleBackend/handlers/maindb/models/nasa"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const (
	nasaAPIHost string = "api.nasa.gov"
	apodAPIPath string = "planetary/apod"
)

// api doc: https://api.nasa.gov/index.html#browseAPI

func apiEndpointBuilder(base, path string, parameters map[string]string) (string, error) {
	urlBase, err := url.Parse(base)
	if err != nil {
		return "", errors.Wrap(err, "failed to parse base of url")
	}

	// Path params
	urlBase.Path += path

	if len(parameters) == 0 {
		return urlBase.String(), nil
	}

	// Query params
	params := url.Values{}
	for k, v := range parameters {
		params.Add(k, v)
	}
	urlBase.RawQuery = params.Encode()
	return urlBase.String(), nil
}

// Apod retrives apod date from Nasa via date
func (h *Handler) Apod(c *gin.Context) {

	// http request to nasa
	parameters := make(map[string]string, 2)
	parameters["api_key"] = h.NasaAPIKey
	date := c.DefaultQuery("date", time.Now().Format("2006-01-02"))

	if len(date) > 0 {
		// check the format of 'date' via regexp
		// YYYY-MM-DD
		var validDate *regexp.Regexp = regexp.MustCompile(`^[0-9]{4}\-[0-1]{1}[0-9]{1}\-[0-3]{1}[0-9]{1}`)
		if !validDate.MatchString(date) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid format of date, the correct is 'YYYY-MM-DD'"})
			return
		}
		parameters["date"] = date
	}
	endpoint, err := apiEndpointBuilder("https://"+nasaAPIHost, apodAPIPath, parameters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// get data from db

	ap := &nasa.Apod{}

	if res := h.MainDB.Take(&ap, nasa.Apod{Date: date}); res.Error != nil {
		if !errors.Is(res.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": res.Error.Error()})
			return
		}
	} else {
		c.JSON(http.StatusOK, ap.Map())
		return
	}

	// send http request
	resp, err := http.Get(endpoint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errors.Wrap(err, "failed to get response from nasa").Error()})
		log.Logger.Debug("error with http request", zap.String("api key", h.NasaAPIKey))
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "status code is not 200 from Nasa"})
		return
	}

	buf := &bytes.Buffer{}
	if _, err := buf.ReadFrom(resp.Body); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errors.Wrap(err, "failed to read data from http response").Error()})
		return
	}

	apod := nasa.Apod{}

	if err := json.Unmarshal(buf.Bytes(), &apod); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errors.Wrap(err, "failed to unmarshal json").Error()})
		return
	}

	c.JSON(http.StatusOK, apod.Map())
	h.MainDB.Create(&apod)
}
