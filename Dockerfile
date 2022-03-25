FROM golang:latest

COPY . /go/src/app/.

RUN mkdir /app && cd /go/src/app && go get -d -v -buildvcs=false ./... && go install -v -buildvcs=false ./... && go build -o /app/inference_server cmd/inference/ur_inference_server.go && cp -R /go/src/app/trained /app/trained && rm -r /go/src/app/*

EXPOSE 8090

WORKDIR /app
CMD ["/app/inference_server"]