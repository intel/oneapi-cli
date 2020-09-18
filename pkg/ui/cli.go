// Copyright 2019 Intel Corporation
// SPDX-License-Identifier: BSD-3-Clause

package ui

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gdamore/tcell"
	"github.com/intel/oneapi-cli/pkg/aggregator"
	"github.com/intel/oneapi-cli/pkg/browser"
	"github.com/intel/oneapi-cli/pkg/deps"
	"github.com/intel/oneapi-cli/pkg/extractor"
	"gitlab.com/tslocum/cview"
)

const depsMissingEnvFmt = `The sample you have chosen requires the following dependencies: %s

Unfortunately, we are unable to determine if they are present.
Did you use setvars to configure your environment? Are you using a container build environment?`

//CLI data
type CLI struct {
	sidebar    *cview.TextView
	app        *cview.Application
	aggregator *aggregator.Aggregator
	userHome   string
	oneAPIRoot string
	home       *cview.List
	langSelect *cview.List
}

const idzURL = "https://software.intel.com/en-us/oneapi"

func optionViewDocsInBrowser(url string) {
	//If it cant open it for some reason, it will silently fail
	browser.OpenBrowser(url)
}

//NewCLI create a new *CLI element for showing the CLI
func NewCLI(a *aggregator.Aggregator, uH string) (cli *CLI, err error) {
	if a == nil {
		return nil, fmt.Errorf("Aggregator passed not valid")
	}
	if uH == "" {
		return nil, fmt.Errorf("User Home not passed")
	}

	oneRootPath, err := deps.GetOneAPIRoot()
	if err != nil {
		log.Printf("Could not find oneAPI environment, will not check for missing dependencies")
	}
	app := cview.NewApplication()

	if app == nil {
		return nil, fmt.Errorf("Failed to create backend application")
	}
	return &CLI{app: cview.NewApplication(), aggregator: a, userHome: uH, oneAPIRoot: oneRootPath}, nil
}

//Show displays the UI
func (cli *CLI) Show() {

	list := cview.NewList().
		AddItem("Create a project", "", '1', func() {
			cli.selectLang()

		}).ShowSecondaryText(false).
		AddItem("View oneAPI docs in browser", "", '2', func() {
			cli.gotoLinkModel()
		}).
		AddItem("Quit", "Press to exit", 'q', func() {
			cli.app.Stop()
		})

	list.SetBorder(true)

	cli.home = list

	//Main Run action!
	if err := cli.app.SetRoot(cli.home, true).Run(); err != nil {
		log.Fatal(err)
	}
}

func (cli *CLI) gotoLinkModel() {

	optionViewDocsInBrowser(idzURL)
	modal := cview.NewModal().
		SetText("View oneAPI docs:\n" + idzURL + "\n").
		AddButtons([]string{"Back"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			cli.app.SetRoot(cli.home, true)
		})

	cli.app.SetRoot(modal, true)
}

func (cli *CLI) goBackPrj() {
	if cli.langSelect != nil {
		cli.app.SetRoot(cli.langSelect, true)
		return
	}
	cli.app.SetRoot(cli.home, true)
}

func (cli *CLI) selectLang() {
	if len(cli.aggregator.GetLanguages()) == 1 {
		cli.selectProject(cli.aggregator.GetLanguages()[0])
		return
	}

	list := cview.NewList().ShowSecondaryText(false)

	var start rune
	start = '1'

	for _, k := range cli.aggregator.GetLanguages() {
		list.AddItem(k, "", start, func() {
			i := list.GetCurrentItem() //List doesnt support a reference
			cli.selectProject(cli.aggregator.GetLanguages()[i])
		})
		start++
	}
	list.AddItem("Back", "", 'b', func() {
		cli.app.SetRoot(cli.home, true)
	})
	list.AddItem("Quit", "Press to exit", 'q', func() {
		cli.app.Stop()
	})
	list.SetBorder(true).SetTitle("Select sample language")

	cli.langSelect = list

	cli.app.SetRoot(cli.langSelect, true)
}

