package project

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var limitStdout sync.Mutex

func captureOutput(f func()) []byte {
	limitStdout.Lock()
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	f()
	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout
	limitStdout.Unlock()
	return out
}

func TestDependsOn(t *testing.T) {
	t.Parallel()
	p := New("../examples/tag", "..")
	if p.DependsOn[0] != fmt.Sprintf("%s%c%s", p.RootDir, os.PathSeparator, "project") {
		t.Errorf(".gally.yml should contain a depends_on tag with project, instead %s", p.DependsOn[0])
		t.FailNow()
	}

	p = New("../examples/notag", "..")
	if p.DependsOn[0] != fmt.Sprintf("%s%c%s", p.RootDir, os.PathSeparator, "repo") {
		t.Errorf(".gally.yml should contain a depends_on tag with repo, instead %s", p.DependsOn[0])
		t.FailNow()
	}
}
func TestEnv(t *testing.T) {
	t.Parallel()
	p := New("../examples/tag", "..")
	if p.Env[0].Name != "NAMESPACE" || p.Env[0].Value != "staging" {
		t.Errorf(".gally.yml should contain NAMESPACE=staging")
		t.FailNow()
	}
}

func TestFindProjectPaths(t *testing.T) {
	t.Parallel()
	paths, err := findPaths("..")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	expected := []string{"../examples/notag", "../examples/tag"}
	if !cmp.Equal(paths, expected) {
		t.Errorf("paths %v must be equal to %v", paths, expected)
		t.FailNow()
	}
}

func TestNew(t *testing.T) {
	t.Parallel()
	p := New("../examples/tag", "..")
	if !strings.HasPrefix(p.BaseDir, p.Dir) {
		t.Errorf("%q directory is not in %q", p.BaseDir, p.Dir)
		t.FailNow()
	}
}

func TestRun(t *testing.T) {
	t.Parallel()
	expected := []byte("world\n")
	c := Project{Scripts: map[string]string{"hello": "echo world"}}
	out := captureOutput(func() {
		if err := c.Run("hello"); err != nil {
			t.Error(err)
			t.FailNow()
		}
	})
	if !cmp.Equal(out, expected) {
		t.Errorf("output must be equal to %q but is equal to %q", expected, out)
		t.FailNow()
	}

	expected = []byte("env.go\nproject.go\nproject_test.go\nstrategies.go\n")
	c = Project{Scripts: map[string]string{"list": "ls"}}
	out = captureOutput(func() {
		if err := c.Run("list"); err != nil {
			t.Error(err)
			t.FailNow()
		}
	})
	if !cmp.Equal(out, expected) {
		t.Errorf("output must be equal to %q but is equal to %q", expected, out)
		t.FailNow()
	}
}

func TestRunBuild(t *testing.T) {
	t.Parallel()
	c := Project{BuildScript: "echo go building $GALLY_VERSION!"}
	out := captureOutput(func() {
		if err := c.runBuild("test"); err != nil {
			t.Error(err)
			t.FailNow()
		}
	})
	expected := []byte("go building test!\n")
	if !cmp.Equal(out, expected) {
		t.Errorf("output must be equal to %q but is equal to %q", expected, out)
		t.FailNow()
	}
}

func TestVersion(t *testing.T) {
	t.Parallel()
	c := Project{BaseDir: ".", Dir: "../examples/tag", VersionScript: "head -1 VERSION"}
	expected := "0.3.5"
	if c.Version() != expected {
		t.Errorf("version must be equal to %q but is actually %q", expected, c.Version())
		t.FailNow()
	}
}
