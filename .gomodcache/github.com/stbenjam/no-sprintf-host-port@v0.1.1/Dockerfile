FROM docker.io/library/golang:1.18 as builder

COPY / /nosprintfhostport
WORKDIR /nosprintfhostport
RUN CGO_ENABLED=0 make

FROM docker.io/library/golang:1.18
COPY --from=builder /nosprintfhostport/nosprintfhostport /usr/bin/
CMD ["nosprintfhostport"]
