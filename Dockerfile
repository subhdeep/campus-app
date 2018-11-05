FROM golang:latest AS builder

# Download and install the latest release of dep
RUN wget https://github.com/golang/dep/releases/download/v0.5.0/dep-linux-amd64 -O /usr/bin/dep
RUN chmod +x /usr/bin/dep

# Copy the code from the host and compile it
WORKDIR $GOPATH/src/github.com/subhdeep/campus-app
COPY Gopkg.toml Gopkg.lock ./
RUN dep ensure -vendor-only -v
COPY . ./
ADD ./config.yaml /config.yaml
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o /app .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /app /config.yaml ./
ENTRYPOINT ["./app"]
