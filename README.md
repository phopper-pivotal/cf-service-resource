# Cloud Foundry Docker Image Resource

An output only resource that will push Docker Image to a
Cloud Foundry deployment.  
based on [concourse/cf-resource](https://github.com/concourse/cf-resource)

## Source Configuration

* `api`: *Required.* The address of the Cloud Controller in the Cloud Foundry
  deployment.
* `username`: *Required.* The username used to authenticate.
* `password`: *Required.* The password used to authenticate.
* `organization`: *Required.* The organization to push the application to.
* `space`: *Required.* The space to push the application to.
* `skip_cert_check`: *Optional.* Check the validity of the CF SSL cert.
  Defaults to `false`.

## Behaviour

### `out`: Deploy Docker Image to a Cloud Foundry

Pushes Docker Image to the Cloud Foundry. 

#### Parameters

* `repository`: *Required.* A path to Docker Image.
* `current_app_name`: *Required.* The name of the application.
* `memory`: *Optional.* Memory limit. (e.g. 256M, 1024M, 1G)
* `disk`: *Optional.* Disk limit. (e.g. 256M, 1024M, 1G)
* `health_check`: *Optional.* `port/none` Default `port`.
* `delete_app`: *Optional.* Default `false`. (not yet implemented)

## Pipeline example

```yaml
---
resource_types:
  - name: cf-sb-resource
    type: docker-image
    source:
      repository: idahobean/cf-sb-resource

resources:
  - name: foobar-cf
    type: cf-sb-resource
    source:
      api: https://api.foo.bar.cfapps.io
      username: USERNAME
      password: PASSWORD
      organization: ORG
      space: SPACE
      skip_cert_check: false

jobs:
- name: job-deploy-docker-image
  public: true
  serial: true
  plan:
  - put: foobar-cf
    params:
      repository: cloudfoundry/lattice-app
      current_app_name: lattice-foobar
      memory: 256M
      disk: 256M
      health_check: none

```
