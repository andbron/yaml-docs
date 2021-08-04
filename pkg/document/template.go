package document

import (
	"io/ioutil"
	"os"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig"
	log "github.com/sirupsen/logrus"

	"github.com/norwoodj/helm-docs/pkg/helm"
)

const defaultDocumentationTemplate = `
{{ template "chart.valuesSection" . }}

{{ template "helm-docs.versionFooter" . }}
`

func getNameTemplate() string {
	nameBuilder := strings.Builder{}
	nameBuilder.WriteString(`{{ define "chart.name" }}`)
	nameBuilder.WriteString("{{ .Name }}")
	nameBuilder.WriteString("{{ end }}")

	return nameBuilder.String()
}

func getHeaderTemplate() string {
	headerTemplateBuilder := strings.Builder{}
	headerTemplateBuilder.WriteString(`{{ define "chart.header" }}`)
	headerTemplateBuilder.WriteString("# {{ .Name }}\n")
	headerTemplateBuilder.WriteString("{{ end }}")

	return headerTemplateBuilder.String()
}

func getDeprecatedTemplate() string {
	deprecatedTemplateBuilder := strings.Builder{}
	deprecatedTemplateBuilder.WriteString(`{{ define "chart.deprecationWarning" }}`)
	deprecatedTemplateBuilder.WriteString("{{ if .Deprecated }}> **:exclamation: This Helm Chart is deprecated!**{{ end }}")
	deprecatedTemplateBuilder.WriteString("{{ end }}")

	return deprecatedTemplateBuilder.String()
}

func getVersionTemplates() string {
	versionBuilder := strings.Builder{}
	versionBuilder.WriteString(`{{ define "chart.version" }}{{ .Version }}{{ end }}\n`)
	versionBuilder.WriteString(`{{ define "chart.versionBadge" }}`)
	versionBuilder.WriteString(`![Version: {{ .Version }}](https://img.shields.io/badge/Version-{{ .Version | replace "-" "--" }}-informational?style=flat-square) `)
	versionBuilder.WriteString("{{ end }}")

	return versionBuilder.String()
}

func getTypeTemplate() string {
	typeBuilder := strings.Builder{}
	typeBuilder.WriteString(`{{ define "chart.type" }}{{ .Type }}{{ end }}\n`)
	typeBuilder.WriteString(`{{ define "chart.typeBadge" }}`)
	typeBuilder.WriteString("{{ if .Type }}![Type: {{ .Type }}](https://img.shields.io/badge/Type-{{ .Type }}-informational?style=flat-square) {{ end }}")
	typeBuilder.WriteString("{{ end }}")

	return typeBuilder.String()
}

func getAppVersionTemplate() string {
	appVersionBuilder := strings.Builder{}
	appVersionBuilder.WriteString(`{{ define "chart.appVersion" }}{{ .AppVersion }}{{ end }}\n`)
	appVersionBuilder.WriteString(`{{ define "chart.appVersionBadge" }}`)
	appVersionBuilder.WriteString(`{{ if .AppVersion }}![AppVersion: {{ .AppVersion }}](https://img.shields.io/badge/AppVersion-{{ .AppVersion | replace "-" "--" }}-informational?style=flat-square) {{ end }}`)
	appVersionBuilder.WriteString("{{ end }}")

	return appVersionBuilder.String()
}

func getBadgesTemplates() string {
	badgeBuilder := strings.Builder{}
	badgeBuilder.WriteString(`{{ define "chart.badgesSection" }}`)
	badgeBuilder.WriteString(`{{ template "chart.versionBadge" . }}{{ template "chart.typeBadge" . }}{{ template "chart.appVersionBadge" . }}`)
	badgeBuilder.WriteString("{{ end }}")

	return badgeBuilder.String()
}

func getDescriptionTemplate() string {
	descriptionBuilder := strings.Builder{}
	descriptionBuilder.WriteString(`{{ define "chart.description" }}`)
	descriptionBuilder.WriteString("{{ if .Description }}{{ .Description }}{{ end }}")
	descriptionBuilder.WriteString("{{ end }}")

	return descriptionBuilder.String()
}

func getHomepageTemplate() string {
	homepageBuilder := strings.Builder{}
	homepageBuilder.WriteString(`{{ define "chart.homepage" }}{{ .Home }}{{ end }}\n`)
	homepageBuilder.WriteString(`{{ define "chart.homepageLine" }}`)
	homepageBuilder.WriteString("{{ if .Home }}**Homepage:** <{{ .Home }}>{{ end }}")
	homepageBuilder.WriteString("{{ end }}")

	return homepageBuilder.String()
}

