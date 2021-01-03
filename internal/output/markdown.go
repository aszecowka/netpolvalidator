package output

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"text/template"

	"github.com/aszecowka/netpolvalidator/internal/model"
)

var md = `
# Network Policy Report

## Violations

Number of violations: {{ len .Violations }}

{{ if .Violations }}
| Namespace | Network Policy Name | Type | Message |
|-----------|---------------------|------|---------|
{{ end }}
{{- range .Violations -}}
| {{.Namespace}} | {{.NetworkPolicyName}} | {{.Type}} | {{.Message}} |
{{- end -}}
`

type Markdown struct{}

type Data struct {
	State      model.ClusterState
	Violations []model.Violation
}

func NewMarkdown() *Markdown {
	return &Markdown{}
}

func (m *Markdown) Generate(ctx context.Context, state model.ClusterState, violations []model.Violation) (io.Writer, error) {
	tpl, err := template.New("net_pol_report").Parse(md)
	if err != nil {
		return nil, fmt.Errorf("while parsing markdown report: %w", err)
	}
	buf := bytes.Buffer{}
	if err := tpl.Execute(&buf, Data{State: state, Violations: violations}); err != nil {
		return nil, fmt.Errorf("while generating markdown report: %w", err)
	}
	return &buf, nil
}
