// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package knownhosts_test

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

func TestHostKeyPreferenceBasedOnKnownHosts(t *testing.T) {
	sshConfig := &ssh.ClientConfig{
		// HostKeyAlgorithms: []string{"ssh-rsa"},
	}

	assert := assertions{t}

	file, err := ioutil.TempFile("", "knownhosts")
	assert.NilError(err)
	defer func() {
		assert.NilError(os.Remove(file.Name()))
	}()
	githubKnownHost := "github.com,140.82.118.3 ssh-rsa " +
		"AAAAB3NzaC1yc2EAAAABIwAAAQEAq2A7hRGmdnm9tUDbO9IDSwBK6TbQa+PXYPCPy6rbTrTt" +
		"w7PHkccKrpp0yVhp5HdEIcKr6pLlVDBfOLX9QUsyCOV0wzfjIJNlGEYsdlLJizHhbn2mUjvS" +
		"AHQqZETYP81eFzLQNnPHt4EVVUh7VfDESU84KezmD5QlWpXLmvU31/yMf+Se8xhHTvKSCZIF" +
		"ImWwoG6mbUoWf9nzpIoaSjB+weqqUUmpaaasXVal72J+UX2B+2RPW3RcT0eOzQgqlJL3RKrT" +
		"JvdsjE3JEAvGq3lGHSZXy28G3skua2SmVi/w4yCE6gbODqnTWlg7+wC604ydGXA8VJiS5ap4" +
		"3JXiUFFAaQ=="
	assert.NilError(ioutil.WriteFile(file.Name(), []byte(githubKnownHost), 0o400))
	assert.NilError(file.Close())

	hostKeyCallback, err := knownhosts.New(file.Name())
	assert.NilError(err)
	sshConfig.HostKeyCallback = hostKeyCallback

	cli, err := ssh.Dial("tcp", "github.com:22", sshConfig)
	defer func() {
		if cli != nil {
			err = cli.Close()
			assert.NilError(err)
		}
	}()
	assert.Equal(err.Error(), "ssh: handshake failed: knownhosts: key mismatch")
}

type assertions struct {
	testing.TB
}

func (a assertions) NilError(err error) {
	a.Helper()
	if err != nil {
		a.Fatal(err)
	}
}

func (a assertions) Equal(got interface{}, want interface{}) {
	a.Helper()
	if !reflect.DeepEqual(got, want) {
		a.Fatalf("expected aren't equal:\n got: %#v \nwant: %#v", got, want)
	}
}
