/*Command protoc-gen-gitaly is designed to be used as a protobuf compiler
plugin to verify Gitaly processes are being followed when writing RPC's.

Usage

The protoc-gen-gitaly linter can be chained into any protoc workflow that
requires verification that Gitaly RPC guidelines are followed. Typically
this can be done by adding the following argument to an existing protoc
command:

  --gitaly_out=.

For example, you may add the linter as an argument to the command responsible
for generating Go code:

  protoc --go_out=. --gitaly_out=. *.proto

Or, you can run the Gitaly linter by itself. To try out, run the following
command while in the project root:

  protoc --gitaly_out=. ./go/internal/linter/testdata/incomplete.proto

You should see some errors printed to screen for improperly written
RPC's in the incomplete.proto file.

Prerequisites

The protobuf compiler (protoc) can be obtained from the GitHub page:
https://github.com/protocolbuffers/protobuf/releases

Background

The protobuf compiler accepts plugins to analyze protobuf files and generate
language specific code.

These plugins require the following executable naming convention:

  protoc-gen-$NAME

Where $NAME is the plugin name of the compiler desired. The protobuf compiler
will search the PATH until an executable with that name is found for a
desired plugin. For example, the following protoc command:

  protoc --gitaly_out=. *.proto

The above will search the PATH for an executable named protoc-gen-gitaly

The plugin accepts a protobuf message in STDIN that describes the parsed
protobuf files. A response is sent back on STDOUT that contains any errors.
*/
package main

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"gitlab.com/gitlab-org/gitaly/v14/proto/go/internal/linter"
	"google.golang.org/protobuf/proto"
)

const (
	gitalyProtoDirArg = "proto_dir"
	gitalypbDirArg    = "gitalypb_dir"
)

func main() {
	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalf("reading input: %s", err)
	}

	req := &plugin.CodeGeneratorRequest{}

	if err := proto.Unmarshal(data, req); err != nil {
		log.Fatalf("parsing input proto: %s", err)
	}

	if err := lintProtos(req); err != nil {
		log.Fatal(err)
	}

	if err := generateProtolistGo(req); err != nil {
		log.Fatal(err)
	}
}

func parseArgs(argString string) (gitalyProtoDir string, gitalypbDir string) {
	for _, arg := range strings.Split(argString, ",") {
		argKeyValue := strings.Split(arg, "=")
		if len(argKeyValue) != 2 {
			continue
		}
		switch argKeyValue[0] {
		case gitalyProtoDirArg:
			gitalyProtoDir = argKeyValue[1]
		case gitalypbDirArg:
			gitalypbDir = argKeyValue[1]
		}
	}

	return gitalyProtoDir, gitalypbDir
}

func lintProtos(req *plugin.CodeGeneratorRequest) error {
	var errMsgs []string
	for _, pf := range req.GetProtoFile() {
		errs := linter.LintFile(pf, req)
		for _, err := range errs {
			errMsgs = append(errMsgs, err.Error())
		}
	}

	resp := &plugin.CodeGeneratorResponse{}

	if len(errMsgs) > 0 {
		errMsg := strings.Join(errMsgs, "\n\t")
		resp.Error = &errMsg
	}

	// Send back the results.
	data, err := proto.Marshal(resp)
	if err != nil {
		return fmt.Errorf("failed to marshal output proto: %s", err)
	}

	_, err = os.Stdout.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write output proto: %s", err)
	}
	return nil
}

func generateProtolistGo(req *plugin.CodeGeneratorRequest) error {
	var err error
	gitalyProtoDir, gitalypbDir := parseArgs(req.GetParameter())

	if gitalyProtoDir == "" {
		return fmt.Errorf("%s not provided", gitalyProtoDirArg)
	}
	if gitalypbDir == "" {
		return fmt.Errorf("%s not provided", gitalypbDirArg)
	}

	var protoNames []string

	if gitalyProtoDir, err = filepath.Abs(gitalyProtoDir); err != nil {
		return fmt.Errorf("failed to get absolute path for %s: %v", gitalyProtoDir, err)
	}

	files, err := ioutil.ReadDir(gitalyProtoDir)
	if err != nil {
		return fmt.Errorf("failed to read %s: %v", gitalyProtoDir, err)
	}

	for _, fi := range files {
		if !fi.IsDir() && strings.HasSuffix(fi.Name(), ".proto") {
			protoNames = append(protoNames, fmt.Sprintf(`"%s"`, fi.Name()))
		}
	}

	f, err := os.Create(filepath.Join(gitalypbDir, "protolist.go"))
	if err != nil {
		return fmt.Errorf("could not create protolist.go: %v", err)
	}
	defer f.Close()

	if err = renderProtoList(f, protoNames); err != nil {
		return fmt.Errorf("could not render go code: %v", err)
	}

	return nil
}

// renderProtoList generate a go file with a list of gitaly protos
func renderProtoList(dest io.WriteCloser, protoNames []string) error {
	var joinFunc = template.FuncMap{"join": strings.Join}
	protoList := `package gitalypb
	// Code generated by protoc-gen-gitaly. DO NOT EDIT

	// GitalyProtos is a list of gitaly protobuf files
	var GitalyProtos = []string{
		{{join . ",\n"}},
	}
	`
	protoListTempl, err := template.New("protoList").Funcs(joinFunc).Parse(protoList)
	if err != nil {
		return fmt.Errorf("could not create go code template: %v", err)
	}

	var rawGo bytes.Buffer

	if err := protoListTempl.Execute(&rawGo, protoNames); err != nil {
		return fmt.Errorf("could not execute go code template: %v", err)
	}

	formattedGo, err := format.Source(rawGo.Bytes())
	if err != nil {
		return fmt.Errorf("could not format go code: %v", err)
	}

	if _, err = io.Copy(dest, bytes.NewBuffer(formattedGo)); err != nil {
		return fmt.Errorf("failed to write protolist.go file: %v", err)
	}

	return nil
}
