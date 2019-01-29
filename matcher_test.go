package xignore

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMatches_Simple(t *testing.T) {
	matcher := NewSystemMatcher()
	result, err := matcher.Matches("testdata/simple", &MatchesOptions{
		Ignorefile: ".xignore",
	})
	require.NoError(t, err)

	require.Equal(t, result.MatchedFiles, []string{".xignore", "empty.log"})
	require.Equal(t, result.UnmatchedFiles, []string{"rain.txt"})
	require.Empty(t, result.ErrorFiles)
	require.Empty(t, result.MatchedDirs)
	require.Empty(t, result.UnmatchedDirs)
	require.Empty(t, result.ErrorDirs)
}

func TestMatches_Simple_WithBeforePatterns(t *testing.T) {
	matcher := NewSystemMatcher()
	result, err := matcher.Matches("testdata/simple", &MatchesOptions{
		Ignorefile:     ".xignore",
		BeforePatterns: []string{"rain.txt"},
	})
	require.NoError(t, err)

	require.Equal(t, result.MatchedFiles, []string{".xignore", "empty.log", "rain.txt"})
	require.Equal(t, result.UnmatchedFiles, []string{})
	require.Empty(t, result.ErrorFiles)
	require.Empty(t, result.MatchedDirs)
	require.Empty(t, result.UnmatchedDirs)
	require.Empty(t, result.ErrorDirs)
}

func TestMatches_Simple_WithAfterPatterns(t *testing.T) {
	matcher := NewSystemMatcher()
	result, err := matcher.Matches("testdata/simple", &MatchesOptions{
		Ignorefile:    ".xignore",
		AfterPatterns: []string{"!.xignore", "!empty.log", "rain.txt"},
	})
	require.NoError(t, err)

	require.Equal(t, result.MatchedFiles, []string{"rain.txt"})
	require.Equal(t, result.UnmatchedFiles, []string{".xignore", "empty.log"})
	require.Empty(t, result.ErrorFiles)
	require.Empty(t, result.MatchedDirs)
	require.Empty(t, result.UnmatchedDirs)
	require.Empty(t, result.ErrorDirs)
}

func TestMatches_Folder(t *testing.T) {
	matcher := NewSystemMatcher()
	result, err := matcher.Matches("testdata/folder", &MatchesOptions{
		Ignorefile: ".xignore",
	})
	require.NoError(t, err)

	require.Equal(t, []string{"foo/bar/1.txt"}, result.MatchedFiles)
	require.Equal(t, []string{".xignore", "foo/bar/tool/lex.txt", "foo/tar/2.txt"}, result.UnmatchedFiles)
	require.Empty(t, result.ErrorFiles)
	require.Equal(t, []string{"foo/bar"}, result.MatchedDirs)
	require.Equal(t, []string{"foo", "foo/bar/tool", "foo/tar"}, result.UnmatchedDirs)
	require.Empty(t, result.ErrorDirs)
}

func TestMatches_Root(t *testing.T) {
	matcher := NewSystemMatcher()
	result, err := matcher.Matches("testdata/root", &MatchesOptions{
		Ignorefile: ".xignore",
	})
	require.NoError(t, err)

	require.Equal(t, []string{"1.txt"}, result.MatchedFiles)
	require.Equal(t, []string{".xignore", "sub/1.txt", "sub/2.txt"}, result.UnmatchedFiles)
	require.Empty(t, result.ErrorFiles)
	require.Empty(t, result.MatchedDirs)
	require.Equal(t, []string{"sub"}, result.UnmatchedDirs)
	require.Empty(t, result.ErrorDirs)
}

func TestMatches_Exclusion(t *testing.T) {
	matcher := NewSystemMatcher()
	result, err := matcher.Matches("testdata/exclusion", &MatchesOptions{
		Ignorefile: ".xignore",
	})
	require.NoError(t, err)

	require.Equal(t, []string{"e1.txt", "e3.txt", "en/e3.txt"}, result.MatchedFiles)
	require.Equal(t, []string{"!", ".xignore", "e2.txt", "en/e1.txt", "en/e2.txt"}, result.UnmatchedFiles)
	require.Empty(t, result.ErrorFiles)
	require.Empty(t, result.MatchedDirs)
	require.Equal(t, []string{"en"}, result.UnmatchedDirs)
	require.Empty(t, result.ErrorDirs)
}

func TestMatches_DisabledNested(t *testing.T) {
	matcher := NewSystemMatcher()
	result, err := matcher.Matches("testdata/nested", &MatchesOptions{
		Ignorefile: ".xignore",
		Nested:     false,
	})
	require.NoError(t, err)

	require.Equal(t, []string{
		"inner/foo.md",
	}, result.MatchedFiles)
	require.Equal(t, []string{
		".xignore", "1.txt",
		"inner/.xignore", "inner/2.lst",
		"inner/inner2/.xignore", "inner/inner2/jess.ini", "inner/inner2/moss.ini",
	}, result.UnmatchedFiles)
	require.Empty(t, result.ErrorFiles)
	require.Empty(t, result.MatchedDirs)
	require.Equal(t, []string{"inner", "inner/inner2"}, result.UnmatchedDirs)
	require.Empty(t, result.ErrorDirs)
}

func TestMatches_Nested(t *testing.T) {
	matcher := NewSystemMatcher()
	result, err := matcher.Matches("testdata/nested", &MatchesOptions{
		Ignorefile: ".xignore",
		Nested:     true,
	})
	require.NoError(t, err)

	require.Equal(t, []string{
		"inner/2.lst", "inner/foo.md", "inner/inner2/moss.ini",
	}, result.MatchedFiles)
	require.Equal(t, []string{
		".xignore", "1.txt",
		"inner/.xignore",
		"inner/inner2/.xignore", "inner/inner2/jess.ini",
	}, result.UnmatchedFiles)
	require.Empty(t, result.ErrorFiles)
	require.Empty(t, result.MatchedDirs)
	require.Equal(t, result.UnmatchedDirs, []string{"inner", "inner/inner2"})
	require.Empty(t, result.ErrorDirs)
}
