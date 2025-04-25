{{- define "rpcs.getValues" -}}
{{- $root := . -}}
{{- $rpcValues := dict -}}
{{- if (hasKey $root "Files") -}}
  {{- if ($root.Files.Glob (printf "rpcs/%s.yaml" $root.name)) -}}
    {{- $rpcValues = $root.Files.Get (printf "rpcs/%s.yaml" $root.name) | fromYaml -}}
  {{- end -}}
{{- end -}}
{{- $rpcValues | toYaml -}}
{{- end -}}

