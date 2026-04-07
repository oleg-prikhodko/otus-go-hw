{{/*
Expand the name of the chart.
*/}}
{{- define "calendar.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "calendar.labels" -}}
helm.sh/chart: {{ include "calendar.name" . }}
app.kubernetes.io/name: {{ .Chart.Name }}
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
App selector labels
*/}}
{{- define "calendar.selectorLabels" -}}
app.kubernetes.io/name: {{ .Chart.Name }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create app name
*/}}
{{- define "calendar.appName" -}}
{{- .Values.app.name -}}
{{- end }} 