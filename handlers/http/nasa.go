package httphandler

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"time"

	"simpleBackend/handlers/maindb/models/nasa"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const (
	nasaAPIHost string = "api.nasa.gov"
	apodAPIPath string = "planetary/apod"
)

// the format of date for Apod: `YYYY-MM-DD`
var validDate *regexp.Regexp = regexp.MustCompile(`^[0-9]{4}\-[0-1]{1}[0-9]{1}\-[0-3]{1}[0-9]{1}`)

// api doc: https://api.nasa.gov/index.html#browseAPI

// apiEndpointBuilder builds api endpoint with query string if any
func apiEndpointBuilder(base, path string, parameters map[string]string) (string, error) {
	urlBase, err := url.Parse(base)
	if err != nil {
		return "", errors.Wrap(err, "failed to parse base of url")
	}

	// Path
	urlBase.Path += path

	// no query string
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

// unmarshalJSONFromIOReader is an util for unmarshal json from io reader to obj
func unmarshalJSONFromIOReader(r io.Reader, obj interface{}) error {
	buf := &bytes.Buffer{}
	if _, err := buf.ReadFrom(r); err != nil {
		return errors.Wrap(err, "failed to read data to buffer")
	}

	if err := json.Unmarshal(buf.Bytes(), obj); err != nil {
		return errors.Wrap(err, "failed to unmarshal json")
	}

	return nil
}

// Apod retrives apod date from Nasa via date
func (h *Handler) Apod(c *gin.Context) {

	// read 'date' from query string and check is whether valid or not
	date := c.DefaultQuery("date", time.Now().Format("2006-01-02"))
	if !validDate.MatchString(date) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid format of date, the correct is 'YYYY-MM-DD'"})
		return
	}

	// check if data in database first
	apod := &nasa.Apod{}
	if res := h.MainDB.Take(&apod, nasa.Apod{Date: date}); res.Error != nil {
		if !errors.Is(res.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": res.Error.Error()})
			return
		}
	} else {
		c.JSON(http.StatusOK, apod.Reponse())
		return
	}

	// if not found in db, so prepare the parameters for nasa api calling
	var parameters map[string]string = map[string]string{
		"api_key": h.NasaAPIKey,
		"date":    date,
	}
	// api endpoint builder
	endpoint, err := apiEndpointBuilder("https://"+nasaAPIHost, apodAPIPath, parameters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// send http request
	resp, err := http.Get(endpoint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errors.Wrap(err, "failed to get response from nasa").Error()})
		return
	}
	defer resp.Body.Close()

	// check response
	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "status code is not 200 from Nasa"})
		return
	}

	if err := unmarshalJSONFromIOReader(resp.Body, &apod); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errors.Wrap(err, "failed to read data from http response").Error()})
		return
	}
	c.JSON(http.StatusOK, apod.Reponse())
	h.MainDB.Create(&apod) // save data
}
