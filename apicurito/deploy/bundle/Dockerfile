FROM scratch

LABEL operators.operatorframework.io.bundle.mediatype.v1=registry+v1
LABEL operators.operatorframework.io.bundle.manifests.v1=manifests/
LABEL operators.operatorframework.io.bundle.metadata.v1=metadata/
LABEL operators.operatorframework.io.bundle.package.v1=apicurito
LABEL operators.operatorframework.io.bundle.channels.v1=apicurito-7.8.x
LABEL operators.operatorframework.io.bundle.channel.default.v1=apicurito-7.8.x
LABEL com.redhat.delivery.operator.bundle=true
LABEL com.redhat.openshift.versions="v4.5,v4.6"

COPY manifests /manifests/
COPY metadata/annotations.yaml /metadata/annotations.yaml

LABEL name="fuse7/fuse-online-operator-metadata" \
      version="7.8.0" \
      maintainer="Otavio Piske <opiske@redhat.com>" \
      summary="Operator which manages the lifecycle of the Apicurito application." \
      description="Operator which manages the lifecycle of the Apicurito application." \
      com.redhat.component="apicurito-operator-metadata-container" \
      io.k8s.description="Operator which manages the lifecycle of the Apicurito application." \
      io.k8s.display-name="Red Hat Apicurito Operator" \
      io.openshift.tags="fuse"
