@echo off
"h:\real_estate\temp_protoc\bin\protoc.exe" --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/recommendation/recommendation.proto
echo Da bien dich xong!
pause