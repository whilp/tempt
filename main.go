// tempt runs a command in a temporary directory and then cleans up.
/*
 * Copyright (c) 2016 Will Maier <wcmaier@m.aier.us>
 *
 * Permission to use, copy, modify, and distribute this software for any
 * purpose with or without fee is hereby granted, provided that the above
 * copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
 * WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
 * MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
 * ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
 * WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
 * ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
 * OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 */
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"syscall"
)

var (
	version   string
	fVersion  = flag.Bool("version", false, "print version and exit")
	fPreserve = flag.Bool("preserve", false, "preserve temporary directory")
)

func main() {
	flag.Usage = usage
	flag.Parse()
	args := flag.Args()

	if *fVersion {
		fmt.Printf("tempt %s\n", version)
		os.Exit(0)
	}

	exit := 0
	err, exit := inner(args[0], args[1:]...)
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(exit)
}

func inner(command string, args ...string) (error, int) {
	exit := 0
	dir, err := ioutil.TempDir("", "tempt-")
	if err != nil {
		return err, exit
	}
	
	// Bail if we can't convert to an absolute path; relative paths 
	// are too risky to blindly remove.
	abs, err := filepath.Abs(dir)
	if err != nil {
		return err, exit
	}

	dir = abs
	defer cleanup(dir)

	err = os.Chdir(dir)
	if err != nil {
		return err, exit
	}

	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Env = os.Environ()

	err = cmd.Run()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			waitStatus := exitError.Sys().(syscall.WaitStatus)
			exit = waitStatus.ExitStatus()
		} else {
			log.Panic(err)
		}
	}

	return nil, exit
}

func cleanup(dir string) {
	if *fPreserve {
		log.Printf("preserving temporary directory %s\n", dir)
		return
	}
	err := os.RemoveAll(dir)
	if err != nil {
		log.Fatal(err)
	}
}

func usage() {
	self := path.Base(os.Args[0])
	fmt.Fprintf(os.Stderr, "usage: %s ARG [ARG ...]\n\n", self)
	fmt.Fprint(os.Stderr, "Run ARG in a temporary directory and then cleanup.\n\n")
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Arguments:")
	fmt.Fprintln(os.Stderr, "  ARG    Command to run, followed by optional arguments")
	fmt.Fprintln(os.Stderr, "Environment variables:")
	fmt.Fprintln(os.Stderr, "  TMPDIR Location in which to create the temporary directory")
	os.Exit(2)
}
