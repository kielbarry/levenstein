package main

import (
	"encoding/json"
	"log"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"
	"io/ioutil"
	"math"
	"github.com/gorilla/mux"
	"github.com/agnivade/levenshtein"
)

type spaHandler struct {
	staticPath string
	indexPath  string
}

type WordAnalysis struct {
	FirstWord string `json:"firstWord"`
	SecondWord  string `json:"secondWord"`
	FirstWordLength float64 `json:"firstWordLength"`
	SecondWordLength float64 `json:"secondWordLength"`
	AbsoluteDifference float64 `json:"absoluteDifference"`
	LevensteinDifference float64 `json:"levensteinDifference"`
	ThirdPartyLevensteinDifference int `json:"thirdPartylevensteinDifference"`
}

func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	path = filepath.Join(h.staticPath, path)
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		http.ServeFile(w, r, filepath.Join(h.staticPath, h.indexPath))
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
}

func getAbsoluteLength(wa WordAnalysis) WordAnalysis {
	wa.FirstWordLength = float64(len(wa.FirstWord))
	wa.SecondWordLength = float64(len(wa.SecondWord))
	wa.AbsoluteDifference = math.Abs(wa.FirstWordLength - wa.SecondWordLength)
	return wa
}

func getLevensteinLength(wa WordAnalysis) WordAnalysis {

	var longer, shorter string
	isSecondLonger := wa.SecondWordLength > wa.FirstWordLength

	switch isSecondLonger {
    case false:
		longer = wa.FirstWord
		shorter = wa.SecondWord
    case true:
		longer = wa.SecondWord
		shorter = wa.FirstWord
	}

	var iterator float64
	iterator = 0
	for pos, char := range shorter {
		charAtOtherStringsIndex := []rune(longer)[pos]
		if (charAtOtherStringsIndex != char) {
			iterator += 1;
		}
	}

	wa.LevensteinDifference = iterator + wa.AbsoluteDifference
	return wa
}



func LevensteinUsingThirdParty(wa WordAnalysis) WordAnalysis {
	wa.ThirdPartyLevensteinDifference = levenshtein.ComputeDistance(wa.FirstWord, wa.SecondWord)
	return wa
}

func analyseWords(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var wordAnalysis WordAnalysis
	json.Unmarshal(reqBody, &wordAnalysis)
	wordAnalysis = getAbsoluteLength(wordAnalysis)
	wordAnalysis = getLevensteinLength(wordAnalysis)
	wordAnalysis = LevensteinUsingThirdParty(wordAnalysis)
	json.NewEncoder(w).Encode(wordAnalysis)
}

func main() {
	router := mux.NewRouter()
	fmt.Printf("\n main running")

	router.HandleFunc("/api/postWords", analyseWords)
	spa := spaHandler{staticPath: "build", indexPath: "asset-management.html"}
	// spa := spaHandler{staticPath: "public", indexPath: "index.html"}
	router.PathPrefix("/").Handler(spa)

	srv := &http.Server{
		Handler: router,
		Addr:    "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}