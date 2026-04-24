package engine

import (
	"fmt"
	"testing"

	"github.com/MathieuMoalic/amumax/src/fsutil"
)

func resetGlobals(t *testing.T) string {
	tmpDir := t.TempDir()
	OutputDir = ""
	InputFile = ""
	InputFileSuffix = ""
	CacheDir = ""
	SkipExists = false
	ForceClean = false
	HideProgresBar = false
	SelfTest = false
	SyncAndLog = false
	fsutil.SetWD(tmpDir)

	return tmpDir
}

func TestOD(t *testing.T) {
	resetGlobals(t)

	t.Run("Panic when not initialized", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("OD() should panic when OutputDir is empty")
			}
		}()
		OD()
	})

	t.Run("Returns trailing slash", func(t *testing.T) {
		OutputDir = "tmp/test"
		res := OD()
		if res != "tmp/test/" {
			t.Errorf("Expected tmp/test/, got %s", res)
		}
	})

	t.Run("Returns as is if already has slash", func(t *testing.T) {
		OutputDir = "tmp/test/"
		res := OD()
		if res != "tmp/test/" {
			t.Errorf("Expected tmp/test/, got %s", res)
		}
	})
}

func TestInitIO(t *testing.T) {
	t.Run("Default output directory", func(t *testing.T) {
		tmpDir := resetGlobals(t)
		mx3 := "test_file.mx3"
		InitIO(mx3, "", "", false, false, false, false, false)

		expected := "test_file.zarr"
		if OutputDir != expected {
			t.Errorf("Expected OutputDir %s, got %s", expected, OutputDir)
		}
		if !fsutil.Exists(expected) {
			t.Errorf("Expected directory %s to be created", expected)
		}
		fmt.Printf("\nTree for directory: %s\n", tmpDir)
		fsutil.PrintTree(tmpDir, "")
	})

	t.Run("Custom output directory", func(t *testing.T) {
		tmpDir := resetGlobals(t)
		mx3 := "test_file.mx3"
		customOD := "custom_out"
		InitIO(mx3, customOD, "", false, false, false, false, false)

		if OutputDir != customOD {
			t.Errorf("Expected OutputDir %s, got %s", customOD, OutputDir)
		}
		if !fsutil.Exists(customOD) {
			t.Errorf("Expected directory %s to be created", customOD)
		}
		fmt.Printf("\nTree for directory: %s\n", tmpDir)
		fsutil.PrintTree(tmpDir, "")
	})

	t.Run("ForceClean wipes directory", func(t *testing.T) {
		tmpDir := resetGlobals(t)
		mx3 := "test_file.mx3"
		od := "clean_me"

		// Create directory and a file inside it
		fsutil.Mkdir(od)
		fsutil.Touch(od + "/somefile")

		fmt.Printf("\nTree before calling InitIO for directory: %s\n", tmpDir)
		fsutil.PrintTree(tmpDir, "")

		InitIO(mx3, od, "", false, true, false, false, false)

		if OutputDir != od {
			t.Errorf("Expected OutputDir %s, got %s", od, OutputDir)
		}
		// Check if the file inside was removed
		if fsutil.Exists(od + "/somefile") {
			t.Errorf("Expected directory %s to be cleaned", od)
		}
		fmt.Printf("\nTree after calling InitIO for directory: %s\n", tmpDir)
		fsutil.PrintTree(tmpDir, "")
	})

	t.Run("Existing directory nests output", func(t *testing.T) {
		tmpDir := resetGlobals(t)
		mx3 := "nest_test.mx3"
		od := "nest_dir"

		fsutil.Mkdir(od)

		InitIO(mx3, od, "", false, false, false, false, false)

		expected := od + "/" + "nest_test.zarr"
		if OutputDir != expected {
			t.Errorf("Expected OutputDir %s, got %s", expected, OutputDir)
		}
		if !fsutil.Exists(expected) {
			t.Errorf("Expected nested directory %s to be created", expected)
		}
		fmt.Printf("\nTree for directory: %s\n", tmpDir)
		fsutil.PrintTree(tmpDir, "")
	})
}
