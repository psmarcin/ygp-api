package video

import (
	"net/http"

	"github.com/sirupsen/logrus"

	"ygp/pkg/utils"

	"github.com/gorilla/mux"
)

var FormatMp4 = "audio/mp4"
var FormatMPEG = "audio/mpeg"

// Handler handles GET /video?videoId= route
func Handler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	videoID := params["videoId"]

	videoURL := GetURL(videoID, FormatMp4)
	if videoURL == "" {
		logrus.Infof("didn't find video (%s) with audio", videoID)
		http.NotFound(w, r)
		return
	}
	w.Header().Del("Content-Type")
	utils.Redirect(w, videoURL)
}

// Handler handles GET /video?videoId= route
func HandlerMP3(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	videoID := params["videoId"]

	videoURL := GetURL(videoID, FormatMPEG)
	if videoURL == "" {
		logrus.Infof("didn't find video (%s) with audio", videoID)
		http.NotFound(w, r)
		return
	}
	w.Header().Del("Content-Type")
	utils.Redirect(w, videoURL)
}

// RedirectHandler is legacy handler for old routes, it redirect to new route
func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	utils.PermanentRedirect(w, "/video/"+r.FormValue("videoId")+".mp3")
}
