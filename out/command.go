package out

import (
	"time"

	"github.com/idahobean/cf-sb-resource"
)

type Command struct {
	paas PAAS
}

func NewCommand(paas PAAS) *Command {
	return &Command{
		paas: paas,
	}
}

func (command *Command) Run(request Request) (Response, error) {
	err := command.paas.Login(
		request.Source.API,
		request.Source.Username,
		request.Source.Password,
		request.Source.SkipCertCheck,
	)
	if err != nil {
		return Response{}, err
	}

	err = command.paas.Target(
		request.Source.Organization,
		request.Source.Space,
	)
	if err != nil {
		return Response{}, err
	}

	err = command.paas.PushApp(
		request.Params.Repository,
		request.Params.CurrentAppName,
		request.Params.Memory,
		request.Params.Disk,
		request.Params.HealthCheck,
	)
	if err != nil {
		return Response{}, err
	}

	return Response{
		Version: resource.Version{
			Timestamp: time.Now(),
		},
		Metadata: []resource.MetadataPair{
			{
				Name:  "organization",
				Value: request.Source.Organization,
			},
			{
				Name:  "space",
				Value: request.Source.Space,
			},
		},
	}, nil
}
