FROM --platform=${TARGETPLATFORM} busybox:latest
COPY tetris /bin/tetris
ENTRYPOINT ["/bin/tetris"]
