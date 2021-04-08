package httphandler

import (
	"bytes"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// api doc: https://api.nasa.gov/index.html#browseAPI

// Apod ...
func (h *Handler) Apod(c *gin.Context) {
	resp, err := http.Get("https://api.nasa.gov/planetary/apod?api_key=" + h.NasaAPIKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "status code is not 200 from Nasa"})
		return
	}

	buf := &bytes.Buffer{}

	if _, err := buf.ReadFrom(resp.Body); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": errors.Wrap(err, "failed to read data from http response")})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": buf.String()})

}
