FROM registry.access.redhat.com/ubi7-dev-preview/ubi-minimal:7.6

ENV OPERATOR=/usr/local/bin/apicurito \
    USER_UID=1001 \
    USER_NAME=apicurito

# install operator binary
COPY build/_output/bin/apicurito ${OPERATOR}
COPY build/conf/config.yaml /conf/config.yaml

COPY build/bin /usr/local/bin
RUN  /usr/local/bin/user_setup

ENTRYPOINT ["/usr/local/bin/entrypoint", "run"]

USER ${USER_UID}
