package cmd

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
	"testing"

	"github.com/icodezjb/atomicswap/logger"
)

const (
	node1DataDir = "testing/integration/node1/"
	node2DataDir = "testing/integration/node2/"
	node1Config  = "testing/integration/node1/config.json"
	node2Config  = "testing/integration/node2/config.json"
)

var (
	noclean       = flag.Bool("noclean", false, "after test not clean the test data ")
	node1InitArgs = []string{
		"init", node1DataDir + "poa.json", "--datadir", node1DataDir,
	}

	node2InitArgs = []string{
		"init", node2DataDir + "poa.json", "--datadir", node2DataDir,
	}

	node1RunArgs = []string{
		"--identity", "node1", "--rpc", "--rpcport", "7545", "--datadir", node1DataDir, "--port",
		"30303", "--nodiscover", "--allow-insecure-unlock", "--mine", "--password", node1DataDir + "password.txt",
		"--unlock", "0xae6e5fee5161cede9bc4d89effbbf9944867127d", "--keystore", node1DataDir,
	}

	node2RunArgs = []string{
		"--identity", "node1", "--rpc", "--rpcport", "8545", "--datadir", node2DataDir, "--port",
		"30304", "--nodiscover", "--allow-insecure-unlock", "--mine", "--password", node2DataDir + "password.txt",
		"--unlock", "0x75a8f951632c2e550906f31b53b7923f45be5157", "--keystore", node2DataDir,
	}
)

func TestMain(m *testing.M) {
	var (
		node1 = newTestGeth(node1DataDir, node1InitArgs, node1RunArgs)
		node2 = newTestGeth(node2DataDir, node2InitArgs, node2RunArgs)
	)

	flag.Parse()

	logger.Info("geth path: %s", node1.path)

	node1.Cleanup()
	node2.Cleanup()

	os.Exit(testMainWrapper(m, node1, node2))
}

func testMainWrapper(m *testing.M, nodes ...*TestCmd) int {
	log.SetOutput(ioutil.Discard) // is noisy otherwise
	defer log.SetOutput(os.Stderr)

	for _, node := range nodes {
		node.MustRun()
	}

	defer func() {
		rc := recover()
		for _, node := range nodes {
			node.KillExit()
		}

		if rc != nil {
			panic(rc)
		}
	}()

	return m.Run()
}

type TestCmd struct {
	path     string
	initArgs []string
	runArgs  []string
	cmd      *exec.Cmd
	mu       sync.Mutex
	logFile  *os.File

	DataDir string
	Cleanup func()
}

func newTestGeth(dir string, initArgs []string, runArgs []string) *TestCmd {
	tc := new(TestCmd)

	tc.DataDir = dir
	tc.initArgs = initArgs
	tc.runArgs = runArgs

	switch runtime.GOOS {
	case "darwin":
		tc.path = "testing/geth-darwin-amd64-1.9.8-d62e9b28"
	case "linux":
		tc.path = "testing/geth-linux-amd64-1.9.8-d62e9b28"
	default:
		logger.FatalError("os '%s' is not support", runtime.GOOS)
	}

	gethLog, err := os.OpenFile(dir+"geth.log", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		logger.FatalError("Fatal to open geth.log: %v", err)
	}

	tc.logFile = gethLog

	tc.Cleanup = func() {
		exec.Command("/usr/bin/pkill", "geth").CombinedOutput() //nolint:errcheck

		files, _ := filepath.Glob(dir + "geth*")
		for _, file := range files {
			os.RemoveAll(file) //nolint:errcheck
		}
	}

	return tc
}

func (tc *TestCmd) Run() error {
	if out, err := exec.Command(tc.path, tc.initArgs...).CombinedOutput(); err != nil {
		return fmt.Errorf("Failed to init node: %v\n%s", err, out)
	}

	tc.cmd = &exec.Cmd{
		Path: tc.path,
		Args: append([]string{tc.path}, tc.runArgs...),
	}

	if tc.logFile != nil {
		tc.cmd.Stdout = tc.logFile
		tc.cmd.Stderr = tc.logFile
	}

	return tc.cmd.Start()
}

func (tc *TestCmd) MustRun() {
	if err := tc.Run(); err != nil {
		panic(fmt.Sprintf("Fatal to geth run: %v, path=%v", err, tc.runArgs))
	}
}

func (tc *TestCmd) KillExit() {
	if err := tc.cmd.Process.Signal(os.Interrupt); err != nil {
		logger.Warn("%v", err)
	}

	if tc.logFile != nil {
		tc.logFile.Close()
	}

	if !*noclean {
		tc.Cleanup()
	}
}

func (tc *TestCmd) Write(b []byte) (n int, err error) {
	lines := bytes.Split(b, []byte("\n"))
	for _, line := range lines {
		if len(line) > 0 {
			logger.Info("(stderr) %s", line)
		}
	}

	tc.mu.Lock()
	tc.logFile.Write(b) //nolint:errcheck
	tc.mu.Unlock()

	return len(b), err
}
