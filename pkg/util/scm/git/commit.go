package git

type GitFileContentMap map[string][]byte

type CommitInfo struct {
	CommitMsg    string
	CommitBranch string
	GitFileMap   GitFileContentMap
}

// GitFilePathInfo contains file local path and remote git path
type GitFilePathInfo struct {
	SourcePath      string
	DestinationPath string
}
