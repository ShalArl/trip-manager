{{- define "common.service" -}}
apiVersion: v1
kind: Service
metadata:
  name: {{ include "common.fullname" . }}
  namespace: {{ .Release.Namespace }}
  labels: {{- include "common.labels" . | nindent 4 }}
  {{- if .Values.service.healthPath }}
  annotations:
    cloud.google.com/backend-config: {{ printf `{"default": "%s-backend-config"}` (include "common.fullname" .) | quote }}
  {{- end }}
spec:
  type: ClusterIP
  ports:
    - port: {{ .Values.service.port }}
      targetPort: http
      protocol: TCP
      name: http
  selector: {{- include "common.selectorLabels" . | nindent 4 }}
{{- end }}

{{- define "common.backendConfig" -}}
{{- if .Values.service.healthPath }}
---
apiVersion: cloud.google.com/v1
kind: BackendConfig
metadata:
  name: {{ include "common.fullname" . }}-backend-config
  namespace: {{ .Release.Namespace }}
  labels: {{- include "common.labels" . | nindent 4 }}
spec:
  healthCheck:
    checkIntervalSec: 15
    timeoutSec: 5
    healthyThreshold: 1
    unhealthyThreshold: 2
    type: HTTP
    requestPath: {{ .Values.service.healthPath }}
    port: {{ .Values.service.port }}
{{- end }}
{{- end }}