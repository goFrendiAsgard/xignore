package xignore

import (
	"os"
	"sort"

	"github.com/spf13/afero"
)

// Matcher xignore matcher
type Matcher struct {
	fs afero.Fs
}

// NewSystemMatcher create matcher for system filesystem
func NewSystemMatcher() *Matcher {
	return &Matcher{afero.NewReadOnlyFs(afero.NewOsFs())}
}

// Matches returns matched files from dir files.
func (m *Matcher) Matches(basedir string, options *MatchesOptions) (*MatchesResult, error) {
	vfs := afero.NewBasePathFs(m.fs, basedir)
	ignorefile := options.Ignorefile
	if ok, err := afero.DirExists(vfs, "/"); !ok || err != nil {
		if err == nil {
			return nil, os.ErrNotExist
		}
		return nil, err
	}

	// Root filemap
	rootMap := stateMap{}
	files, errorFiles := collectFiles(vfs)
	// Init all files state
	rootMap.mergeFiles(files, false)

	// Apply before patterns
	beforePatterns, err := makePatterns(options.BeforePatterns)
	if err != nil {
		return nil, err
	}
	err = rootMap.applyPatterns(vfs, files, beforePatterns)
	if err != nil {
		return nil, err
	}

	// Apply ignorefile patterns
	ierrFiles, err := rootMap.applyIgnorefile(vfs, ignorefile, options.Nested)
	if err != nil {
		return nil, err
	}
	for _, efile := range ierrFiles {
		errorFiles = append(errorFiles, efile)
	}

	// Apply after patterns
	afterPatterns, err := makePatterns(options.AfterPatterns)
	if err != nil {
		return nil, err
	}
	err = rootMap.applyPatterns(vfs, files, afterPatterns)
	if err != nil {
		return nil, err
	}

	return makeResult(vfs, basedir, rootMap, errorFiles)
}

func makeResult(vfs afero.Fs, basedir string,
	fileMap stateMap, errorFiles []string) (*MatchesResult, error) {
	matchedFiles := []string{}
	unmatchedFiles := []string{}
	matchedDirs := []string{}
	unmatchedDirs := []string{}
	errorDirs := []string{}
	for f, matched := range fileMap {
		if f == "" {
			continue
		}
		isDir, err := afero.IsDir(vfs, f)
		if err != nil {
			errorDirs = append(errorDirs, f)
			return nil, err
		}
		if isDir {
			if matched {
				matchedDirs = append(matchedDirs, f)
			} else {
				unmatchedDirs = append(unmatchedDirs, f)
			}
		} else {
			if matched {
				matchedFiles = append(matchedFiles, f)
			} else {
				unmatchedFiles = append(unmatchedFiles, f)
			}
		}
	}

	sort.Strings(matchedFiles)
	sort.Strings(unmatchedFiles)
	sort.Strings(errorFiles)
	sort.Strings(matchedDirs)
	sort.Strings(unmatchedDirs)
	sort.Strings(errorDirs)
	return &MatchesResult{
		BaseDir:        basedir,
		MatchedFiles:   matchedFiles,
		UnmatchedFiles: unmatchedFiles,
		ErrorFiles:     errorFiles,
		MatchedDirs:    matchedDirs,
		UnmatchedDirs:  unmatchedDirs,
		ErrorDirs:      errorDirs,
	}, nil
}
