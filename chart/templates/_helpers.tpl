{{/*
Expand the name of the chart.
*/}}
{{- define "phasor.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
*/}}
{{- define "phasor.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "phasor.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "phasor.labels" -}}
helm.sh/chart: {{ include "phasor.chart" . }}
{{ include "phasor.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "phasor.selectorLabels" -}}
app.kubernetes.io/name: {{ include "phasor.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Backend fullname
*/}}
{{- define "phasor.backend.fullname" -}}
{{- printf "%s-backend" (include "phasor.fullname" .) }}
{{- end }}

{{/*
Backend labels
*/}}
{{- define "phasor.backend.labels" -}}
{{ include "phasor.labels" . }}
app.kubernetes.io/component: backend
{{- end }}

{{/*
Backend selector labels
*/}}
{{- define "phasor.backend.selectorLabels" -}}
{{ include "phasor.selectorLabels" . }}
app.kubernetes.io/component: backend
{{- end }}

{{/*
Frontend fullname
*/}}
{{- define "phasor.frontend.fullname" -}}
{{- printf "%s-frontend" (include "phasor.fullname" .) }}
{{- end }}

{{/*
Frontend labels
*/}}
{{- define "phasor.frontend.labels" -}}
{{ include "phasor.labels" . }}
app.kubernetes.io/component: frontend
{{- end }}

{{/*
Frontend selector labels
*/}}
{{- define "phasor.frontend.selectorLabels" -}}
{{ include "phasor.selectorLabels" . }}
app.kubernetes.io/component: frontend
{{- end }}
