package engine

// Management of output directory.

import (
	"os"
	"strings"

	"github.com/MathieuMoalic/amumax/src/fsutil"
	"github.com/MathieuMoalic/amumax/src/log"
	"github.com/MathieuMoalic/amumax/src/zarr"
)

var (
	OutputDir       string // Output directory
	InputFile       string
	InputFileSuffix string

	CacheDir       string
	SkipExists     bool
	ForceClean     bool
	HideProgresBar bool
	SelfTest       bool
	SyncAndLog     bool
)

func OD() string {
	if OutputDir == "" {
		panic("output not yet initialized")
	}
	if !strings.HasSuffix(OutputDir, "/") {
		return OutputDir + "/"
	}
	return OutputDir
}

// InitIO SetOD sets the output directory where auto-saved files will be stored.
// The -o flag can also be used for this purpose.

func InitIO(mx3Path, od, cachedir string, skipexists, forceclean, hideprogressbar, selftest, syncandlog bool) {
	InputFile = mx3Path
	InputFileSuffix = strings.TrimSuffix(mx3Path, ".mx3")
	CacheDir = cachedir
	SkipExists = skipexists
	ForceClean = forceclean
	HideProgresBar = hideprogressbar
	SelfTest = selftest
	SyncAndLog = syncandlog

	if OutputDir != "" {
		panic("output directory already set")
	}

	if od == "" {
		OutputDir = InputFileSuffix + ".zarr"
	} else {
		OutputDir = od
	}
	if fsutil.IsDir(OutputDir) {
		if SkipExists {
			// if directory exists and --skip-exist flag is set, skip the directory
			log.Log.Warn("Directory `%s` exists, skipping because of --skip-exist flag.", OutputDir)
			os.Exit(0)
		} else if ForceClean {
			// if directory exists and --force-clean flag is set, remove the directory
			log.Log.Warn("Cleaning `%s`", OutputDir)
			log.Log.PanicIfError(fsutil.Remove(OutputDir))
			log.Log.PanicIfError(fsutil.Mkdir(OutputDir))
		} else if od != "" {
			// If the -o directory exists and neither --skip-exist nor --force-clean flag is set, put the output InputFile.zarr in the -o directory
			OutputDir = od + "/" + InputFileSuffix + ".zarr"
			log.Log.Warn("Directory `%s` exists, putting output in `%s` because neither --skip-exist nor --force-clean flag is set.", od, OutputDir)
		}
	} else if !fsutil.Exists(OutputDir) {
		// If the -o directory does not exist, create it in the plain -o directory.
		log.Log.PanicIfError(fsutil.Mkdir(OutputDir))
	}
	zarr.InitZgroup("", OD())
}
