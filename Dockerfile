FROM scratch

COPY bin/gottado /gottado
COPY templates/ templates/

EXPOSE 3000

ENTRYPOINT ["/gottado"]