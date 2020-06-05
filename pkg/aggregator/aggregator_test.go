// Copyright 2019 Intel Corporation
// SPDX-License-Identifier: BSD-3-Clause
package aggregator

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

const testJSONdate = "[{\"path\":\"testrepo/simple-test-test\",\"sha\":\"2c755297a2073d7f317440e8429d274b284a9051\",\"example\":{\"name\":\"Simple Test Test\",\"category\":\"Unit Test\",\"categories\":[\"TestCat\"],\"description\":\"I am a simple test\",\"author\":\"Intel Corporation\",\"date\":\"1970-01-01\",\"tag\":\"test\",\"sample_readme_uri\":\"https://test.com\"}}]"

type testNewAggregatorData struct {
	dir           string
	ts            *httptest.Server
	testLanguages []string
}

func setupAggregatorTest(t *testing.T) *testNewAggregatorData {
	t.Helper()
	var td testNewAggregatorData

	//Get Cache to use
	dir, err := ioutil.TempDir("", "example")
	if err != nil {
		t.Error(err)
	}
	td.dir = dir

	// Get a test http server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, testJSONdate)
	}))
	td.ts = ts

	td.testLanguages = []string{"cpp"}

	return &td
}

func (tdata *testNewAggregatorData) removeLock(t *testing.T) {
	t.Helper()
	os.Remove(filepath.Join(tdata.dir, AggregatorLocalAPILevel, cacheLockName))
}

func (tdata *testNewAggregatorData) cleanup() {
	os.RemoveAll(tdata.dir)
	tdata.ts.Close()
}

func TestNewAggregator(t *testing.T) {
	td := setupAggregatorTest(t)
	defer td.cleanup()

	_, err := NewAggregator("", "", []string{}, true, true)
	if err == nil {
		t.Errorf("this NewAggregator setup should have failed! With empty URL")
	}
	td.removeLock(t)

	_, err = NewAggregator("1asd://sd", "", []string{}, true, true)
	if err == nil {
		t.Errorf("this NewAggregator setup should have failed! With malformed URL ")
	}
	td.removeLock(t)

	_, err = NewAggregator("http://abcIShouldNotExist.intel.com/", "", []string{}, true, true)
	if err == nil {
		t.Errorf("this NewAggregator setup should have failed! With empty directory passed")
	}
	td.removeLock(t)

	// I wanted to test failing to create the local cache directory but I could think of a
	// way todo it: a) crossplatform b) running as admin could be valid usecase

	_, err = NewAggregator("http://abcIShouldNotExist.intel.com/", td.dir, []string{}, true, true)
	if err == nil {
		t.Errorf("lenth or nil lanauges array should have failed")
	}
	td.removeLock(t)

	_, err = NewAggregator("http://abcIShouldNotExist.intel.com/", td.dir, td.testLanguages, true, true)
	if err == nil {
		t.Errorf("should not be able to find cpp.json here, err should be network or http related")
	}
	td.removeLock(t)

	//404 Test (non 200)
	badTS := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))

	defer badTS.Close()

	_, err = NewAggregator(badTS.URL, td.dir, td.testLanguages, true, true)
	if err == nil {
		t.Errorf("should have failed with 404 HTTP code")
	}
	td.removeLock(t)

	//Server does not return JSON test
	badJSONTS := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "This is not JSON ¯\\_(ツ)_/¯ ")
	}))
	defer badJSONTS.Close()

	_, err = NewAggregator(badJSONTS.URL, td.dir, td.testLanguages, true, true)
	if err == nil {
		t.Errorf("should have failed garbage JSON")
	}
	td.removeLock(t)

	_, err = NewAggregator(td.ts.URL, td.dir, td.testLanguages, true, true)
	if err != nil {
		t.Error(err)
	}
	td.removeLock(t)

}

func TestGetLocalPath(t *testing.T) {
	td := setupAggregatorTest(t)
	defer td.cleanup()

	a, err := NewAggregator(td.ts.URL, td.dir, td.testLanguages, true, true)
	if err != nil {
		t.Error(err)
	}

	expected := filepath.Join(td.dir, AggregatorLocalAPILevel)
	returned := a.GetLocalPath()

	if returned != expected {
		t.Errorf("directory passed %s was not returned %s as returned by the aggregator", expected, returned)
	}
}

