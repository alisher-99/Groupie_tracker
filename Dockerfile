FROM golang:alpine AS appbuild
WORKDIR /1st-stage
COPY . /1st-stage/
RUN go build -o main cmd/main.go

FROM alpine
LABEL maintainer="alisher-99"
WORKDIR /2nd-stage
COPY --from=appbuild /1st-stage/main /2nd-stage/
COPY --from=appbuild /1st-stage/ui/ /2nd-stage/ui
COPY --from=appbuild /1st-stage/config/config.json /2nd-stage/config/
CMD [ "./main" ]