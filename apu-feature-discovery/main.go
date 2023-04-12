// Copyright (c) 2019, NVIDIA CORPORATION. All rights reserved.

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/LYNXI/apu-feature-discovery/apihelper"
)

const (
	// Bin : Name of the binary
	Bin = "apu-feature-discovery"
)

var (
	// Version : Version of the binary
	// This will be set using ldflags at compile time
	Version = "0.1.0"
	// MachineTypePath : Path to the file describing the machine type
	// This will be override during unit testing
	MachineTypePath               = "/sys/class/dmi/id/product_name"
	OutputFilePath                = "/etc/kubernetes/node-feature-discovery/features.d/apu"
	SleepInterval   time.Duration = 1 * time.Minute
)

func main() {
	log.SetPrefix(Bin + ": ")

	if Version == "" {
		log.Print("Version is not set.")
		log.Fatal("Be sure to compile with '-ldflags \"-X main.Version=${AFD_VERSION}\"' and to set $AFD_VERSION")
	}

	log.Printf("Running %s in version %s", Bin, Version)

	lynxipcilib := NewLynxiPCILib()

	log.Print("Start running")
	err := run(lynxipcilib)
	if err != nil {
		log.Printf("Unexpected error: %v", err)
	}
	log.Print("Exiting")
}

func run(lynxipci LynxiPCI) error {
	defer func() {
		err := removeOutputFile(OutputFilePath)
		if err != nil {
			log.Printf("Warning: Error removing output file: %v", err)
		}
	}()

	labelHelp := apihelper.NewHelper(os.Getenv("NODE_NAME"))
	if err := labelHelp.Init(); err != nil {
		log.Printf("apihelper init Error %v", err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	exitChan := make(chan bool)

	go func() {
		select {
		case s := <-sigChan:
			log.Printf("Received signal \"%v\", shutting down.", s)
			exitChan <- true
		}
	}()

	afdLabels := make(map[string]string)
	afdLabels["lynxi.com/afd.timestamp"] = fmt.Sprintf("%d", time.Now().Unix())

L:
	for {
		log.Print("getAPULabels begin")
		APULabels, err := getAPULabels(lynxipci)
		if err != nil {
			return fmt.Errorf("Error generating APU labels: %v", err)
		}

		if len(APULabels) == 0 {
			log.Printf("Warning: no labels generated from any source")
		}

		labelHelp.UpdateNodeLabels(APULabels)
		log.Print("Sleeping for ", SleepInterval)

		select {
		case <-exitChan:
			break L
		case <-time.After(SleepInterval):
			break
		}
	}

	return nil
}

func getAPULabels(lynxipci LynxiPCI) (map[string]string, error) {
	devices, err := lynxipci.Devices()
	if err != nil {
		return nil, fmt.Errorf("Unable to get APU devices: %v", err)
	}
	labels := make(map[string]string)
	if len(devices) > 0 {
		labels["lynxi.com/apu.present"] = "true"
	} else {
		labels["lynxi.com/apu.present"] = "false"
	}
	return labels, nil
}

// writeLabelsToFile writes a set of labels to the specified path. The file is written atomocally
func writeLabelsToFile(path string, labelSets ...map[string]string) error {
	output := new(bytes.Buffer)
	for _, labels := range labelSets {
		for k, v := range labels {
			fmt.Fprintf(output, "%s=%s\n", k, v)
		}
	}
	err := writeFileAtomically(path, output.Bytes(), 0644)
	if err != nil {
		return fmt.Errorf("Error atomically writing file '%s': %v", path, err)
	}
	return nil
}

func writeFileAtomically(path string, contents []byte, perm os.FileMode) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("Failed to retrieve absolute path of output file: %v", err)
	}

	absDir := filepath.Dir(absPath)
	tmpDir := filepath.Join(absDir, "gfd-tmp")

	err = os.Mkdir(tmpDir, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return fmt.Errorf("Failed to create temporary directory: %v", err)
	}
	defer func() {
		if err != nil {
			os.RemoveAll(tmpDir)
		}
	}()

	tmpFile, err := ioutil.TempFile(tmpDir, "gfd-")
	if err != nil {
		return fmt.Errorf("Fail to create temporary output file: %v", err)
	}
	defer func() {
		if err != nil {
			tmpFile.Close()
			os.Remove(tmpFile.Name())
		}
	}()

	err = ioutil.WriteFile(tmpFile.Name(), contents, perm)
	if err != nil {
		return fmt.Errorf("Error writing temporary file '%v': %v", tmpFile.Name(), err)
	}

	err = os.Rename(tmpFile.Name(), path)
	if err != nil {
		return fmt.Errorf("Error moving temporary file to '%v': %v", path, err)
	}

	err = os.Chmod(path, perm)
	if err != nil {
		return fmt.Errorf("Error setting permissions on '%v': %v", path, err)
	}

	return nil
}

func removeOutputFile(path string) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("Failed to retrieve absolute path of output file: %v", err)
	}

	absDir := filepath.Dir(absPath)
	tmpDir := filepath.Join(absDir, "gfd-tmp")

	err = os.RemoveAll(tmpDir)
	if err != nil {
		return fmt.Errorf("Failed to remove temporary output directory: %v", err)
	}

	err = os.Remove(absPath)
	if err != nil {
		return fmt.Errorf("Failed to remove output file: %v", err)
	}

	return nil
}