func TestGetURL(t *testing.T) {
	td := setupAggregatorTest(t)
	defer td.cleanup()

	a, err := NewAggregator(td.ts.URL, td.dir, td.testLanguages, true, true)
	if err != nil {
		t.Error(err)
	}

	returned := a.GetURL()

	if a.GetURL() != td.ts.URL {
		t.Errorf("URL passed %s was not returned %s as returned by the aggregator", td.ts.URL, returned)
	}
}

func TestGetLanguages(t *testing.T) {
	td := setupAggregatorTest(t)
	defer td.cleanup()

	a, err := NewAggregator(td.ts.URL, td.dir, td.testLanguages, true, true)
	if err != nil {
		t.Error(err)
	}

	returned := a.GetLanguages()
	if !reflect.DeepEqual(returned, td.testLanguages) { //Maybe should not test order of []string
		t.Errorf("URL passed %s was not returned %s as returned by the aggregator", td.testLanguages, returned)
	}
}

func TestBadTLS(t *testing.T) {
	td := setupAggregatorTest(t)
	defer td.cleanup()
	badTLS := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "I am a baaad guy")
	}))
	defer badTLS.Close()

	_, err := NewAggregator(badTLS.URL, td.dir, td.testLanguages, true, true)
	if err == nil {
		t.Errorf("should have failed due to invalid certificate")
	}
}

func TestOSFiltering(t *testing.T) {
	td := setupAggregatorTest(t)
	defer td.cleanup()
	td.ts.Close()

	//Make a more specifc test server
	const a = "[{\"path\":\"hpc-toolkit-samples-0af2a44aa341bf20ea53d5b908c93d467f65aacf/Nbody\",\"sha\":\"0af2a44aa341bf20ea53d5b908c93d467f65aacf\",\"example\":{\"name\":\"nbody\",\"categories\":[\"Intel\u00AE oneAPI HPC Toolkit/Segment Samples\"],\"description\":\"An N-body simulation is a simulation of a dynamical system of particles, usually under the influence of physical forces, such as gravity. This nbody sample code is implemented using C++ and SYCL language for CPU and GPU.\"}},{\"path\":\"hpc-toolkit-samples-0af2a44aa341bf20ea53d5b908c93d467f65aacf/Particle_Diffusion\",\"sha\":\"0af2a44aa341bf20ea53d5b908c93d467f65aacf\",\"example\":{\"name\":\"Particle-Diffusion\",\"categories\":[\"Intel\u00AE oneAPI HPC Toolkit/Segment Samples\"],\"description\":\"This code sample shows a simple (non-optimized) implementation of a Monte Carlo simulation of the diffusion of water molecules in tissue.\"}},{\"path\":\"hpc-toolkit-samples-0af2a44aa341bf20ea53d5b908c93d467f65aacf/iso3dfd_dpcpp\",\"sha\":\"0af2a44aa341bf20ea53d5b908c93d467f65aacf\",\"example\":{\"name\":\"ISO3DFD\",\"categories\":[\"Intel\u00AE oneAPI HPC Toolkit/Segment Samples\"],\"description\":\"A finite difference stencil kernel for solving 3D acoustic isotropic wave equation\",\"toolchain\":[\"dpcpp\"],\"os\":[\"noknownOS\"],\"sample_readme_uri\":\"https://software.intel.com/en-us/articles/code-samples-for-intel-oneapibeta-toolkits\"}},{\"path\":\"hpc-toolkit-samples-0af2a44aa341bf20ea53d5b908c93d467f65aacf/mandelbrot\",\"sha\":\"0af2a44aa341bf20ea53d5b908c93d467f65aacf\",\"example\":{\"name\":\"Mandelbrot\",\"categories\":[\"Intel\u00AE oneAPI HPC Toolkit/Segment Samples\"],\"description\":\"mandelbrot sample.\",\"os\":[\"noknownOS\"]}}]"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, a)
	}))
	td.ts = ts
	defer td.ts.Close()

	filtered, err := NewAggregator(td.ts.URL, td.dir, td.testLanguages, false, true)
	if err != nil {
		t.Errorf("failed to setup aggregator with good configs")
	}
	td.removeLock(t)

	if len(filtered.Samples[td.testLanguages[0]]) != 2 {
		t.Errorf("aggregator should have only seen two samples %v", len(filtered.Samples[td.testLanguages[0]]))
	}

	unFiltered, err := NewAggregator(td.ts.URL, td.dir, td.testLanguages, true, true)
	if err != nil {
		t.Errorf("failed to setup aggregator with good configs")
	}
	if len(unFiltered.Samples[td.testLanguages[0]]) != 4 {
		t.Errorf("aggregator should have only seen two samples")
	}

}
