FROM ubuntu

ARG DOCKER_PACKAGE_PATH
ENV DOCKER_PACKAGE_PATH ${DOCKER_PACKAGE_PATH}

ENV SERVICE_NAME=ddns
ENV SERVICE_HOME=/opt/${SERVICE_NAME}
ENV PATH=$PATH:${SERVICE_HOME}/bin

COPY ${DOCKER_PACKAGE_PATH} ${SERVICE_HOME}
RUN apt update
RUN apt install -y --no-install-recommends ca-certificates curl
RUN update-ca-certificates

WORKDIR ${SERVICE_HOME}/bin
ENTRYPOINT [ "ddns" ]