func (cli *CLI) selectProject(language string) {
	cli.sidebar = cview.NewTextView().SetWordWrap(true).
		SetChangedFunc(func() {
			cli.app.Draw()
		})
	cli.sidebar.SetDynamicColors(true)

	cli.sidebar.Box.SetBorder(true).SetTitle("Description")

	inst := cview.NewTextView()
	inst.SetBorder(true)
	inst.SetText("Press Backspace to return to previous screen!")

	flex := cview.NewFlex().
		AddItem(cli.tree(language), 0, 1, true).
		AddItem(cview.NewFlex().SetDirection(cview.FlexRow).
			AddItem(cli.sidebar, 0, 9, false).
			AddItem(inst, 3, 0, false), 0, 1, false)

	cli.app.SetRoot(flex, true)
}

func newSampleNode(s aggregator.Sample) *cview.TreeNode {
	node := cview.NewTreeNode(s.Fields.Name).SetSelectable(true)
	node.SetReference(s)
	return node
}

//This take a sample and a category to add it to. TODO revist this
func categoriesSeach(parent *cview.TreeNode, cats []string, s aggregator.Sample) {
	if len(cats) == 0 {
		parent.AddChild(newSampleNode(s))
		return
	}
	for _, csearch := range parent.GetChildren() {
		if csearch.GetText() == cats[0] {
			if len(cats) == 1 {
				csearch.AddChild(newSampleNode(s))
				return
				//continue
			}
			categoriesSeach(csearch, cats[1:], s) //recurse removing current from array-pop
			return
		}
	}
	if len(cats) > 0 {
		csearch := cview.NewTreeNode(cats[0]).SetColor(tcell.ColorOrange)
		parent.AddChild(csearch)
		categoriesSeach(csearch, cats[1:], s)
		return
	}
}

type byNodeText []*cview.TreeNode

func (s byNodeText) Len() int {
	return len(s)
}
func (s byNodeText) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s byNodeText) Less(i, j int) bool {
	return strings.ToLower(s[i].GetText()) < strings.ToLower(s[j].GetText())
}

func (cli *CLI) tree(language string) cview.Primitive {

	root := cview.NewTreeNode("Samples").
		SetColor(tcell.ColorOrange)
	tree := cview.NewTreeView().
		SetRoot(root).
		SetCurrentNode(root)

	var missingCat []*cview.TreeNode

	for _, s := range cli.aggregator.Samples[language] {
		if len(s.Fields.Categories) == 0 {
			missingCat = append(missingCat, newSampleNode(s))
			continue //skip samples without any categories
		}

		//root.AddChild(node)
		for _, c := range s.Fields.Categories {
			cats := strings.Split(c, "/")
			categoriesSeach(root, cats, s)
		}

	}

	if len(missingCat) > 0 {
		other := cview.NewTreeNode("Other").SetColor(tcell.ColorOrange)
		root.AddChild(other)
		for _, missingNode := range missingCat {
			other.AddChild(missingNode)
		}
	}

	root.Walk(func(node *cview.TreeNode, parent *cview.TreeNode) bool {
		if len(node.GetChildren()) > 1 {
			sort.Sort(byNodeText(node.GetChildren()))
		}
		return true
	})

	tree.SetChangedFunc(func(node *cview.TreeNode) {
		reference := node.GetReference()
		a, ok := reference.(aggregator.Sample)
		if !ok {

			cli.sidebar.Clear()
			return
		}
		var sideTextExtra string
		if len(a.Fields.Dependencies) > 0 {
			if cli.oneAPIRoot == "" {
				sideTextExtra = fmt.Sprintf(depsMissingEnvFmt, a.Fields.Dependencies)
			} else {
				sideTextExtra, _ = deps.CheckDeps(a.Fields.Dependencies, cli.oneAPIRoot)
			}
		}
		sideTextExtra = cview.Escape(sideTextExtra)

		newText := fmt.Sprintf("%s\n\n[red]%s", a.Fields.Description, sideTextExtra)
		cli.sidebar.SetText(newText)

	}).SetSelectedFunc(func(node *cview.TreeNode) {
		reference := node.GetReference()
		a, ok := reference.(aggregator.Sample)
		if !ok {
			return
		}
		cli.askPath(a, language, "")
	})
	tree.Box.SetBorder(true).SetTitle("Samples")
	tree.SetTopLevel(1)

	tree.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

		if event.Key() == tcell.KeyBackspace2 ||
			event.Key() == tcell.KeyBackspace {
			cli.goBackPrj()
		}

		return event

	})

	return tree
}

func isPathEmpty(path string) bool {
	if !aggregator.FileExists(path) {
		return true
	}
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return true
	}
	return (len(files) == 0)
}

