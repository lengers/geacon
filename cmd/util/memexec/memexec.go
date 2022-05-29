package memexec

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

// Exec is an in-memory executable code unit.
type Exec struct {
	f *os.File
}

// New creates new memory execution object that can be
// used for executing commands on a memory based binary.
func New(b []byte) (*Exec, error) {
	f, err := open(b)
	if err != nil {
		return nil, err
	}
	return &Exec{f: f}, nil
}

// Command is an equivalent of `exec.Command`,
// except that the path to the executable is be omitted.
func (m *Exec) Command(arg ...string) *exec.Cmd {
	return exec.Command(m.f.Name(), arg...)
}

func (m *Exec) CommandAsUser(username []byte, password []byte, arg ...string) *exec.Cmd {

	path := strings.ReplaceAll(m.f.Name(), "self", fmt.Sprint(os.Getpid()))
	chmodArgs := []string{"-c", fmt.Sprintf("chmod 555 %s", path)}
	chmodCmd := exec.Command("/bin/sh", chmodArgs...)
	chmodCmd.Run()

	cradle := fmt.Sprintf("\"echo '%s' | su - %s -c %s\"", password, username, path)
	fmt.Printf("CRADLE: %s\n", cradle)
	time.Sleep(60 * time.Second)
	args := []string{"-c", cradle}
	return exec.Command("/bin/sh", args...)
}

func (m *Exec) File() *os.File {
	return m.f
}

// Close closes Exec object.
//
// Any further command will fail, it's client's responsibility
// to control the flow by using synchronization algorithms.
func (m *Exec) Close() error {
	return clean(m.f)
}
