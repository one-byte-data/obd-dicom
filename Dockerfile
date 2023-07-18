FROM golang:alpine AS builder

ARG TARGETOS
ARG TARGETARCH

RUN apk add gcc g++

WORKDIR /workspace

COPY . /workspace/

RUN echo ${TARGETARCH}
RUN echo ${TARGETOS}

RUN go build -v -a -o build/${TARGETOS}/${TARGETARCH}/obd-dicom cmd/obd-dicom/main.go

RUN go build -v -a -o build/${TARGETOS}/${TARGETARCH}/compare cmd/compare/main.go

FROM alpine:latest

ARG TARGETOS
ARG TARGETARCH

RUN mkdir /app

WORKDIR /app

COPY --from=builder /workspace/build/${TARGETOS}/${TARGETARCH}/compare /app/

COPY --from=builder /workspace/build/${TARGETOS}/${TARGETARCH}/obd-dicom /app/

ENTRYPOINT [ "/app/odb-dicom", "-scp", "-calledae", "DICOM_SCP", "-port", "1040", "-datastore", "/datastore" ]
