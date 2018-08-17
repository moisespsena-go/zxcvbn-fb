package zxcvbn_fb

import (
	"bytes"
	"encoding/json"
	"os/exec"
	"syscall"
)

type ZxcvbnResult struct {
	Score    int
	Feedback struct {
		Warning     string
		Suggestions []string
	}
}

func Zxcvbn(password string, userInputs ...string) (*ZxcvbnResult, error) {
	var args []string
	for _, ui := range userInputs {
		args = append(args, "--user-input", ui)
	}

	cmd := exec.Command("zxcvbn", args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Pdeathsig: syscall.SIGTERM,
	}
	var in, out bytes.Buffer
	in.WriteString(password + "\n")
	cmd.Stdout, cmd.Stdin = &out, &in
	err := cmd.Start()
	if err != nil {
		return nil, err
	}
	err = cmd.Wait()
	if err != nil {
		return nil, err
	}
	dec := json.NewDecoder(&out)
	var result ZxcvbnResult
	err = dec.Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
