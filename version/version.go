package version

const Version = "0.1.0"

var GitCommit string

func BuildVersion() string {
	version := Version
	if GitCommit != "" {
		version += " - " + GitCommit
	}

	return version
}
