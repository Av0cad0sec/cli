// Package sgxs contains helper functions for dealing with ELF and SGXS binaries.
package sgxs

import (
	"os/exec"
	"strconv"
)

// Elf2Sgxs converts an ELF binary built for the SGX ABI into an SGXS binary.
//
// It requires the `ftxsgx-elf2sgxs` utility to be installed.
func Elf2Sgxs(elfSgxPath, sgxsPath string, heapSize, stackSize, threads uint64) error {
	args := []string{
		elfSgxPath,
		"--heap-size", strconv.FormatUint(heapSize, 10),
		"--stack-size", strconv.FormatUint(stackSize, 10),
		"--threads", strconv.FormatUint(threads, 10),
		"--output", sgxsPath,
	}

	cmd := exec.Command("ftxsgx-elf2sgxs", args...)
	return cmd.Run()
}
