
HOST=xxx
USER=xxx
curl -XPUT http://${HOST}/voice_lover_audio -u${USER} -H 'Content-Type: application/json' -d '{
  "mappings":{
    "default":{
      "properties" : {
                "albums" : {
                  "type" : "long"
                },
                "audit_status" : {
                  "type" : "long"
                },
                "cover" : {
                  "type" : "text",
                  "fields" : {
                    "keyword" : {
                      "type" : "keyword",
                      "ignore_above" : 256
                    }
                  }
                },
                "create_time" : {
                  "type" : "long"
                },
                "desc" : {
                  "type" : "text",
                  "fields" : {
                    "keyword" : {
                      "type" : "keyword",
                      "ignore_above" : 256
                    }
                  }
                },
                "has_album" : {
                  "type" : "long"
                },
                "id" : {
                  "type" : "long"
                },
                "labels" : {
                  "type" : "text",
                  "analyzer" : "ik_smart"
                },
                "op_uid" : {
                  "type" : "long"
                },
                "pub_uid" : {
                  "type" : "long"
                },
                "resource" : {
                  "type" : "text",
                  "fields" : {
                    "keyword" : {
                      "type" : "keyword",
                      "ignore_above" : 256
                    }
                  }
                },
                "seconds" : {
                  "type" : "long"
                },
                "source" : {
                  "type" : "long"
                },
                "title" : {
                  "type" : "text",
                  "analyzer" : "ik_smart"
                }
      }
    }
  }
}'