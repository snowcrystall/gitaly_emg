package command

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m, goleak.IgnoreTopFunction("go.opencensus.io/stats/view.(*worker).start"))
}

func TestNewCommandExtraEnv(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	extraVar := "FOOBAR=123456"
	buff := &bytes.Buffer{}
	cmd, err := New(ctx, exec.Command("/usr/bin/env"), nil, buff, nil, extraVar)

	require.NoError(t, err)
	require.NoError(t, cmd.Wait())

	require.Contains(t, strings.Split(buff.String(), "\n"), extraVar)
}

func TestNewCommandExportedEnv(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	testCases := []struct {
		key   string
		value string
	}{
		{
			key:   "HOME",
			value: "/home/git",
		},
		{
			key:   "PATH",
			value: "/usr/bin:/bin:/usr/sbin:/sbin:/usr/local/bin",
		},
		{
			key:   "LD_LIBRARY_PATH",
			value: "/path/to/your/lib",
		},
		{
			key:   "TZ",
			value: "foobar",
		},
		{
			key:   "GIT_TRACE",
			value: "true",
		},
		{
			key:   "GIT_TRACE_PACK_ACCESS",
			value: "true",
		},
		{
			key:   "GIT_TRACE_PACKET",
			value: "true",
		},
		{
			key:   "GIT_TRACE_PERFORMANCE",
			value: "true",
		},
		{
			key:   "GIT_TRACE_SETUP",
			value: "true",
		},
		{
			key:   "all_proxy",
			value: "http://localhost:4000",
		},
		{
			key:   "http_proxy",
			value: "http://localhost:5000",
		},
		{
			key:   "HTTP_PROXY",
			value: "http://localhost:6000",
		},
		{
			key:   "https_proxy",
			value: "https://localhost:5000",
		},
		{
			key:   "HTTPS_PROXY",
			value: "https://localhost:6000",
		},
		{
			key:   "no_proxy",
			value: "https://excluded:5000",
		},
		{
			key:   "NO_PROXY",
			value: "https://excluded:5000",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.key, func(t *testing.T) {
			oldValue, exists := os.LookupEnv(tc.key)
			defer func() {
				if !exists {
					require.NoError(t, os.Unsetenv(tc.key))
					return
				}
				require.NoError(t, os.Setenv(tc.key, oldValue))
			}()
			require.NoError(t, os.Setenv(tc.key, tc.value))

			buff := &bytes.Buffer{}
			cmd, err := New(ctx, exec.Command("/usr/bin/env"), nil, buff, nil)
			require.NoError(t, err)
			require.NoError(t, cmd.Wait())

			expectedEnv := fmt.Sprintf("%s=%s", tc.key, tc.value)
			require.Contains(t, strings.Split(buff.String(), "\n"), expectedEnv)
		})
	}
}

func TestNewCommandUnexportedEnv(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	unexportedEnvKey, unexportedEnvVal := "GITALY_UNEXPORTED_ENV", "foobar"

	oldValue, exists := os.LookupEnv(unexportedEnvKey)
	defer func() {
		if !exists {
			require.NoError(t, os.Unsetenv(unexportedEnvKey))
			return
		}
		require.NoError(t, os.Setenv(unexportedEnvKey, oldValue))
	}()

	require.NoError(t, os.Setenv(unexportedEnvKey, unexportedEnvVal))

	buff := &bytes.Buffer{}
	cmd, err := New(ctx, exec.Command("/usr/bin/env"), nil, buff, nil)

	require.NoError(t, err)
	require.NoError(t, cmd.Wait())

	require.NotContains(t, strings.Split(buff.String(), "\n"), fmt.Sprintf("%s=%s", unexportedEnvKey, unexportedEnvVal))
}

func TestRejectEmptyContextDone(t *testing.T) {
	defer func() {
		p := recover()
		if p == nil {
			t.Error("expected panic, got none")
			return
		}

		if _, ok := p.(contextWithoutDonePanic); !ok {
			panic(p)
		}
	}()

	_, err := New(context.Background(), exec.Command("true"), nil, nil, nil)
	require.NoError(t, err)
}

func TestNewCommandTimeout(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	defer func(ch chan struct{}, t time.Duration) {
		spawnTokens = ch
		spawnConfig.Timeout = t
	}(spawnTokens, spawnConfig.Timeout)

	// This unbuffered channel will behave like a full/blocked buffered channel.
	spawnTokens = make(chan struct{})
	// Speed up the test by lowering the timeout
	spawnTimeout := 200 * time.Millisecond
	spawnConfig.Timeout = spawnTimeout

	testDeadline := time.After(1 * time.Second)
	tick := time.After(spawnTimeout / 2)

	errCh := make(chan error)
	go func() {
		_, err := New(ctx, exec.Command("true"), nil, nil, nil)
		errCh <- err
	}()

	var err error
	timePassed := false

wait:
	for {
		select {
		case err = <-errCh:
			break wait
		case <-tick:
			timePassed = true
		case <-testDeadline:
			t.Fatal("test timed out")
		}
	}

	require.True(t, timePassed, "time must have passed")
	require.Error(t, err)
	require.Contains(t, err.Error(), "process spawn timed out after")
}

