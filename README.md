# yaml-docs
=========
[![Go Report Card](https://goreportcard.com/badge/github.com/andbron/yaml-docs)](https://goreportcard.com/report/github.com/andbron/yaml-docs)

The yaml-docs tool auto-generates documentation from yaml files into markdown files as a table with each of the files
values, their defaults, and an optional description parsed from comments.

**Note:** This was originally forked from [norwoodj/helm-docs]( https://github.com/norwoodj/helm-docs ), which is build
around [Helm](helm.sh) charts, and there are a number of hold-overs in the examples and naming.


## Usage

```default
Usage:
  yaml-docs [flags]

Flags:
  -d, --dry-run                    don't actually render any markdown files just print to stdout passed
  -h, --help                       help for yaml-docs
  -l, --log-level string           Level of logs that should printed, one of (panic, fatal, error, warning, info, debug, trace) (default "info")
  -o, --output-file string         markdown file path relative to input template to which rendered documentation will be written (default "README.md")
  -s, --sort-values-order string   order in which to sort the values table ("alphanum" or "file") (default "alphanum")
  -t, --template-files strings     gotemplate file paths relative to each chart directory from which documentation will be generated (default [README.md.gotmpl])
  -f, --values-file strings        yaml values file to be parsed into values table. Can be specified multiple times
```

The markdown generation is entirely [gotemplate](https://golang.org/pkg/text/template) driven. The tool parses metadata
from charts and generates a number of sub-templates that can be referenced in a template file (by default `README.md.gotmpl`).
If no template file is provided, the tool has a default internal template that will generate a reasonably formatted README.

The most useful aspect of this tool is the auto-detection of field descriptions from comments:

```yaml
config:
  databasesToCreate:
    # -- default database for storage of database metadata
    - postgres

    # -- database for the [hashbash](https://github.com/norwoodj/hashbash-backend-go) project
    - hashbash

  usersToCreate:
    # -- admin user
    - {name: root, admin: true}

    # -- user with access to the database with the same name
    - {name: hashbash, readwriteDatabases: [hashbash]}

statefulset:
  image:
    # -- Image to use for deploying, must support an entrypoint which creates users/databases from appropriate config files
    repository: jnorwood/postgresql
    tag: "11"

  # -- Additional volumes to be mounted into the database container
  extraVolumes:
    - name: data
      emptyDir: {}
```

Resulting in a resulting README section like so:

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| config.databasesToCreate[0] | string | `"postgresql"` | default database for storage of database metadata |
| config.databasesToCreate[1] | string | `"hashbash"` | database for the [hashbash](https://github.com/norwoodj/hashbash-backend-go) project |
| config.usersToCreate[0] | object | `{"admin":true,"name":"root"}` | admin user |
| config.usersToCreate[1] | object | `{"name":"hashbash","readwriteDatabases":["hashbash"]}` | user with access to the database with the same name |
| statefulset.extraVolumes | list | `[{"emptyDir":{},"name":"data"}]` | Additional volumes to be mounted into the database container |
| statefulset.image.repository | string | `"jnorwood/postgresql:11"` | Image to use for deploying, must support an entrypoint which creates users/databases from appropriate config files |
| statefulset.image.tag | string | `"18.0831"` |  |

You'll notice that some complex fields (lists and objects) are documented while others aren't, and that some simple fields
like `statefulset.image.tag` are documented even without a description comment. The rules for what is and isn't documented in
the final table will be described in detail later in this document.

## Installation
<!-- helm-docs can be installed using [homebrew](https://brew.sh/): -->

<!-- ```bash -->
<!-- brew install norwoodj/tap/helm-docs -->
<!-- ``` -->

<!-- or [scoop](https://scoop.sh): -->

<!-- ```bash -->
<!-- scoop install helm-docs -->
<!-- ``` -->

<!-- This will download and install the [latest release](https://github.com/norwoodj/helm-docs/releases/latest) -->
<!-- of the tool. -->

To build from source in this repository:

```bash
cd cmd/yaml-docs
go build
```

Or install from source:

```bash
GO111MODULE=on go get github.com/andbron/yaml-docs/cmd/yaml-docs
```

## Usage

### Pre-commit hook

If you want to automatically generate `README.md` files with a pre-commit hook, make sure you
[install the pre-commit binary](https://pre-commit.com/#install), and add a [.pre-commit-config.yaml file](./.pre-commit-config.yaml)
to your project. Then run:

```bash
pre-commit install
pre-commit install-hooks
```

Future changes to your project's `yaml` or `README.md.gotmpl` files will cause an update to documentation when you commit.

### Running the binary directly

To run and generate documentation into READMEs:

```bash
yaml-docs -f my-values.yaml
# OR
yaml-docs --dry-run # prints generated documentation to stdout rather than modifying READMEs
```

<!-- ### Using docker -->

<!-- You can mount a directory with charts under `/helm-docs` within the container. -->

<!-- Then run: -->

<!-- ```bash -->
<!-- docker run --rm --volume "$(pwd):/helm-docs" -u $(id -u) jnorwood/helm-docs:latest -->
<!-- ``` -->

## Markdown Rendering

`--template-files` specifies the list of gotemplate files that should be used in rendering the resulting markdown file
for each chart found.
By default `--template-files=README.md.gotmpl`.

Files are always interpreted as being _relative to the working directory_.

If any of the specified template files is not found for a chart (you'll notice most of the example charts do not have a README.md.gotmpl)
file, then the internal default template is used instead.

The tool also includes the [sprig templating library](https://github.com/Masterminds/sprig), so those functions can be used
in the templates you supply.

### values.yaml metadata
This tool can parse descriptions and defaults of values from `values.yaml` files. The defaults are pulled directly from
the yaml in the file. 

It was formerly the case that descriptions had to be specified with the full path of the yaml field. This is no longer
the case, although it is still supported. Where before you would document a values.yaml like so:

```yaml
controller:
  publishService:
    # controller.publishService.enabled -- Whether to expose the ingress controller to the public world
    enabled: false

  # controller.replicas -- Number of nginx-ingress pods to load balance between.
  # Do not set this below 2.
  replicas: 2
```

You may now equivalently write:
```yaml
controller:
  publishService:
    # -- Whether to expose the ingress controller to the public world
    enabled: false

  # -- Number of nginx-ingress pods to load balance between.
  # Do not set this below 2.
  replicas: 2
```

New-style comments are much the same as the old-style comments, except that while old comments for a field could appear
anywhere in the file, new-style comments must appear **on the line(s) immediately preceding the field being documented.**

I invite you to check out the [example-charts](./example-charts) to see how this is done in practice. The `but-auto-comments`
examples in particular document the new comment format.

Note that comments can continue on the next line. In that case leave out the double dash, and the lines will simply be
appended with a space in-between, as in the `controller.replicas` field in the example above

The following rules are used to determine which values will be added to the values table in the README:

* By default, only _leaf nodes_, that is, fields of type `int`, `string`, `float`, `bool`, empty lists, and empty maps
  are added as rows in the values table. These fields will be added even if they do not have a description comment
* Lists and maps which contain elements will not be added as rows in the values table _unless_ they have a description
  comment which refers to them
* Adding a description comment for a non-empty list or map in this way makes it so that leaf nodes underneath the
  described field will _not_ be automatically added to the values table. In order to document both a non-empty list/map
  _and_ a leaf node within that field, description comments must be added for both

e.g. In this case, both `controller.livenessProbe` and `controller.livenessProbe.httpGet.path` will be added as rows in
the values table, but `controller.livenessProbe.httpGet.port` will not
```yaml
controller:
  # -- Configure the healthcheck for the ingress controller
  livenessProbe:
    httpGet:
      # -- This is the liveness check endpoint
      path: /healthz
      port: http
```

Results in:

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| controller.livenessProbe | object | `{"httpGet":{"path":"/healthz","port":8080}}` | Configure the healthcheck for the ingress controller |
| controller.livenessProbe.httpGet.path | string | `"/healthz"` | This is the liveness check endpoint |

If we remove the comment for `controller.livenessProbe` however, both leaf nodes `controller.livenessProbe.httpGet.path`
and `controller.livenessProbe.httpGet.port` will be added to the table, with or without description comments:

```yaml
controller:
  livenessProbe:
    httpGet:
      # -- This is the liveness check endpoint
      path: /healthz
      port: http
```

Results in:

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| controller.livenessProbe.httpGet.path | string | `"/healthz"` | This is the liveness check endpoint |
| controller.livenessProbe.httpGet.port | string | `"http"` | |


### nil values
If you would like to define a key for a value, but leave the default empty, you can still specify a description for it
as well as a type. This is possible with both the old and the new comment format:
```yaml
controller:
  # -- (int) Number of nginx-ingress pods to load balance between
  replicas:
  
  # controller.image -- (string) Number of nginx-ingress pods to load balance between
  image:
```
This could be useful when wanting to enforce user-defined values for the chart, where there are no sensible defaults.

### Default values/column
In cases where you do not want to include the default value from `values.yaml`, or where the real default is calculated
inside the chart, you can change the contents of the column like so:

```yaml
service:
  # -- Add annotations to the service, this is going to be a long comment across multiple lines
  # but that's fine, these will be concatenated and the @default will be rendered as the default for this field
  # @default -- the chart will add some internal annotations automatically
  annotations: []
```

The order is important. The first comment line(s) must be the one specifying the key or using the auto-detection feature
and the description for the field. The `@default` comment must follow.

See [here](./example-charts/custom-template/values.yaml) for an example.

### Spaces and Dots in keys
In the old-style comment, if a key name contains any "." or " " characters, that section of the path must be quoted in
description comments e.g.

```yaml
service:
  annotations:
    # service.annotations."external-dns.alpha.kubernetes.io/hostname" -- Hostname to be assigned to the ELB for the service
    external-dns.alpha.kubernetes.io/hostname: stupidchess.jmn23.com

configMap:
  # configMap."not real config param" -- A completely fake config parameter for a useful example
  not real config param: value
```
