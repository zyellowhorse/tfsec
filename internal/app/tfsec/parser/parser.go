package parser

import (
	"fmt"
	"io/fs"
	"strings"

	"github.com/aquasecurity/tfsec/internal/app/tfsec/block"
	"github.com/hashicorp/hcl/v2"

	"github.com/aquasecurity/tfsec/internal/app/tfsec/debug"
	"github.com/aquasecurity/tfsec/internal/app/tfsec/metrics"

	"io/ioutil"
	"os"
	"path/filepath"
)

// Parser is a tool for parsing terraform templates at a given file system location
type Parser struct {
	initialPath    string
	tfvarsPaths    []string
	excludePaths   []string
	stopOnFirstTf  bool
	stopOnHCLError bool
	workspaceName  string
	skipDownloaded bool
}

// New creates a new Parser
func New(initialPath string, options ...Option) *Parser {
	p := &Parser{
		initialPath:   initialPath,
		stopOnFirstTf: true,
		workspaceName: "default",
	}

	for _, option := range options {
		option(p)
	}

	return p
}

func (parser *Parser) parseDirectoryFiles(files []*hcl.File) (block.Blocks, error) {
	var blocks block.Blocks

	for _, file := range files {
		fileBlocks, err := LoadBlocksFromFile(file)
		if err != nil {
			if parser.stopOnHCLError {
				return nil, err
			}
			_, _ = fmt.Fprintf(os.Stderr, "WARNING: HCL error: %s\n", err)
			continue
		}
		if len(fileBlocks) > 0 {
			debug.Log("Added %d blocks from %s...", len(fileBlocks), fileBlocks[0].DefRange.Filename)
		}
		for _, fileBlock := range fileBlocks {
			blocks = append(blocks, block.NewHCLBlock(fileBlock, nil, nil))
		}
	}

	return blocks, nil
}

// ParseDirectory parses all terraform files within a given directory
func (parser *Parser) ParseDirectory() ([]block.Module, error) {

	debug.Log("Finding Terraform subdirectories...")
	t := metrics.Start(metrics.DiskIO)
	subdirectories, err := parser.getSubdirectories(parser.initialPath)
	if err != nil {
		return nil, err
	}
	t.Stop()

	var blocks block.Blocks

	for _, dir := range subdirectories {
		if parser.skipDownloaded && strings.Contains(dir, ".terraform") {
			fmt.Printf("skipping download module file %s\n", dir)
			continue
		}
		debug.Log("Beginning parse for directory '%s'...", dir)
		files, err := LoadDirectory(dir, parser.stopOnHCLError)
		if err != nil {
			return nil, err
		}

		parsedBlocks, err := parser.parseDirectoryFiles(files)
		if err != nil {
			return nil, err
		}

		blocks = append(blocks, parsedBlocks...)
	}

	metrics.Add(metrics.BlocksLoaded, len(blocks))

	if len(blocks) == 0 && parser.stopOnFirstTf {
		return nil, nil
	}

	tfPath := parser.initialPath
	if len(subdirectories) > 0 && parser.stopOnFirstTf {
		tfPath = subdirectories[0]
		debug.Log("Project root set to '%s'...", tfPath)
	}

	debug.Log("Loading TFVars...")

	inputVars, err := LoadTFVars(parser.tfvarsPaths)
	if err != nil {
		return nil, err
	}

	var modulesMetadata *ModulesMetadata
	if parser.skipDownloaded {
		debug.Log("Skipping module metadata loading, --exclude-downloaded-modules passed")
	} else {
		debug.Log("Loading module metadata...")
		t = metrics.Start(metrics.DiskIO)
		modulesMetadata, _ = LoadModuleMetadata(tfPath)
		t.Stop()
	}

	debug.Log("Evaluating expressions...")
	workingDir, _ := os.Getwd()
	evaluator := NewEvaluator(tfPath, tfPath, workingDir, blocks, inputVars, modulesMetadata, nil, parser.stopOnHCLError, parser.workspaceName)
	modules, err := evaluator.EvaluateAll()
	if err != nil {
		return nil, err
	}
	return modules, nil

}

func (parser *Parser) getSubdirectories(path string) ([]string, error) {
	entries, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	entries = parser.RemoveExcluded(path, entries)

	var results []string
	for _, entry := range entries {

		if !entry.IsDir() && (filepath.Ext(entry.Name()) == ".tf" || strings.HasSuffix(entry.Name(), ".tf.json")) {
			debug.Log("Found qualifying subdirectory containing .tf files: %s", path)
			results = append(results, path)
			if parser.stopOnFirstTf {
				return results, nil
			}
		}
	}

	for _, entry := range entries {

		if entry.IsDir() {
			dirs, err := parser.getSubdirectories(filepath.Join(path, entry.Name()))
			if err != nil {
				return nil, err
			}
			results = append(results, dirs...)
		}
	}

	return results, nil
}

func (parser *Parser) RemoveExcluded(path string, entries []fs.FileInfo) (valid []fs.FileInfo) {
	if len(parser.excludePaths) == 0 {
		return entries
	}

	for _, entry := range entries {
		var remove bool
		fullPath := filepath.Join(path, entry.Name())
		for _, excludePath := range parser.excludePaths {
			if fullPath == excludePath {
				remove = true
			}
		}
		if !remove {
			valid = append(valid, entry)
		} else {
			debug.Log("Excluding path %s", fullPath)
		}
	}
	return valid
}