func getMaintainersTemplate() string {
	maintainerBuilder := strings.Builder{}
	maintainerBuilder.WriteString(`{{ define "chart.maintainersHeader" }}## Maintainers{{ end }}`)

	maintainerBuilder.WriteString(`{{ define "chart.maintainersTable" }}`)
	maintainerBuilder.WriteString("| Name | Email | Url |\n")
	maintainerBuilder.WriteString("| ---- | ------ | --- |\n")
	maintainerBuilder.WriteString("  {{- range .Maintainers }}")
	maintainerBuilder.WriteString("\n| {{ .Name }} | {{ .Email }} | {{ .Url }} |")
	maintainerBuilder.WriteString("  {{- end }}")
	maintainerBuilder.WriteString("{{ end }}")

	maintainerBuilder.WriteString(`{{ define "chart.maintainersSection" }}`)
	maintainerBuilder.WriteString("{{ if .Maintainers }}")
	maintainerBuilder.WriteString(`{{ template "chart.maintainersHeader" . }}`)
	maintainerBuilder.WriteString("\n\n")
	maintainerBuilder.WriteString(`{{ template "chart.maintainersTable" . }}`)
	maintainerBuilder.WriteString("{{ end }}")
	maintainerBuilder.WriteString("{{ end }}")

	return maintainerBuilder.String()
}

func getSourceLinkTemplates() string {
	sourceLinkBuilder := strings.Builder{}
	sourceLinkBuilder.WriteString(`{{ define "chart.sourcesHeader" }}## Source Code{{ end}}`)

	sourceLinkBuilder.WriteString(`{{ define "chart.sourcesList" }}`)
	sourceLinkBuilder.WriteString("{{- range .Sources }}")
	sourceLinkBuilder.WriteString("\n* <{{ . }}>")
	sourceLinkBuilder.WriteString("{{- end }}")
	sourceLinkBuilder.WriteString("{{ end }}")

	sourceLinkBuilder.WriteString(`{{ define "chart.sourcesSection" }}`)
	sourceLinkBuilder.WriteString("{{ if .Sources }}")
	sourceLinkBuilder.WriteString(`{{ template "chart.sourcesHeader" . }}`)
	sourceLinkBuilder.WriteString("\n")
	sourceLinkBuilder.WriteString(`{{ template "chart.sourcesList" . }}`)
	sourceLinkBuilder.WriteString("{{ end }}")
	sourceLinkBuilder.WriteString("{{ end }}")

	return sourceLinkBuilder.String()
}

func getRequirementsTableTemplates() string {
	requirementsSectionBuilder := strings.Builder{}
	requirementsSectionBuilder.WriteString(`{{ define "chart.requirementsHeader" }}## Requirements{{ end }}`)

	requirementsSectionBuilder.WriteString(`{{ define "chart.kubeVersion" }}{{ .KubeVersion }}{{ end }}\n`)
	requirementsSectionBuilder.WriteString(`{{ define "chart.kubeVersionLine" }}`)
	requirementsSectionBuilder.WriteString("{{ if .KubeVersion }}Kubernetes: `{{ .KubeVersion }}`{{ end }}")
	requirementsSectionBuilder.WriteString("{{ end }}")

	requirementsSectionBuilder.WriteString(`{{ define "chart.requirementsTable" }}`)
	requirementsSectionBuilder.WriteString("| Repository | Name | Version |\n")
	requirementsSectionBuilder.WriteString("|------------|------|---------|")
	requirementsSectionBuilder.WriteString("  {{- range .Dependencies }}")
	requirementsSectionBuilder.WriteString("    {{- if .Alias }}")
	requirementsSectionBuilder.WriteString("\n| {{ .Repository }} | {{ .Alias }}({{ .Name }}) | {{ .Version }} |")
	requirementsSectionBuilder.WriteString("    {{- else }}")
	requirementsSectionBuilder.WriteString("\n| {{ .Repository }} | {{ .Name }} | {{ .Version }} |")
	requirementsSectionBuilder.WriteString("    {{- end }}")
	requirementsSectionBuilder.WriteString("  {{- end }}")
	requirementsSectionBuilder.WriteString("{{ end }}")

	requirementsSectionBuilder.WriteString(`{{ define "chart.requirementsSection" }}`)
	requirementsSectionBuilder.WriteString("{{ if or .Dependencies .KubeVersion }}")
	requirementsSectionBuilder.WriteString(`{{ template "chart.requirementsHeader" . }}`)
	requirementsSectionBuilder.WriteString("\n\n")
	requirementsSectionBuilder.WriteString("{{ if .KubeVersion }}")
	requirementsSectionBuilder.WriteString(`{{ template "chart.kubeVersionLine" . }}`)
	requirementsSectionBuilder.WriteString("\n\n")
	requirementsSectionBuilder.WriteString("{{ end }}")
	requirementsSectionBuilder.WriteString("{{ if .Dependencies }}")
	requirementsSectionBuilder.WriteString(`{{ template "chart.requirementsTable" . }}`)
	requirementsSectionBuilder.WriteString("{{ end }}")
	requirementsSectionBuilder.WriteString("{{ end }}")
	requirementsSectionBuilder.WriteString("{{ end }}")

	return requirementsSectionBuilder.String()
}

