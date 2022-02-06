package utils

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

func TestFindFileInDirAndParents(t *testing.T) {
	const goModFileName = "go.mod"
	wd, err := os.Getwd()
	assert.NoError(t, err)
	projectRoot := filepath.Join(wd, "testdata", "project")

	// Find the file in the current directory
	root, err := FindFileInDirAndParents(projectRoot, goModFileName)
	assert.NoError(t, err)
	assert.Equal(t, projectRoot, root)

	// Find the file in the current directory's parent
	projectSubDirectory := filepath.Join(projectRoot, "dir")
	root, err = FindFileInDirAndParents(projectSubDirectory, goModFileName)
	assert.NoError(t, err)
	assert.Equal(t, projectRoot, root)

	// Look for a file that doesn't exist
	_, err = FindFileInDirAndParents(projectRoot, "notexist")
	assert.Error(t, err)
}

func TestListFilesByFilterFunc(t *testing.T) {
	testDir := filepath.Join("testdata", "listextension")
	expected := []string{filepath.Join(testDir, "a.proj"),
		filepath.Join(testDir, "b.csproj"),
		filepath.Join(testDir, "someproj.csproj")}

	// List files with extension that satisfy the filter function.
	filterFunc := func(filePath string) (bool, error) {
		ext := strings.TrimLeft(filepath.Ext(filePath), ".")
		return regexp.MatchString(`.*proj$`, ext)
	}
	files, err := ListFilesByFilterFunc(testDir, filterFunc)
	if err != nil {
		assert.NoError(t, err)
		return
	}
	assert.ElementsMatch(t, expected, files)
}
