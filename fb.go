package zxcvbn_fb

import (
	"bytes"
	"encoding/json"
	"os"
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
	var args = append([]string{
		"-c",
		`import json
import sys

from zxcvbn import zxcvbn

class JSONEncoder(json.JSONEncoder):
    def default(self, o):
        try:
            return super(JSONEncoder, self).default(o)
        except TypeError:
            return str(o)

res = zxcvbn(sys.argv[1], user_inputs=sys.argv[2:])
json.dump(res, sys.stdout, indent=2, cls=JSONEncoder)
sys.stdout.write('\n')
`,
		password,
	}, userInputs...)

	cmd := exec.Command("python3", args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
	cmd.Stderr = os.Stderr

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
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
