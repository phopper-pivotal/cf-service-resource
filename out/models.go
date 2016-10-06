package out

import "github.com/idahobean/cf-sb-resource"

type Request struct {
	Source resource.Source `json:"source"`
	Params Params          `json:"params"`
}

type Params struct {
	Repository           string            `json:"repository"`
	CurrentAppName       string            `json:"current_app_name"`
	Memory               string            `json:"memory"`
	Disk                 string            `json:"disk"`
	HealthCheck          string            `json:"health_check"`
}

type Response struct {
	Version  resource.Version        `json:"version"`
	Metadata []resource.MetadataPair `json:"metadata"`
}
