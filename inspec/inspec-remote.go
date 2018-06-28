package inspec

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"

	"github.com/armon/circbuf"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	linereader "github.com/mitchellh/go-linereader"
)

const (
	// maxBufSize limits how much output we collect from a local
	// invocation. This is to prevent TF memory usage from growing
	// to an enormous amount due to a faulty process.
	maxBufSize = 8 * 1024
)

func runRemote(ctx context.Context, s *terraform.InstanceState, data *schema.ResourceData, o terraform.UIOutput, profiles string) error {
	cmd := "inspec exec " + profiles + " -t aws://"
	return inspecExec(ctx, o, cmd, "", nil)
}

func inspecExec(ctx context.Context, o terraform.UIOutput, command string, workingdir string, env []string) error {

	var cmdargs []string

	if runtime.GOOS == "windows" {
		cmdargs = []string{"cmd", "/C"}
	} else {
		cmdargs = []string{"/bin/sh", "-c"}
	}
	cmdargs = append(cmdargs, command)

	// Setup the reader that will read the output from the command.
	// We use an os.Pipe so that the *os.File can be passed directly to the
	// process, and not rely on goroutines copying the data which may block.
	// See golang.org/issue/18874
	pr, pw, err := os.Pipe()
	if err != nil {
		return fmt.Errorf("failed to initialize pipe for output: %s", err)
	}

	var cmdEnv []string
	cmdEnv = os.Environ()
	cmdEnv = append(cmdEnv, env...)

	// Setup the command
	cmd := exec.Command(cmdargs[0], cmdargs[1:]...)
	cmd.Stderr = pw
	cmd.Stdout = pw
	// Dir specifies the working directory of the command.
	// If Dir is the empty string (this is default), runs the command
	// in the calling process's current directory.
	cmd.Dir = workingdir
	// Env specifies the environment of the command.
	// By default will use the calling process's environment
	cmd.Env = cmdEnv

	output, _ := circbuf.NewBuffer(maxBufSize)

	// Write everything we read from the pipe to the output buffer too
	tee := io.TeeReader(pr, output)

	// copy the teed output to the UI output
	copyDoneCh := make(chan struct{})
	go copyOutputChan(o, tee, copyDoneCh)

	// Output what we're about to run
	o.Output(fmt.Sprintf("Executing: %q", cmdargs))

	// Start the command
	err = cmd.Start()
	if err == nil {
		err = cmd.Wait()
	}

	// Close the write-end of the pipe so that the goroutine mirroring output
	// ends properly.
	pw.Close()

	// Cancelling the command may block the pipe reader if the file descriptor
	// was passed to a child process which hasn't closed it. In this case the
	// copyOutput goroutine will just hang out until exit.
	select {
	case <-copyDoneCh:
	case <-ctx.Done():
	}

	if err != nil {
		return fmt.Errorf("Error running command '%s': %v. Output: %s",
			command, err, output.Bytes())
	}

	return nil
}

func copyOutputChan(o terraform.UIOutput, r io.Reader, doneCh chan<- struct{}) {
	defer close(doneCh)
	lr := linereader.New(r)
	for line := range lr.Ch {
		o.Output(line)
	}
}
