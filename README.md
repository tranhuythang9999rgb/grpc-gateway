1. Install the toolchain.

   ```shell
   go install \
       "google.golang.org/protobuf/cmd/protoc-gen-go" \
       "google.golang.org/grpc/cmd/protoc-gen-go-grpc" \
       "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway" \
       "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2"
   ```
2 gen 
``
protoc -I ./proto \
  --go_out ./proto --go_opt paths=source_relative \
  --go-grpc_out ./proto --go-grpc_opt paths=source_relative \
  --grpc-gateway_out ./proto --grpc-gateway_opt paths=source_relative \
  --openapiv2_out ./proto --openapiv2_opt logtostderr=true,allow_merge=true,merge_file_name=car \
  ./proto/api/api.proto

``
api change port
``

curl --location '127.0.0.1:8090/v1/example/echo' \
--header 'Content-Type: application/json' \
--header 'Accept: application/json' \
--data '{
  "name": "<string>"
}'

``