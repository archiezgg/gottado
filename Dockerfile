FROM scratch

COPY bin/gottado /gottado

ENTRYPOINT ["/gottado"]