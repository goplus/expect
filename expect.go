package expect

import (
	"os"
	"regexp"
	"time"

	"golang.org/x/crypto/ssh"

	expect "github.com/google/goexpect"
	term "github.com/google/goterm/term"
)

const (
	GopPackage = true
)

type gExpecter interface {
	MainEntry()
	Close() error
}

type Project struct {
	e    *expect.GExpect
	done <-chan error
}

func (p *Project) Close() error {
	if p.e != nil {
		return p.e.Close()
	}
	return nil
}

// Spawn starts a new process and collects the output. Arguments may not contain spaces.
func (p *Project) Spawn__0(
	command string, timeout time.Duration, handle func(err error), opts ...expect.Option) (err error) {
	p.Close()
	p.e, p.done, err = expect.Spawn(command, timeout, opts...)
	check(err, handle)
	return
}

// Spawn starts a new process and collects the output. Arguments may contain spaces.
func (p *Project) Spawn__1(
	command []string, timeout time.Duration, handle func(err error), opts ...expect.Option) (err error) {
	p.Close()
	p.e, p.done, err = expect.SpawnWithArgs(command, timeout, opts...)
	check(err, handle)
	return
}

// Spawn starts an interactive SSH session,ties it to a PTY and collects the output.
func (p *Project) Spawn__2(
	sshClient *ssh.Client, timeout time.Duration, handle func(err error), opts ...expect.Option) (err error) {
	p.Close()
	p.e, p.done, err = expect.SpawnSSH(sshClient, timeout, opts...)
	check(err, handle)
	return
}

// Spawn starts an interactive SSH session and ties it to a local PTY, optionally requests a remote PTY.
func (p *Project) Spawn__3(
	sshClient *ssh.Client, timeout time.Duration, term term.Termios, handle func(err error), opts ...expect.Option) (err error) {
	p.Close()
	p.e, p.done, err = expect.SpawnSSHPTY(sshClient, timeout, term, opts...)
	check(err, handle)
	return
}

// Spawn is used to write generic Spawners.
func (p *Project) Spawn__4(
	opt *expect.GenOptions, timeout time.Duration, handle func(err error), opts ...expect.Option) (err error) {
	p.Close()
	p.e, p.done, err = expect.SpawnGeneric(opt, timeout, opts...)
	check(err, handle)
	return
}

// Spawn spawns an expect.Batcher.
func (p *Project) Spawn__5(
	b []expect.Batcher, timeout time.Duration, handle func(err error), opts ...expect.Option) (err error) {
	p.Close()
	p.e, p.done, err = expect.SpawnFake(b, timeout, opts...)
	check(err, handle)
	return
}

// Expect reads spawned processes output looking for input regular expression.
// Timeout set to 0 makes Expect return the current buffer.
func (p *Project) Expect__0(pattern string, timeout ...time.Duration) (out string, match []string, err error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return
	}
	return p.e.Expect(re, duration(timeout))
}

// Expect checks each Case against the accumulated out buffer, sending specified
// string back. Leaving Send empty will Send nothing to the process.
// Substring expansion can be used eg.
// 	Case{`vf[0-9]{2}.[a-z]{3}[0-9]{2}\.net).*UP`,`show arp \1`}
// 	Given: vf11.hnd01.net            UP      35 (4)        34 (4)          CONNECTED         0              0/0
// 	Would send: show arp vf11.hnd01.net
func (p *Project) Expect__1(cs []expect.Caser, timeout ...time.Duration) (string, []string, int, error) {
	return p.e.ExpectSwitchCase(cs, duration(timeout))
}

// Expect takes an array of BatchEntries and runs through them in order. For every Expect
// command a BatchRes entry is created with output buffer and sub matches.
// Failure of any of the batch commands will stop the execution, returning the results up to the
// failure.
func (p *Project) Expect__2(bat []expect.Batcher, timeout ...time.Duration) ([]expect.BatchRes, error) {
	return p.e.ExpectBatch(bat, duration(timeout))
}

// Expect reads spawned processes output looking for input regular expression.
// Timeout set to 0 makes Expect return the current buffer.
func (p *Project) Expect__3(re *regexp.Regexp, timeout ...time.Duration) (string, []string, error) {
	return p.e.Expect(re, duration(timeout))
}

// Send sends a string to spawned process.
func (p *Project) Send__0(in string) error {
	return p.e.Send(in)
}

// Send sends a signal to the Expect controlled process.
// Only works on Process Expecters.
func (p *Project) Send__1(sig os.Signal) error {
	return p.e.SendSignal(sig)
}

// Gopt_Project_Main is required by Go+ compiler as the entry of a .expect project.
func Gopt_Project_Main(proj gExpecter) {
	proj.MainEntry()
	proj.Close()
}

func duration(timeout []time.Duration) (dur time.Duration) {
	if timeout != nil {
		return timeout[0]
	}
	return -1
}

func check(err error, handle func(err error)) {
	if err != nil {
		if handle == nil {
			panic(err)
		}
		handle(err)
	}
}
