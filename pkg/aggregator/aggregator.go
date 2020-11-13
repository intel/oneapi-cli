// Copyright 2019 Intel Corporation
// SPDX-License-Identifier: BSD-3-Clause

package aggregator

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

type sampleWorkItem struct {
	language string
	s        Sample
	retry    int
}

//Aggregator struct representing the sample store. Not really thread safe
type Aggregator struct {
	baseURL     *url.URL
	localPath   string
	languages   []string
	jobs        chan sampleWorkItem
	results     chan error
	wg          sync.WaitGroup
	sampleCount sync.WaitGroup
	Samples     Samples
	Online      bool
	ignoreOS    bool
	Bulk        bool
}

const defaultRetry = 3

//Samples a map containing an array of avaible samples for that language
// i.e. Samples['cpp"]
type Samples map[string][]Sample

//AggregatorLocalAPILevel the current level of the local file cache. use BaseDir plus this
const AggregatorLocalAPILevel = "v1"

const cacheLockName = "lock"

//ErrCacheLock Is thrown when aggregator's local cache is locked.
var ErrCacheLock = errors.New("aggregator cache is locked")

//HTTPTimeout timeout in seconds for HTTP operations
const HTTPTimeout = 10

//NewAggregator Gives you a Aggregator.
func NewAggregator(URL string, FilePath string, languages []string, ignoreOS bool, bulk bool) (*Aggregator, error) {
	var a Aggregator
	if URL == "" {
		return nil, fmt.Errorf("no sample url passed")
	}
	u, err := checkURL(URL)
	if err != nil {
		return nil, err
	}
	a.baseURL = u

	if FilePath == "" {
		return nil, fmt.Errorf("no base directory passed")
	}

	a.ignoreOS = ignoreOS
	a.Bulk = bulk

	//Add Current file APP level
	a.localPath = filepath.Join(FilePath, AggregatorLocalAPILevel)

	a.languages = languages
	//Create Directory for local path
	if err := os.MkdirAll(a.localPath, 0750); err != nil {
		return nil, err //Package Tests do not cover this
	}

	if a.isLocked() {
		return nil, ErrCacheLock
	}

	if languages == nil || len(languages) < 1 {
		return nil, fmt.Errorf("No Languages are being selected")
	}

	a.Samples = make(map[string][]Sample)

	if err := a.Update(); err != nil {
		return nil, err
	}

	return &a, nil
}

func (a *Aggregator) isLocked() bool {
	return FileExists(filepath.Join(a.localPath, cacheLockName))
}

func (a *Aggregator) lock() {
	_, err := os.Create(filepath.Join(a.localPath, cacheLockName))
	if err != nil {
		log.Fatalf("failed to create cache lock file - %v", err)
	}
}

func sampleWorker(a *Aggregator, jobs <-chan sampleWorkItem, results chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()
	for j := range jobs {
		j.retry--
		err := a.workSample(j) // Try sync the sample
		if err == nil {
			a.sampleCount.Done()
			continue // Sync was good, move on
		}
		if j.retry > 0 {
			a.jobs <- j //Resumbit if retry is
			continue
		}
		a.sampleCount.Done()
		results <- err // Ran out of retries, submit failure to results
	}
}

func (a *Aggregator) workSample(w sampleWorkItem) error {
	_, err := GetTarBall(a.localPath, a.baseURL.String(), w.language, w.s.Path)
	if err != nil {
		return err
	}
	return nil
}

func (a *Aggregator) setupWorkers(n int) {
	a.jobs = make(chan sampleWorkItem, 50)
	a.results = make(chan error, 100)

	for i := 0; i <= n; i++ {
		a.wg.Add(1)
		go sampleWorker(a, a.jobs, a.results, &a.wg)
	}
}

//syncLanguages interates over configured lanauges, and if a newer version is available online

