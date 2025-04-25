{{- define "cmds.getValues" -}}
{{- $root := . -}}
{{- $cmdValues := dict -}}
{{- if (hasKey $root "Files") -}}
  {{- if ($root.Files.Glob (printf "cmds/%s.yaml" $root.name)) -}}
    {{- $cmdValues = $root.Files.Get (printf "cmds/%s.yaml" $root.name) | fromYaml -}}
  {{- end -}}
{{- end -}}
{{- $cmdValues | toYaml -}}
{{- end -}}

