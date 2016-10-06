package out

import (
	"os"
	"os/exec"
)

type PAAS interface {
	Login(api string, username string, password string, insecure bool) error
	Target(organization string, space string) error
	PushApp(repository string, currentAppName string, memory string, disk string, healthCheck string) error
}

type CloudFoundry struct{}

func NewCloudFoundry() *CloudFoundry {
	return &CloudFoundry{}
}

func (cf *CloudFoundry) Login(api string, username string, password string, insecure bool) error {
	args := []string{"api", api}
	if insecure {
		args = append(args, "--skip-ssl-validation")
	}

	err := cf.cf(args...).Run()
	if err != nil {
		return err
	}

	return cf.cf("auth", username, password).Run()
}

func (cf *CloudFoundry) Target(organization string, space string) error {
	return cf.cf("target", "-o", organization, "-s", space).Run()
}

func (cf *CloudFoundry) PushApp(repository string, currentAppName string, memory string, disk string, healthCheck string) error {
	args := []string{}
	options := []string{}

	if memory != "" {
		options = append(options, "-m", memory)
	}
	if disk != "" {
		options = append(options, "-k", disk)
	}
	options = append(options, "-u", healthCheck)

	args = append(args, "push", currentAppName, "-o", repository)
	args = append(args, options...)
	return cf.cf(args...).Run()
}

func (cf *CloudFoundry) cf(args ...string) *exec.Cmd {
	cmd := exec.Command("cf", args...)
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(), "CF_COLOR=true")

	return cmd
}
