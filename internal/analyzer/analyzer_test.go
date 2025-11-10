package analyzer

import (
	"testing"
)

func TestCategorizeCommits(t *testing.T) {
	tests := []struct {
		name     string
		commits  []string
		wantCats CommitCategories
	}{
		{
			name: "all conventional commit types",
			commits: []string{
				"feat: add new feature",
				"fix: resolve bug",
				"docs: update README",
				"chore: update dependencies",
				"refactor: improve code structure",
				"test: add unit tests",
			},
			wantCats: CommitCategories{
				Features: []string{"add new feature"},
				Fixes:    []string{"resolve bug"},
				Docs:     []string{"update README"},
				Chores:   []string{"update dependencies", "add unit tests"},
				Changes:  []string{"improve code structure"},
				Breaking: []string{},
			},
		},
		{
			name: "breaking changes",
			commits: []string{
				"breaking: remove old API",
				"feat: add BREAKING CHANGE in body",
				"fix!: breaking fix",
			},
			wantCats: CommitCategories{
				Features: []string{},
				Fixes:    []string{},
				Docs:     []string{},
				Chores:   []string{},
				Changes:  []string{},
				Breaking: []string{
					"breaking: remove old API",
					"feat: add BREAKING CHANGE in body",
					"fix!: breaking fix",
				},
			},
		},
		{
			name: "uncategorized commits",
			commits: []string{
				"update version",
				"merge pull request",
				"initial commit",
			},
			wantCats: CommitCategories{
				Features: []string{},
				Fixes:    []string{},
				Docs:     []string{},
				Chores:   []string{},
				Changes: []string{
					"update version",
					"merge pull request",
					"initial commit",
				},
				Breaking: []string{},
			},
		},
		{
			name: "feature prefix variations",
			commits: []string{
				"feat: new feature",
				"feature: another feature",
				"FEAT: uppercase feature",
			},
			wantCats: CommitCategories{
				Features: []string{
					"new feature",
					"another feature",
					"uppercase feature",
				},
				Fixes:    []string{},
				Docs:     []string{},
				Chores:   []string{},
				Changes:  []string{},
				Breaking: []string{},
			},
		},
		{
			name:    "empty commits",
			commits: []string{},
			wantCats: CommitCategories{
				Features: []string{},
				Fixes:    []string{},
				Docs:     []string{},
				Chores:   []string{},
				Changes:  []string{},
				Breaking: []string{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CategorizeCommits(tt.commits)

			// Compare each category
			if !equalStringSlices(got.Features, tt.wantCats.Features) {
				t.Errorf("Features = %v, want %v", got.Features, tt.wantCats.Features)
			}
			if !equalStringSlices(got.Fixes, tt.wantCats.Fixes) {
				t.Errorf("Fixes = %v, want %v", got.Fixes, tt.wantCats.Fixes)
			}
			if !equalStringSlices(got.Docs, tt.wantCats.Docs) {
				t.Errorf("Docs = %v, want %v", got.Docs, tt.wantCats.Docs)
			}
			if !equalStringSlices(got.Chores, tt.wantCats.Chores) {
				t.Errorf("Chores = %v, want %v", got.Chores, tt.wantCats.Chores)
			}
			if !equalStringSlices(got.Changes, tt.wantCats.Changes) {
				t.Errorf("Changes = %v, want %v", got.Changes, tt.wantCats.Changes)
			}
			if !equalStringSlices(got.Breaking, tt.wantCats.Breaking) {
				t.Errorf("Breaking = %v, want %v", got.Breaking, tt.wantCats.Breaking)
			}
		})
	}
}

func TestExtractMessage(t *testing.T) {
	tests := []struct {
		name     string
		commit   string
		prefixes []string
		want     string
	}{
		{
			name:     "feat prefix",
			commit:   "feat: add new feature",
			prefixes: []string{"feat:"},
			want:     "add new feature",
		},
		{
			name:     "multiple prefixes",
			commit:   "feature: add new feature",
			prefixes: []string{"feat:", "feature:"},
			want:     "add new feature",
		},
		{
			name:     "preserve case in message",
			commit:   "FEAT: Add New Feature",
			prefixes: []string{"feat:"},
			want:     "Add New Feature",
		},
		{
			name:     "no matching prefix",
			commit:   "no prefix here",
			prefixes: []string{"feat:"},
			want:     "no prefix here",
		},
		{
			name:     "whitespace handling",
			commit:   "fix:   fix with spaces  ",
			prefixes: []string{"fix:"},
			want:     "fix with spaces",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractMessage(tt.commit, tt.prefixes...)
			if got != tt.want {
				t.Errorf("extractMessage() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestCommitCategories_CountTotal(t *testing.T) {
	tests := []struct {
		name string
		cats CommitCategories
		want int
	}{
		{
			name: "all categories populated",
			cats: CommitCategories{
				Features: []string{"a", "b"},
				Fixes:    []string{"c"},
				Docs:     []string{"d"},
				Chores:   []string{"e", "f"},
				Changes:  []string{"g"},
				Breaking: []string{"h"},
			},
			want: 8,
		},
		{
			name: "empty categories",
			cats: CommitCategories{},
			want: 0,
		},
		{
			name: "only features",
			cats: CommitCategories{
				Features: []string{"a", "b", "c"},
			},
			want: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.cats.CountTotal()
			if got != tt.want {
				t.Errorf("CountTotal() = %d, want %d", got, tt.want)
			}
		})
	}
}

// Helper function to compare string slices
func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