func getValuesTableTemplates() string {
	valuesSectionBuilder := strings.Builder{}
	valuesSectionBuilder.WriteString(`{{ define "chart.valuesHeader" }}## Values{{ end }}`)

	valuesSectionBuilder.WriteString(`{{ define "chart.valuesTable" }}`)
	valuesSectionBuilder.WriteString("| Key | Type | Default | Description |\n")
	valuesSectionBuilder.WriteString("|-----|------|---------|-------------|\n")
	valuesSectionBuilder.WriteString("  {{- range .Values }}")
	valuesSectionBuilder.WriteString("\n| {{ .Key }} | {{ .Type }} | {{ if .Default }}{{ .Default }}{{ else }}{{ .AutoDefault }}{{ end }} | {{ if .Description }}{{ .Description }}{{ else }}{{ .AutoDescription }}{{ end }} |")
	valuesSectionBuilder.WriteString("  {{- end }}")
	valuesSectionBuilder.WriteString("{{ end }}")

	valuesSectionBuilder.WriteString(`{{ define "chart.valuesSection" }}`)
	valuesSectionBuilder.WriteString("{{ if .Values }}")
	valuesSectionBuilder.WriteString(`{{ template "chart.valuesHeader" . }}`)
	valuesSectionBuilder.WriteString("\n\n")
	valuesSectionBuilder.WriteString(`{{ template "chart.valuesTable" . }}`)
	valuesSectionBuilder.WriteString("{{ end }}")
	valuesSectionBuilder.WriteString("{{ end }}")

	return valuesSectionBuilder.String()
}

func getHelmDocsVersionTemplates() string {
	versionSectionBuilder := strings.Builder{}
	versionSectionBuilder.WriteString(`{{ define "helm-docs.version" }}{{ if .HelmDocsVersion }}{{ .HelmDocsVersion }}{{ end }}{{ end }}`)
	versionSectionBuilder.WriteString(`{{ define "helm-docs.versionFooter" }}`)
	versionSectionBuilder.WriteString("{{ if .HelmDocsVersion }}\n")
	versionSectionBuilder.WriteString("----------------------------------------------\n")
	versionSectionBuilder.WriteString("Autogenerated from chart metadata using [helm-docs v{{ .HelmDocsVersion }}](https://github.com/norwoodj/helm-docs/releases/v{{ .HelmDocsVersion }})")
	versionSectionBuilder.WriteString("{{ end }}")
	versionSectionBuilder.WriteString("{{ end }}")

	return versionSectionBuilder.String()
}

func getDocumentationTemplate(templateFiles []string) (string, error) {
	templateFilesForChart := make([]string, 0)

	var templateNotFound bool

	for _, templateFile := range templateFiles {
		var fullTemplatePath string

		fullTemplatePath = templateFile

		if _, err := os.Stat(fullTemplatePath); os.IsNotExist(err) {
			log.Debugf("Did not find template file %s, using default template", templateFile)

			templateNotFound = true
			continue
		}

		templateFilesForChart = append(templateFilesForChart, fullTemplatePath)
	}

	log.Debugf("Using template files %s", templateFiles)
	allTemplateContents := make([]byte, 0)
	for _, templateFileForChart := range templateFilesForChart {
		templateContents, err := ioutil.ReadFile(templateFileForChart)
		if err != nil {
			return "", err
		}
		allTemplateContents = append(allTemplateContents, templateContents...)
	}

	if templateNotFound {
		allTemplateContents = append(allTemplateContents, []byte(defaultDocumentationTemplate)...)
	}

	return string(allTemplateContents), nil
}

func getDocumentationTemplates(templateFiles []string) ([]string, error) {
	documentationTemplate, err := getDocumentationTemplate(templateFiles)

	if err != nil {
		log.Errorf("Failed to read documentation templates %s: %s", templateFiles, err)
		return nil, err
	}

	return []string{
		getNameTemplate(),
		getHeaderTemplate(),
		getDeprecatedTemplate(),
		getAppVersionTemplate(),
		getBadgesTemplates(),
		getDescriptionTemplate(),
		getVersionTemplates(),
		getTypeTemplate(),
		getSourceLinkTemplates(),
		getRequirementsTableTemplates(),
		getValuesTableTemplates(),
		getHomepageTemplate(),
		getMaintainersTemplate(),
		getHelmDocsVersionTemplates(),
		documentationTemplate,
	}, nil
}

func newChartDocumentationTemplate(chartDocumentationInfo helm.ChartDocumentationInfo, templateFiles []string) (*template.Template, error) {
	documentationTemplate := template.New(chartDocumentationInfo.ChartDirectory)
	documentationTemplate.Funcs(sprig.TxtFuncMap())
	goTemplateList, err := getDocumentationTemplates(templateFiles)

	if err != nil {
		return nil, err
	}

	for _, t := range goTemplateList {
		_, err := documentationTemplate.Parse(t)

		if err != nil {
			return nil, err
		}
	}

	return documentationTemplate, nil
}
