package inspec

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/hashicorp/terraform/communicator"
	"github.com/hashicorp/terraform/communicator/remote"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	linereader "github.com/mitchellh/go-linereader"
	"github.com/rs/zerolog/log"
)

func runLocal(ctx context.Context, s *terraform.InstanceState, data *schema.ResourceData, o terraform.UIOutput, profiles string) error {
	// Get a new communicator
	comm, err := communicator.New(s)
	if err != nil {
		return err
	}

	// Collect the scripts
	scripts, err := inspecScript(o, data, profiles)
	if err != nil {
		return err
	}
	for _, s := range scripts {
		defer s.Close()
	}

	// Copy and execute each script
	if err := runScripts(ctx, o, comm, scripts); err != nil {
		return err
	}

	return nil
}

func inspecScript(o terraform.UIOutput, d *schema.ResourceData, profiles string) ([]io.ReadCloser, error) {
	lines := []string{
		"curl -L https://omnitruck.chef.io/install.sh | sudo bash -s -- -P inspec",
		"inspec exec " + profiles,
	}

	script := strings.Join(lines, "\n")
	scripts := []io.ReadCloser{ioutil.NopCloser(bytes.NewReader([]byte(script)))}
	return scripts, nil
}

// runScripts is used to copy and execute a set of scripts
func runScripts(
	ctx context.Context,
	o terraform.UIOutput,
	comm communicator.Communicator,
	scripts []io.ReadCloser) error {

	retryCtx, cancel := context.WithTimeout(ctx, comm.Timeout())
	defer cancel()

	// Wait and retry until we establish the connection
	err := communicator.Retry(retryCtx, func() error {
		return comm.Connect(o)
	})
	if err != nil {
		return err
	}

	// Wait for the context to end and then disconnect
	go func() {
		<-ctx.Done()
		comm.Disconnect()
	}()

	for _, script := range scripts {
		var cmd *remote.Cmd

		outR, outW := io.Pipe()
		errR, errW := io.Pipe()
		defer outW.Close()
		defer errW.Close()

		go copyToOutput(o, outR)
		go copyToOutput(o, errR)

		remotePath := comm.ScriptPath()

		if err := comm.UploadScript(remotePath, script); err != nil {
			return fmt.Errorf("Failed to upload script: %v", err)
		}

		cmd = &remote.Cmd{
			Command: remotePath,
			Stdout:  outW,
			Stderr:  errW,
		}
		if err := comm.Start(cmd); err != nil {
			return fmt.Errorf("Error starting script: %v", err)
		}

		if err := cmd.Wait(); err != nil {
			return err
		}

		// Upload a blank follow up file in the same path to prevent residual
		// script contents from remaining on remote machine
		empty := bytes.NewReader([]byte(""))
		if err := comm.Upload(remotePath, empty); err != nil {
			// This feature is best-effort.
			log.Printf("[WARN] Failed to upload empty follow up script: %v", err)
		}
	}

	return nil
}

func copyToOutput(o terraform.UIOutput, r io.Reader) {
	lr := linereader.New(r)
	for line := range lr.Ch {
		o.Output(line)
	}
}
