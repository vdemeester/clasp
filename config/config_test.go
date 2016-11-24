package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestFileValidate(t *testing.T) {
	testCases := []struct {
		target        string
		sources       []string
		setupFunc     func(t *testing.T, target string, sources []string, perm os.FileMode) (string, string, func(t *testing.T))
		perm          os.FileMode
		expectedError bool
	}{
		{
			target:        "anything",
			sources:       []string{"afile"},
			setupFunc:     setupNonExistingTargetFolder,
			perm:          0777,
			expectedError: true,
		},
		{
			target:        "anything",
			sources:       []string{"afile"},
			setupFunc:     setupTargetAsAFolder,
			perm:          0777,
			expectedError: true,
		},
		{
			target:        "anything",
			sources:       []string{"afile"},
			setupFunc:     setupTargetAsALink,
			perm:          0777,
			expectedError: true,
		},
		{
			target:        "anything",
			sources:       []string{},
			setupFunc:     setupNonExistingSources,
			perm:          0777,
			expectedError: true,
		},
		{
			target:        "anything",
			sources:       []string{},
			setupFunc:     setupSourcesAsAFile,
			perm:          0777,
			expectedError: true,
		},
		{
			target:        "anything",
			sources:       []string{},
			setupFunc:     setupTemporaryFiles,
			perm:          0777,
			expectedError: false,
		},
		{
			target:        "anything",
			sources:       []string{"one", "two", "three"},
			setupFunc:     setupTemporaryFiles,
			perm:          0777,
			expectedError: false,
		},
	}

	for _, tc := range testCases {
		// FIXME(vdemeester) name the test better
		t.Run(fmt.Sprintf("%s for %v (%v)", tc.target, tc.sources, tc.perm), func(t *testing.T) {
			targetFile, sourceDir, cleanFunc := tc.setupFunc(t, tc.target, tc.sources, tc.perm)
			defer cleanFunc(t)

			file := NewFile(targetFile, sourceDir)
			err := file.Validate()
			if tc.expectedError && err == nil {
				t.Errorf("expected an error, got nothing.")
			}
			if !tc.expectedError && err != nil {
				t.Errorf("expected no error, got %s", err.Error())
			}
		})
	}
}

func setupTargetAsAFolder(t *testing.T, target string, sources []string, perm os.FileMode) (string, string, func(t *testing.T)) {
	targetDir, err := ioutil.TempDir("", "targetDir")
	if err != nil {
		t.Fatal(err)
	}
	return targetDir, "/any/source/dir/really", func(t *testing.T) {
		os.RemoveAll(targetDir)
	}
}

func setupTargetAsALink(t *testing.T, target string, sources []string, perm os.FileMode) (string, string, func(t *testing.T)) {
	targetDir, err := ioutil.TempDir("", "targetDir")
	if err != nil {
		t.Fatal(err)
	}
	targetPath := filepath.Join(targetDir, target)
	linkedTargetPath := filepath.Join(targetDir, "linked")
	err = ioutil.WriteFile(linkedTargetPath, []byte("anything"), perm)
	if err != nil {
		t.Fatal(err)
	}
	err = os.Symlink(linkedTargetPath, targetPath)
	if err != nil {
		t.Fatal(err)
	}
	return targetPath, "/any/source/dir/really", func(t *testing.T) {
		os.RemoveAll(targetDir)
	}
}

func setupSourcesAsAFile(t *testing.T, target string, sources []string, perm os.FileMode) (string, string, func(t *testing.T)) {
	sourceDir, err := ioutil.TempDir("", "sourceDir")
	if err != nil {
		t.Fatal(err)
	}
	sourcePath := filepath.Join(sourceDir, "file")
	err = ioutil.WriteFile(sourcePath, []byte("anything"), perm)
	if err != nil {
		t.Fatal(err)
	}
	return target, sourcePath, func(t *testing.T) {
		err := os.RemoveAll(sourceDir)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func setupNonExistingSources(t *testing.T, target string, sources []string, perm os.FileMode) (string, string, func(t *testing.T)) {
	targetDir, err := ioutil.TempDir("", "targetDir")
	if err != nil {
		t.Fatal(err)
	}
	targetFile := filepath.Join(targetDir, target)
	return targetFile, "/any/source/dir/really", func(t *testing.T) {
		err := os.RemoveAll(targetDir)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func setupNonExistingTargetFolder(t *testing.T, target string, sources []string, perm os.FileMode) (string, string, func(t *testing.T)) {
	sourceDir, err := ioutil.TempDir("", "sourceDir")
	if err != nil {
		t.Fatal(err)
	}

	for index, source := range sources {
		content := fmt.Sprintf("%d-%s", index, source)
		sourcePath := filepath.Join(sourceDir, source)
		err = ioutil.WriteFile(sourcePath, []byte(content), perm)
		if err != nil {
			t.Fatal(err)
		}
	}

	return "/any/target/really", sourceDir, func(t *testing.T) {
		err := os.RemoveAll(sourceDir)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func setupTemporaryFiles(t *testing.T, target string, sources []string, perm os.FileMode) (string, string, func(t *testing.T)) {
	targetDir, err := ioutil.TempDir("", "targetDir")
	if err != nil {
		t.Fatal(err)
	}
	sourceDir, err := ioutil.TempDir("", "sourceDir")
	if err != nil {
		t.Fatal(err)
	}

	for index, source := range sources {
		content := fmt.Sprintf("%d-%s", index, source)
		err = ioutil.WriteFile(source, []byte(content), perm)
		if err != nil {
			t.Fatal(err)
		}
	}

	targetFile := filepath.Join(targetDir, target)

	return targetFile, sourceDir, func(t *testing.T) {
		err := os.RemoveAll(targetDir)
		if err != nil {
			t.Fatal(err)
		}
		err = os.RemoveAll(sourceDir)
		if err != nil {
			t.Fatal(err)
		}
	}
}