func (a *Aggregator) syncLanguagesIndex() error {
	for _, language := range a.languages {
		localPath := filepath.Join(a.localPath, language+".json")
		update := false
		a.Online = true

		remoteHash, remote, indexErr := sha512URL(a.baseURL.String() + "/" + language + ".json")
		if indexErr != nil {
			log.Print("failed to connect to sample aggregator, attempting to use local cache\n")
			a.Online = false
		}
		if FileExists(localPath) {
			localHash, err := localHash(localPath)
			if err != nil {
				return err
			}
			if !bytes.Equal(remoteHash, localHash) {
				update = true

			}
		} else {
			if !a.Online {
				log.Println("operating offline and local cache does not exist")
				return indexErr
			}
			update = true
		}
		if update && a.Online {
			err := ioutil.WriteFile(localPath, remote, 0644)
			if err != nil {
				return err
			}

		}
		//Ensure Directory for local path of language exists
		if err := os.MkdirAll(filepath.Join(a.localPath, language), 0750); err != nil {
			return err
		}
	}

	return nil
}

//Update updates the local cache
func (a *Aggregator) Update() error {
	err := a.syncLanguagesIndex()
	if err != nil {
		return err
	}

	var outerrors error

	var errWorker sync.WaitGroup
	if a.Bulk {
		//Setup threading pools for downloading all samples
		a.setupWorkers(5) //Start workerpool with 5

		//Start Error collection go routine. Will just return a generic error to
		//this function, will print real errors to fmt. for now
		errWorker.Add(1)
		go func() {
			for e := range a.results {
				if e != nil {
					fmt.Println(e)
					outerrors = fmt.Errorf("error occured on worker")
				}
			}
			if outerrors != nil {
				a.lock() //Poison the cache
			}
			errWorker.Done()
		}()
	}

	for _, language := range a.languages {
		localPath := filepath.Join(a.localPath, language+".json")
		if !FileExists(localPath) {
			return fmt.Errorf("unable to find configured language json (%s)", language)
		}
		languageIndex, err := ioutil.ReadFile(localPath)
		if err != nil {
			return err
		}
		var collected []Sample
		jsonErr := json.Unmarshal(languageIndex, &collected)
		if jsonErr != nil {
			return (jsonErr)
		}

		if !a.ignoreOS {
			collected = filterOnOS(collected)
		}

		if a.Bulk {
			for _, sample := range collected {
				a.sampleCount.Add(1)
				a.jobs <- sampleWorkItem{language, sample, defaultRetry}
			}
		}
		a.Samples[language] = collected

	}

	if a.Bulk {
		a.sampleCount.Wait()
		close(a.jobs)
		a.wg.Wait()      //wait for job channel to be completed.
		close(a.results) //tell error channel workers are done.
		errWorker.Wait() //wait for the erros to be fully processed.
	}

	return outerrors

}

func filterOnOS(c []Sample) (filtered []Sample) {
	for _, s := range c {
		if len(s.Fields.OS) > 0 {
			keep := false
			for _, os := range s.Fields.OS {
				if strings.EqualFold(os, runtime.GOOS) {
					keep = true
				}
			}
			if keep {
				filtered = append(filtered, s)
			}
		} else {
			filtered = append(filtered, s)
		}
	}
	return filtered
}

//GetLocalPath returns the path local path
func (a *Aggregator) GetLocalPath() string {
	return a.localPath
}

//GetURL gets the base URL used for fetching
func (a *Aggregator) GetURL() string {
	return a.baseURL.String()
}

//GetLanguages gets the base URL used for fetching
func (a *Aggregator) GetLanguages() []string {
	return a.languages
}

//GetTarBall Path of the tarball
func GetTarBall(base string, baseURL string, language string, path string) (tar string, err error) {
	tarPath := filepath.Join(base, language, path, language+".tar.gz")

	if FileExists(tarPath) {
		return tarPath, nil
	}
	//Download tarball

	url := baseURL + "/" + path + "/" + language + ".tar.gz"
	if err := downloadFileDirect(tarPath, url); err != nil {
		return "", fmt.Errorf("failed to download sample '%s' - %v", path, err)
	}

	return tarPath, nil
}
