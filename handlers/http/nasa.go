package httphandler

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"time"

	"simpleBackend/handlers/maindb/models/nasa"
	"simpleBackend/log"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// just a test, it will be lead race condition
//var dangerCache map[string]map[string]string = make(map[string]map[string]string)

const (
	nasaAPIHost string = "api.nasa.gov"
	apodAPIPath string = "planetary/apod"
)

// the format of date for Apod: `YYYY-MM-DD`
var validDate *regexp.Regexp = regexp.MustCompile(`^\d{4}\-(0[1-9]|1[012])\-(0[1-9]|[12][0-9]|3[01])$`)

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

	// check if data in cache first
	h.RedisCache.Expire(context.TODO(), date, 1*time.Hour).Result() // set expiration for key

	hgetTimeout, hgetCancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer hgetCancel()
	hGetallRes, err := h.RedisCache.HGetAll(hgetTimeout, date).Result()
	if err != nil && len(hGetallRes) > 0 {
		c.JSON(http.StatusOK, hGetallRes)
		return
	}
	if err != nil {
		// just log
		log.Logger.Info("failed to call 'HGetAll' of Redis", zap.Error(err))
	}

	// check if data in database second
	apod := &nasa.Apod{}
	if res := h.MainDB.Take(&apod, nasa.Apod{Date: date}); res.Error != nil {
		if !errors.Is(res.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": res.Error.Error()})
			return
		}
	} else {
		response := apod.Reponse()
		c.JSON(http.StatusOK, response)

		// cache response
		hsetTimeouit, hsetCancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer hsetCancel()
		if _, err := h.RedisCache.HSet(hsetTimeouit, date, response).Result(); err != nil {
			// just log
			log.Logger.Info("failed to call 'HSET' of Redis server", zap.Error(err))
		}
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

	// check status code
	switch resp.StatusCode {
	case http.StatusOK: // This is blank
	case http.StatusNotFound:
		c.JSON(http.StatusNotFound, gin.H{"error": "data was not found"})
		return
	default:
		buf := &bytes.Buffer{}
		buf.ReadFrom(resp.Body)
		log.Logger.Info("error response", zap.String("content", buf.String()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unknown error:" + buf.String()})
		return

	}

	if err := unmarshalJSONFromIOReader(resp.Body, &apod); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errors.Wrap(err, "failed to read data from http response").Error()})
		return
	}
	c.JSON(http.StatusOK, apod.Reponse())
	h.MainDB.Create(&apod) // save data

	// cache response
	hsetTimeouit, hsetCancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer hsetCancel()
	if _, err := h.RedisCache.HSet(hsetTimeouit, date, apod.Reponse()).Result(); err != nil {
		// just log
		log.Logger.Info("failed to call 'HSET' of Redis server", zap.Error(err))
	}

}
