package version

import (
	"runtime"
	"runtime/debug"
)

var (
	Version   = "dev"
	BuildDate = ""
	Commit    = ""
)

func init() {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return
	}

	if Version == "dev" && info.Main.Version != "" {
		Version = info.Main.Version
	}

	for _, s := range info.Settings {
		switch s.Key {
		case "vcs.revision":
			Commit = s.Value
		case "vcs.time":
			BuildDate = s.Value
		}
	}
}

type Info struct {
	Version   string
	Platform  string
	GoVersion string
	Commit    string
	BuildDate string
}

func Get() Info {
	return Info{
		Version:   Version,
		Platform:  runtime.GOOS + "/" + runtime.GOARCH,
		GoVersion: runtime.Version(),
		Commit:    Commit,
		BuildDate: BuildDate,
	}
}

func Platform() string {
	return runtime.GOOS + "/" + runtime.GOARCH
}

func GoVersion() string {
	return runtime.Version()
}
