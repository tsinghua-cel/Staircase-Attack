package versions

import (
	"fmt"
)

var (
	TagVersion string
	AppName    string
	BuildTime  string
	GoVersion  string
	GitBranch  string
	CommitSha1 string
)

func Version() string {
	return fmt.Sprintf("%s-%s-v%s\n", AppName, BuildTime, TagVersion)
}

func DetailVersion() string {
	var info = make([]byte, 0)
	info = append(info, fmt.Sprintf("AppName:\t%s\n", AppName)...)
	info = append(info, fmt.Sprintf("Version:\t%s\n", TagVersion)...)
	info = append(info, fmt.Sprintf("GitBranch:\t%s\n", GitBranch)...)
	info = append(info, fmt.Sprintf("GitCommit:\t%s\n", CommitSha1)...)
	info = append(info, fmt.Sprintf("BuildTime:\t%s\n", BuildTime)...)
	return string(info)
}
