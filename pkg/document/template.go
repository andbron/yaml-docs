package document

import (
	"io/ioutil"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig"
	log "github.com/sirupsen/logrus"

)

const defaultDocumentationTemplate = `
{{ template "docs.valuesSection" . }}

{{ template "yaml-docs.versionFooter" . }}
`

// String providing three templates: value
func getValuesTableTemplates() string {
	valuesSectionBuilder := strings.Builder{}
	valuesSectionBuilder.WriteString(`{{ define "docs.valuesHeader" }}## Values{{ end }}`)

	valuesSectionBuilder.WriteString(`{{ define "docs.valuesTable" }}`)
	valuesSectionBuilder.WriteString("| Key | Type | Description |\n")
	valuesSectionBuilder.WriteString("|-----|------|-------------|\n")
	valuesSectionBuilder.WriteString("  {{- range .Values }}")
	valuesSectionBuilder.WriteString("\n| {{ .Key }} | {{ .Type }} | {{ .Description }} |")
	valuesSectionBuilder.WriteString("  {{- end }}")
	valuesSectionBuilder.WriteString("{{ end }}")

	valuesSectionBuilder.WriteString(`{{ define "docs.valuesSection" }}`)
	valuesSectionBuilder.WriteString("{{ if .Values }}")
	valuesSectionBuilder.WriteString(`{{ template "docs.valuesHeader" . }}`)
	valuesSectionBuilder.WriteString("\n\n")
	valuesSectionBuilder.WriteString(`{{ template "docs.valuesTable" . }}`)
	valuesSectionBuilder.WriteString("{{ end }}")
	valuesSectionBuilder.WriteString("{{ end }}")

	return valuesSectionBuilder.String()
}

func getYamlDocsVersionTemplates() string {
	versionSectionBuilder := strings.Builder{}
	versionSectionBuilder.WriteString(`{{ define "yaml-docs.version" }}{{ if .YamlDocsVersion }}{{ .YamlDocsVersion }}{{ end }}{{ end }}`)
	versionSectionBuilder.WriteString(`{{ define "yaml-docs.versionFooter" }}`)
	versionSectionBuilder.WriteString("{{ if .YamlDocsVersion }}\n")
	versionSectionBuilder.WriteString("----------------------------------------------\n")
	versionSectionBuilder.WriteString("Autogenerated using [yaml-docs v{{ .YamlDocsVersion }}](https://https://github.com/andbron/yaml-docs/releases/v{{ .YamlDocsVersion }})")
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
		getValuesTableTemplates(),
		getYamlDocsVersionTemplates(),
		documentationTemplate,
	}, nil
}


func newDocumentationTemplate(templateFiles []string) (*template.Template, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	documentationTemplate := template.New(path.Base(cwd))
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
