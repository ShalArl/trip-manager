{{- define "common.service" -}}
apiVersion: v1
kind: Service
metadata:
  name: {{ include "common.fullname" . }}
  namespace: {{ .Release.Namespace }}
  labels: {{- include "common.labels" . | nindent 4 }}
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
        port: {{ .Values.service.port }}
        requestPath: {{ .Values.service.healthPath }}
  targetRef:
    group: ""
    kind: Service
    name: {{ include "common.fullname" . }}
{{- end }}
{{- end }}