# Authlet
Authlet is local authentication service for iothub

# Using grpc
protoc -I authlet/ authlet/authlet.proto --go_out=plugins=grpc:authlet