func (cli *CLI) askPath(sample aggregator.Sample, language string, path string) {
	if path == "" {
		pwd, err := os.Getwd()
		if err != nil {
			log.Println(err)
		}
		path = filepath.Join(pwd, filepath.Base(sample.Path))
	}

	text := cview.NewTextView().SetWordWrap(true).
		SetChangedFunc(func() {
			cli.app.Draw()
		}).SetDynamicColors(true)
	form := cview.NewForm().
		AddInputField("Destination", path, 55, nil, func(t string) {
			path = t
		}).
		AddButton("Create", func() {
			path, err := cli.calcPath(path)
			if err != nil {
				return
			}
			if !isPathEmpty(path) {
				cli.confirmOverwrite(sample, language, path)
				return
			}
			outPath, err := cli.createProject(sample, language, path)

			if err != nil {
				cli.app.Stop()
				log.Fatal(err)
			}
			cli.successModal(outPath)

		}).AddButton("Back", func() {
		cli.selectProject(language)
	})

	text.SetBorderPadding(0, 0, 1, 1)

	flex := cview.NewFlex().SetDirection(cview.FlexRow).
		AddItem(form, 7, 0, true)

	flex.SetBorder(true).SetTitle("Create Project").SetTitleAlign(cview.AlignLeft)

	cli.app.SetRoot(flex, true)
}

func (cli *CLI) confirmOverwrite(sample aggregator.Sample, language string, path string) {

	text := fmt.Sprintf("Path %s is not empty, Creating the sample may overwrite some files.", path)
	buttons := []string{"Back", "Confirm"}

	modal := cview.NewModal().
		SetText(text).
		AddButtons(buttons).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {

			if buttonLabel != "Back" {
				outPath, err := cli.createProject(sample, language, path)

				if err != nil {
					cli.app.Stop()
					log.Fatal(err)
				}
				cli.successModal(outPath)
				return
			}
			cli.askPath(sample, language, path)
		})
	cli.app.SetRoot(modal, true)

}

func (cli *CLI) successModal(path string) {

	text := fmt.Sprintf("Sucessfully created project in %s", path)
	buttons := []string{"Quit"}
	printReadmeText := "View Readme and Quit"

	//check for readme,
	readme, err := ioutil.ReadFile(filepath.Join(path, "README.md"))
	if err == nil {
		buttons = append(buttons, printReadmeText)
	} else {
		readme, err = ioutil.ReadFile(filepath.Join(path, "readme.md"))
		if err == nil {
			buttons = append(buttons, printReadmeText)
		}
	}

	modal := cview.NewModal().
		SetText(text).
		AddButtons(buttons).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			cli.app.Stop()
			if buttonLabel != "Quit" {

				//The folloiwing is to reset the terminal
				s, e := tcell.NewScreen()
				if e != nil {
					fmt.Fprintf(os.Stderr, "%v\n", e)
					os.Exit(1)
				}
				if e = s.Init(); e != nil {
					fmt.Fprintf(os.Stderr, "%v\n", e)
					os.Exit(1)
				}
				s.Clear()
				s.Fini()
				fmt.Printf("\n")
				fmt.Printf("%s\n", string(readme))
			}
		})
	cli.app.SetRoot(modal, true)
}

func (cli *CLI) calcPath(projectPath string) (string, error) {
	//Expand env vars the user might have passed through.
	projectPath = os.ExpandEnv(projectPath)
	//Check if tilda ~ is being used, then use the home dir
	if len(projectPath) > 0 && projectPath[0] == '~' {
		//userHome comes from the frontend cli but check the HOME is not empty just incase
		projectPath = filepath.Join(cli.userHome, projectPath[1:]) //Prepend home path and trim tilda
	}

	projectPath, err := filepath.Abs(projectPath)
	if err != nil {
		return "", err
	}
	return projectPath, nil
}

func (cli *CLI) createProject(selectedSample aggregator.Sample, lang string, projectPath string) (output string, err error) {
	//Maybe here we might check if the tarball does not exists and the trigger the aggregator to atempt an update

	tarPath, err := aggregator.GetTarBall(cli.aggregator.GetLocalPath(), cli.aggregator.GetURL(), lang, selectedSample.Path)
	if err != nil {
		return "", err
	}

	err = extractor.ExtractTarGz(tarPath, projectPath)
	if err != nil {
		return "", err
	}
	return projectPath, nil
}
