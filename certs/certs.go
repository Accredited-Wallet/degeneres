package certs

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func Certs(cfg Config) {
	fmt.Println("Generating certs...")

	opensslConfig, err := filepath.Abs(cfg.OpensslConfig)
	if err != nil {
		fmt.Println("Failed getting absolute path:", err)
		return
	}

	os.Mkdir(cfg.OutputPath, os.ModePerm)

	if err := os.Chdir(cfg.OutputPath); err != nil {
		fmt.Println("Failed chdir to outputPath:", cfg.OutputPath)
		return
	}

	keys := []string{
		"ca.cer",
		"ca.key",
		"server.csr",
		"server.key",
		"server.cer",
		"serial",
	}
	for _, key := range keys {
		os.Remove(key)
	}

	// # Courtesy of https://github.com/deckarep/EasyCert
	cmds := []string{
		// "openssl genrsa -out ca.key 4096",
		// `openssl req -x509 -new -key ca.key -out ca.cer -days 90 -subj /CN="rms1000watt"`,
		// "openssl genrsa -out server.key 4096",
		// "openssl req -new -out server.csr -key server.key -config " + opensslConfig,
		// "openssl x509 -req -in server.csr -out server.cer -days 90 -CAkey ca.key -CA ca.cer -CAcreateserial -CAserial serial -extensions v3_ext -extfile ./openssl.cnf",
		"openssl req -x509 -out server.cer -newkey rsa:4096 -keyout server.key -days 365 -nodes -config " + opensslConfig + " -extensions v3_ext",
	}

	for _, cmd := range cmds {
		if err := execute(cmd); err != nil {
			return
		}
	}
}

func execute(cmd string) (err error) {
	cmdArr := strings.Split(cmd, " ")
	if len(cmdArr) < 2 {
		return errors.New("bad command provided")
	}

	name := cmdArr[0]
	args := cmdArr[1:]

	outBytes, err := exec.Command(name, args...).CombinedOutput()
	fmt.Println(string(outBytes))

	if err != nil {
		fmt.Printf("Failed executing cmd: '%s': %s\n", cmd, err)
		return err
	}

	return
}