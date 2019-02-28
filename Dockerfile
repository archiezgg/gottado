FROM scratch

COPY bin/gottado /gottado

EXPOSE 3000

ENTRYPOINT ["/gottado"]