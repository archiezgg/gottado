FROM scratch

COPY gottado /gottado

ENTRYPOINT ["/gottado"]