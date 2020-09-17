// Copyright 2019 Intel Corporation
// SPDX-License-Identifier: BSD-3-Clause
package ui

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/gdamore/tcell"
	"github.com/intel/oneapi-cli/pkg/aggregator"
)

func mkTestScreen(t *testing.T, charset string) tcell.SimulationScreen {
	t.Helper()
	s := tcell.NewSimulationScreen(charset)
	if s == nil {
		t.Fatalf("Failed to get simulation screen")
	}
	if e := s.Init(); e != nil {
		t.Fatalf("Failed to initialize screen: %v", e)
	}
	return s
}

type testNewAggregatorData struct {
	cache         string
	home          string
	b             string
	ts            *httptest.Server
	testLanguages []string
	aggregator    *aggregator.Aggregator
}

func setupAggregatorTest(t *testing.T) *testNewAggregatorData {
	t.Helper()
	var td testNewAggregatorData

	//Get Cache to use
	dir, err := ioutil.TempDir("", "example")
	if err != nil {
		t.Error(err)
	}
	td.b = dir
	td.cache = filepath.Join(dir, "cache")
	td.home = filepath.Join(dir, "home")

	// Get a test http server
	fs := http.FileServer(http.Dir("testdata"))
	ts := httptest.NewServer(fs)
	td.ts = ts

	td.testLanguages = []string{"cpp", "python"}

	td.aggregator, err = aggregator.NewAggregator(td.ts.URL, td.cache, td.testLanguages, true, true)
	if err != nil {
		t.Error(err)
	}

	return &td
}

func (tdata *testNewAggregatorData) cleanup() {
	os.RemoveAll(tdata.b)
	tdata.ts.Close()
}

func TestA(t *testing.T) {
	td := setupAggregatorTest(t)
	defer td.cleanup()

	cli, err := NewCLI(td.aggregator, td.home)
	if err != nil {
		t.Error(err)
	}

	s := mkTestScreen(t, "")
	cli.app = cli.app.SetScreen(s)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		cli.Show()
	}()

	s.InjectKey(tcell.KeyRune, 'q', tcell.ModNone)

	wg.Wait()
}

func TestB(t *testing.T) {
	td := setupAggregatorTest(t)
	defer td.cleanup()
	cli, err := NewCLI(td.aggregator, td.home)
	if err != nil {
		t.Error(err)
	}
	s := mkTestScreen(t, "")

	cli.app = cli.app.SetScreen(s)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		cli.Show()
	}()

	s.InjectKey(tcell.KeyRune, 'q', tcell.ModNone)

	wg.Wait()
}

func TestFullFlow(t *testing.T) {
	td := setupAggregatorTest(t)
	defer td.cleanup()
	cli, err := NewCLI(td.aggregator, td.home)
	if err != nil {
		t.Error(err)
	}
	s := mkTestScreen(t, "")

	cli.app = cli.app.SetScreen(s)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		cli.Show() //Let the CLI run in another routine
	}()

	s.InjectKey(tcell.KeyRune, '1', tcell.ModNone)  //Main Screen
	s.InjectKey(tcell.KeyRune, '1', tcell.ModNone)  //Langauge Select
	s.InjectKey(tcell.KeyDown, 'd', tcell.ModNone)  //Get the zebra
	s.InjectKey(tcell.KeyEnter, 'd', tcell.ModNone) //Enter the sample

	ws, err := ioutil.TempDir("", "ws")
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(ws)

	//Remove default input with CTRL+U
	s.InjectKey(tcell.KeyCtrlU, 'd', tcell.ModNone)

	for i := 0; i < len(ws); i++ {
		s.InjectKey(tcell.KeyRune, rune(ws[i]), tcell.ModNone)
		time.Sleep(10 * time.Millisecond) //Need to give time to the key presses :/
	}

	s.InjectKey(tcell.KeyTAB, 'd', tcell.ModNone)   //Tab
	s.InjectKey(tcell.KeyEnter, 'd', tcell.ModNone) //Create
	s.InjectKey(tcell.KeyEnter, 'd', tcell.ModNone) //Dismiss success dialog

	wg.Wait()

	//Test for known file in "sample"

	knownZebra := filepath.Join(ws, "this-is-a-zebra.md")

	if !fileExists(t, knownZebra) {
		t.Errorf("sample creation flow failed! could not find %s", knownZebra)
	}

}

func fileExists(t *testing.T, path string) bool {
	t.Helper()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}