func TestCommand_Wait_interrupts_after_context_timeout(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctx, timeout := context.WithTimeout(ctx, time.Second)
	defer timeout()

	cmd, err := New(ctx, exec.CommandContext(ctx, "sleep", "3"), nil, nil, nil)
	require.NoError(t, err)

	completed := make(chan error, 1)
	go func() { completed <- cmd.Wait() }()

	select {
	case err := <-completed:
		require.Error(t, err)
		s, ok := ExitStatus(err)
		require.True(t, ok)
		require.Equal(t, -1, s)
	case <-time.After(2 * time.Second):
		require.FailNow(t, "process is running too long")
	}
}

func TestNewCommandWithSetupStdin(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	value := "Test value"
	output := bytes.NewBuffer(nil)

	cmd, err := New(ctx, exec.Command("cat"), SetupStdin, nil, nil)
	require.NoError(t, err)

	_, err = fmt.Fprintf(cmd, "%s", value)
	require.NoError(t, err)

	// The output of the `cat` subprocess should exactly match its input
	_, err = io.CopyN(output, cmd, int64(len(value)))
	require.NoError(t, err)
	require.Equal(t, value, output.String())

	require.NoError(t, cmd.Wait())
}

func TestNewCommandNullInArg(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_, err := New(ctx, exec.Command("sh", "-c", "hello\x00world"), nil, nil, nil)
	require.Error(t, err)
	require.EqualError(t, err, `detected null byte in command argument "hello\x00world"`)
}

func TestNewNonExistent(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cmd, err := New(ctx, exec.Command("command-non-existent"), nil, nil, nil)
	require.Nil(t, cmd)
	require.Error(t, err)
}

func TestCommandStdErr(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var stdout, stderr bytes.Buffer
	expectedMessage := `hello world\nhello world\nhello world\nhello world\nhello world\n`

	logger := logrus.New()
	logger.SetOutput(&stderr)

	ctx = ctxlogrus.ToContext(ctx, logrus.NewEntry(logger))

	cmd, err := New(ctx, exec.Command("./testdata/stderr_script.sh"), nil, &stdout, nil)
	require.NoError(t, err)
	require.Error(t, cmd.Wait())

	assert.Empty(t, stdout.Bytes())
	require.Equal(t, expectedMessage, extractMessage(stderr.String()))
}

func TestCommandStdErrLargeOutput(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var stdout, stderr bytes.Buffer

	logger := logrus.New()
	logger.SetOutput(&stderr)

	ctx = ctxlogrus.ToContext(ctx, logrus.NewEntry(logger))

	cmd, err := New(ctx, exec.Command("./testdata/stderr_many_lines.sh"), nil, &stdout, nil)
	require.NoError(t, err)
	require.Error(t, cmd.Wait())

	assert.Empty(t, stdout.Bytes())
	msg := strings.ReplaceAll(extractMessage(stderr.String()), "\\n", "\n")
	require.LessOrEqual(t, len(msg), maxStderrBytes)
}

func TestCommandStdErrBinaryNullBytes(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var stdout, stderr bytes.Buffer

	logger := logrus.New()
	logger.SetOutput(&stderr)

	ctx = ctxlogrus.ToContext(ctx, logrus.NewEntry(logger))

	cmd, err := New(ctx, exec.Command("./testdata/stderr_binary_null.sh"), nil, &stdout, nil)
	require.NoError(t, err)
	require.Error(t, cmd.Wait())

	assert.Empty(t, stdout.Bytes())
	msg := strings.SplitN(extractMessage(stderr.String()), "\\n", 2)[0]
	require.Equal(t, strings.Repeat("\\x00", maxStderrLineLength), msg)
}

func TestCommandStdErrLongLine(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var stdout, stderr bytes.Buffer

	logger := logrus.New()
	logger.SetOutput(&stderr)

	ctx = ctxlogrus.ToContext(ctx, logrus.NewEntry(logger))

	cmd, err := New(ctx, exec.Command("./testdata/stderr_repeat_a.sh"), nil, &stdout, nil)
	require.NoError(t, err)
	require.Error(t, cmd.Wait())

	assert.Empty(t, stdout.Bytes())
	require.Contains(t, stderr.String(), fmt.Sprintf("%s\\n%s", strings.Repeat("a", maxStderrLineLength), strings.Repeat("b", maxStderrLineLength)))
}

func TestCommandStdErrMaxBytes(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var stdout, stderr bytes.Buffer

	logger := logrus.New()
	logger.SetOutput(&stderr)

	ctx = ctxlogrus.ToContext(ctx, logrus.NewEntry(logger))

	cmd, err := New(ctx, exec.Command("./testdata/stderr_max_bytes_edge_case.sh"), nil, &stdout, nil)
	require.NoError(t, err)
	require.Error(t, cmd.Wait())

	assert.Empty(t, stdout.Bytes())
	require.Equal(t, maxStderrBytes, len(strings.ReplaceAll(extractMessage(stderr.String()), "\\n", "\n")))
}

var logMsgRegex = regexp.MustCompile(`msg="(.+?)"`)

func extractMessage(logMessage string) string {
	subMatches := logMsgRegex.FindStringSubmatch(logMessage)
	if len(subMatches) != 2 {
		return ""
	}

	return subMatches[1]
}
