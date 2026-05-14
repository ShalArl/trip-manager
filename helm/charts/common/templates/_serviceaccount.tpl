{{- define "common.serviceaccount" -}}
apiVersion: v1
kind: ServiceAccount
metadata:
  name: {{ include "common.fullname" . }}
  namespace: {{ .Release.Namespace }}
  labels: {{ include "common.labels" . | nindent 4 }}
  {{- if .Values.serviceAccount.gcpServiceAccount }}
  annotations:
    iam.gke.io/gcp-service-account: {{ .Values.serviceAccount.gcpServiceAccount }}
  {{- end }}
{{- end }}