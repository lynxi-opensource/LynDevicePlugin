// // Copyright (c) 2019, NVIDIA CORPORATION. All rights reserved.

// package main

// import (
// 	"fmt"
// 	"io/ioutil"
// 	"os"
// 	"regexp"
// 	"strings"
// 	"testing"
// 	"time"

// 	"github.com/stretchr/testify/require"
// )

// func TestRunOneshot(t *testing.T) {
// 	nvmlMock := NewTestNvmlMock()
// 	vgpuMock := NewTestVGPUMock()
// 	conf := Conf{true, true, "none", "./gfd-test-oneshot", time.Second, false}

// 	MachineTypePath = "/tmp/machine-type"
// 	machineType := []byte("product-name\n")
// 	err := ioutil.WriteFile("/tmp/machine-type", machineType, 0644)
// 	require.NoError(t, err, "Write machine type mock file")

// 	defer func() {
// 		err = os.Remove(MachineTypePath)
// 		require.NoError(t, err, "Removing machine type mock file")
// 	}()

// 	err = run(nvmlMock, vgpuMock, conf)
// 	require.NoError(t, err, "Error from run function")

// 	outFile, err := os.Open(conf.OutputFilePath)
// 	require.NoError(t, err, "Opening output file")

// 	defer func() {
// 		err = outFile.Close()
// 		require.NoError(t, err, "Closing output file")
// 		err = os.Remove(conf.OutputFilePath)
// 		require.NoError(t, err, "Removing output file")
// 	}()

// 	result, err := ioutil.ReadAll(outFile)
// 	require.NoError(t, err, "Reading output file")

// 	err = checkResult(result, "tests/expected-output.txt")
// 	require.NoError(t, err, "Checking result")

// 	err = checkResult(result, "tests/expected-output-vgpu.txt")
// 	require.NoError(t, err, "Checking result for vgpu labels")
// }

// func TestRunWithNoTimestamp(t *testing.T) {
// 	nvmlMock := NewTestNvmlMock()
// 	vgpuMock := NewTestVGPUMock()
// 	conf := Conf{true, true, "none", "./gfd-test-with-no-timestamp", time.Second, true}

// 	MachineTypePath = "/tmp/machine-type"
// 	machineType := []byte("product-name\n")
// 	err := ioutil.WriteFile("/tmp/machine-type", machineType, 0644)
// 	require.NoError(t, err, "Write machine type mock file")

// 	defer func() {
// 		err = os.Remove(MachineTypePath)
// 		require.NoError(t, err, "Removing machine type mock file")
// 	}()

// 	err = run(nvmlMock, vgpuMock, conf)
// 	require.NoError(t, err, "Error from run function")

// 	outFile, err := os.Open(conf.OutputFilePath)
// 	require.NoError(t, err, "Opening output file")

// 	defer func() {
// 		err = outFile.Close()
// 		require.NoError(t, err, "Closing output file")
// 		err = os.Remove(conf.OutputFilePath)
// 		require.NoError(t, err, "Removing output file")
// 	}()

// 	result, err := ioutil.ReadAll(outFile)
// 	require.NoError(t, err, "Reading output file")

// 	err = checkResult(result, "tests/expected-output.txt")
// 	require.NoError(t, err, "Checking result")
// 	require.NotContains(t, string(result), "nvidia.com/gfd.timestamp=", "Checking absent timestamp")

// 	err = checkResult(result, "tests/expected-output-vgpu.txt")
// 	require.NoError(t, err, "Checking result for vgpu labels")
// }

// func TestRunSleep(t *testing.T) {
// 	nvmlMock := NewTestNvmlMock()
// 	vgpuMock := NewTestVGPUMock()
// 	conf := Conf{false, true, "none", "./gfd-test-loop", time.Second, false}

// 	MachineTypePath = "/tmp/machine-type"
// 	machineType := []byte("product-name\n")
// 	err := ioutil.WriteFile("/tmp/machine-type", machineType, 0644)
// 	require.NoError(t, err, "Write machine type mock file")

// 	defer func() {
// 		err = os.Remove(MachineTypePath)
// 		require.NoError(t, err, "Removing machine type mock file")
// 		err = os.Remove(conf.OutputFilePath)
// 		require.NoError(t, err, "Removing output file")
// 	}()

// 	var runError error
// 	go func() {
// 		runError = run(nvmlMock, vgpuMock, conf)
// 	}()

