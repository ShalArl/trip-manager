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

{{- define "common.backendPolicy" -}}
{{- if .Values.service.healthPath }}
---
apiVersion: networking.gke.io/v1
kind: HealthCheckPolicy
metadata:
  name: {{ include "common.fullname" . }}-hc-policy
  namespace: {{ .Release.Namespace }}
  labels: {{- include "common.labels" . | nindent 4 }}
spec:
  default:
    config:
      type: HTTP
      httpHealthCheck:
        portSpecification: USE_FIXED_PORT
        port: {{ .Values.service.port }}
        requestPath: {{ .Values.service.healthPath }}
  targetRef:
    group: ""
    kind: Service
    name: {{ include "common.fullname" . }}
{{- end }}
{{- end }}