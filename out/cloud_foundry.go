package out

import (
	"os"
	"os/exec"
)

type PAAS interface {
	Login(api string, username string, password string, insecure bool) error
	Target(organization string, space string) error
	Service(service string) error
	CreateService(service string, plan string, instanceName string, parametersAsJson string) error
	UpdateService(service string, plan string, parametersAsJson string) error
  DeleteService(service string)
	BindService(currentAppName string, instanceName string) error
	RestageApp(currentAppName string) error
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

func (cf *CloudFoundry) Service(service string) error {
	err := cf.cf("service", service).Run()

	if err != nil {
		return err
	}
}

func (cf *CloudFoundry) CreateService(service string, plan string, instanceName string, parametersAsJson string) error {
	args := []string{}
	args = append(args, "create-service", service, plan, instanceName

	if parametersAsJson != "" {
		args = append("-c", parametersAsJson)
	}

	return cf.cf(args...).Run()
}

func (cf *CloudFoundry) UpdateService(service string, plan string, parametersAsJson string) error {
  args := []string{};
  args = append(args, "update-service", service)

  if plan != "" {
     args = append("-p", plan)
	}

  if parametersAsJson != "" {
     args = append("-c", parametersAsJson)
	}

  return cf.cf(args...).Run()
}

func (cf *CloudFoundry) DeleteService(service string) error {
   args := []string{}
   args = append(args, "delete-service", service)

   return cf.cf(args...).Run()
}

func (cf *CloudFoundry) BindService(currentAppName string, instanceName string) error {
	args := []string{}
	args = append(args, "bind-service", currentAppName, instanceName)

	return cf.cf(args...).Run()
}

func (cf *CloudFoundry) RestageApp(currentAppName string) error {
	args := []string{}
	args = append(args, "restage", currentAppName)

	return cf.cf(args...).Run()
}

func (cf *CloudFoundry) cf(args ...string) *exec.Cmd {
	cmd := exec.Command("cf", args...)
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(), "CF_COLOR=true")

	return cmd
}