// 	outFileModificationTime := make([]int64, 2)
// 	timestampLabels := make([]string, 2)
// 	// Read two iterations of the output file
// 	for i := 0; i < 2; i++ {
// 		outFile, err := waitForFile(conf.OutputFilePath, 5, time.Second)
// 		require.NoErrorf(t, err, "Open output file: %d", i)

// 		var outFileStat os.FileInfo
// 		var ts int64

// 		for attempt := 0; i > 0 && attempt < 3; attempt++ {
// 			// We ensure that the output file has been modified. Note, we expect the contents to remain the
// 			// same so we check the modification timestamp of the file.
// 			outFileStat, err = os.Stat(conf.OutputFilePath)
// 			require.NoError(t, err, "Getting output file info")

// 			ts = outFileStat.ModTime().Unix()
// 			if ts > outFileModificationTime[0] {
// 				break
// 			}
// 			// We wait for conf.SleepInterval, as the labels should be updated at least once in that period
// 			time.Sleep(conf.SleepInterval)
// 		}
// 		outFileModificationTime[i] = ts

// 		output, err := ioutil.ReadAll(outFile)
// 		require.NoErrorf(t, err, "Read output file: %d", i)

// 		err = outFile.Close()
// 		require.NoErrorf(t, err, "Close output file: %d", i)

// 		err = checkResult(output, "tests/expected-output.txt")
// 		require.NoErrorf(t, err, "Checking result: %d", i)
// 		err = checkResult(output, "tests/expected-output-vgpu.txt")
// 		require.NoErrorf(t, err, "Checking result for vgpu labels: %d", i)

// 		labels, err := buildLabelMapFromOutput(output)
// 		require.NoErrorf(t, err, "Building map of labels from output file: %d", i)

// 		require.Containsf(t, labels, "nvidia.com/gfd.timestamp", "Missing timestamp: %d", i)
// 		timestampLabels[i] = labels["nvidia.com/gfd.timestamp"]

// 		require.Containsf(t, labels, "nvidia.com/vgpu.present", "Missing vgpu present label: %d", i)
// 		require.Containsf(t, labels, "nvidia.com/vgpu.host-driver-version", "Missing vGPU host driver version label: %d", i)
// 		require.Containsf(t, labels, "nvidia.com/vgpu.host-driver-branch", "Missing vGPU host driver branch label: %d", i)
// 	}
// 	require.Greater(t, outFileModificationTime[1], outFileModificationTime[0], "Output file not modified")
// 	require.Equal(t, timestampLabels[1], timestampLabels[0], "Timestamp label changed")

// 	require.NoError(t, runError, "Error from run")
// }

// func TestFailOnNVMLInitError(t *testing.T) {
// 	nvmlMock := NewTestNvmlMock()
// 	vgpuMock := NewTestVGPUMock()
// 	conf := Conf{true, true, "none", "./gfd-test-loop", 500 * time.Millisecond, false}

// 	MachineTypePath = "/tmp/machine-type"
// 	machineType := []byte("product-name\n")
// 	err := ioutil.WriteFile("/tmp/machine-type", machineType, 0644)
// 	require.NoError(t, err, "Write machine type mock file")

// 	defer func() {
// 		err = os.Remove(MachineTypePath)
// 		require.NoError(t, err, "Removing machine type mock file")
// 	}()

// 	defer func() {
// 		// Remove the output file created by any "success" cases below
// 		err = os.Remove(conf.OutputFilePath)
// 		require.NoError(t, err, "Removing output file")
// 	}()

// 	// Test for case (errorOnInit = true, failOnInitError = true, no other errors)
// 	nvmlMock.errorOnInit = true
// 	conf.FailOnInitError = true
// 	conf.MigStrategy = "none"
// 	err = run(nvmlMock, vgpuMock, conf)
// 	require.Error(t, err, "Expected error from NVML Init")

// 	// Test for case (errorOnInit = true, failOnInitError = true, some other error)
// 	nvmlMock.errorOnInit = true
// 	conf.FailOnInitError = true
// 	conf.MigStrategy = "bogus"
// 	err = run(nvmlMock, vgpuMock, conf)
// 	require.Error(t, err, "Expected error from NVML Init")

// 	// Test for case (errorOnInit = true, failOnInitError = false, no other errors)
// 	nvmlMock.errorOnInit = true
// 	conf.FailOnInitError = false
// 	conf.MigStrategy = "none"
// 	err = run(nvmlMock, vgpuMock, conf)
// 	require.NoError(t, err, "Expected to skip error from NVML Init")

// 	// Test for case (errorOnInit = true, failOnInitError = false, some other error)
// 	nvmlMock.errorOnInit = true
// 	conf.FailOnInitError = false
// 	conf.MigStrategy = "bogus"
// 	err = run(nvmlMock, vgpuMock, conf)
// 	require.NoError(t, err, "Expected to skip error from NVML Init")

// 	// Test for case (errorOnInit = false, failOnInitError = true, no other errors)
// 	nvmlMock.errorOnInit = false
// 	conf.FailOnInitError = true
// 	conf.MigStrategy = "none"
// 	err = run(nvmlMock, vgpuMock, conf)
// 	require.NoError(t, err, "Expected no errors")

// 	// Test for case (errorOnInit = false, failOnInitError = true, some other error)
// 	nvmlMock.errorOnInit = false
// 	conf.FailOnInitError = true
// 	conf.MigStrategy = "bogus"
// 	err = run(nvmlMock, vgpuMock, conf)
// 	require.Error(t, err, "Expected error since MIGStrategy is 'bogus'")

// 	// Test for case (errorOnInit = false, failOnInitError = false, no other errors)
// 	nvmlMock.errorOnInit = false
// 	conf.FailOnInitError = false
// 	conf.MigStrategy = "none"
// 	err = run(nvmlMock, vgpuMock, conf)
// 	require.NoError(t, err, "Expected no errors")

// 	// Test for case (errorOnInit = false, failOnInitError = false, some other error)
// 	nvmlMock.errorOnInit = false
// 	conf.FailOnInitError = false
// 	conf.MigStrategy = "bogus"
// 	err = run(nvmlMock, vgpuMock, conf)
// 	require.Error(t, err, "Expected error since MIGStrategy is 'bogus'")
// }

// func buildLabelMapFromOutput(output []byte) (map[string]string, error) {
// 	labels := make(map[string]string)

// 	lines := strings.Split(strings.TrimRight(string(output), "\n"), "\n")
// 	for _, line := range lines {
// 		split := strings.Split(line, "=")
// 		if len(split) != 2 {
// 			return nil, fmt.Errorf("Unexpected format in line: '%v'", line)
// 		}
// 		key := split[0]
// 		value := split[1]

// 		if v, ok := labels[key]; ok {
// 			return nil, fmt.Errorf("Duplicate label '%v': %v (overwrites %v)", key, v, value)
// 		}
// 		labels[key] = value
// 	}

// 	return labels, nil
// }

// func checkResult(result []byte, expectedOutputPath string) error {
// 	expected, err := ioutil.ReadFile(expectedOutputPath)
// 	if err != nil {
// 		return fmt.Errorf("Opening expected output file: %v", err)
// 	}

// 	var expectedRegexps []*regexp.Regexp
// 	for _, line := range strings.Split(strings.TrimRight(string(expected), "\n"), "\n") {
// 		expectedRegexps = append(expectedRegexps, regexp.MustCompile(line))
// 	}

// LOOP:
// 	for _, line := range strings.Split(strings.TrimRight(string(result), "\n"), "\n") {
// 		if expectedOutputPath == "tests/expected-output-vgpu.txt" {
// 			if !strings.Contains(line, "vgpu") {
// 				// ignore other labels when vgpu file is specified
// 				continue
// 			}
// 		} else {
// 			if strings.Contains(line, "vgpu") {
// 				// ignore vgpu labels when non vgpu file is specified
// 				continue
// 			}
// 		}
// 		for _, regex := range expectedRegexps {
// 			if regex.MatchString(line) {
// 				continue LOOP
// 			}
// 		}
// 		return fmt.Errorf("Line does not match any regexp: %v", string(line))
// 	}
// 	return nil
// }

// func waitForFile(fileName string, iter int, sleepInterval time.Duration) (*os.File, error) {
// 	for i := 0; i < iter-1; i++ {
// 		file, err := os.Open(fileName)
// 		if err != nil && os.IsNotExist(err) {
// 			time.Sleep(sleepInterval)
// 			continue
// 		}
// 		if err != nil {
// 			return nil, err
// 		}
// 		return file, nil
// 	}
// 	return os.Open(fileName)
// }
