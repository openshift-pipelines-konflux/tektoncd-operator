ARG GO_BUILDER=brew.registry.redhat.io/rh-osbs/openshift-golang-builder:v1.23
ARG RUNTIME=registry.access.redhat.com/ubi8/ubi-minimal@sha256:e4ca1025efad7cc944fb991c716067b12e25ecbe05e6e81715a337179e417aa8

FROM $GO_BUILDER as builder

WORKDIR /go/src/github.com/tektoncd/operator
COPY upstream .
# fixme: handle patches (maybe ? probably not needed though)
# COPY patches patches/
# RUN set -e; for f in patches/*.patch; do echo ${f}; [[ -f ${f} ]] || continue; git apply ${f}; done
# ENV CHANGESET_REV=$CI_OPERATOR_UPSTREAM_COMMIT
ENV GODEBUG="http2server=0"
RUN go build -ldflags="-X 'knative.dev/pkg/changeset.rev=${CHANGESET_REV:0:7}'" -mod=vendor -o /tmp/openshift-pipelines-operator-proxy \
    ./cmd/openshift/proxy-webhook
# RUN /bin/sh -c 'echo $CI_PIPELINE_UPSTREAM_COMMIT > /tmp/HEAD'

FROM $RUNTIME

ENV OPERATOR_PROXY=/usr/local/bin/openshift-pipelines-operator-proxy \
    KO_DATA_PATH=/kodata

COPY --from=builder /tmp/openshift-pipelines-operator-proxy ${OPERATOR_PROXY}
# COPY --from=builder /tmp/HEAD ${KO_DATA_PATH}/HEAD

LABEL \
      com.redhat.component="openshift-pipelines-operator-proxy-rhel9-container" \
      name="openshift-pipelines/pipelines-operator-proxy-rhel9" \
      version="next" \
      summary="Red Hat OpenShift Pipelines Operator Proxy" \
      maintainer="pipelines-extcomm@redhat.com" \
      description="Red Hat OpenShift Pipelines Operator Proxy" \
      io.k8s.display-name="Red Hat OpenShift Pipelines Operator Proxy" \
      io.k8s.description="Red Hat OpenShift Pipelines Operator Proxy" \
      io.openshift.tags="operator,tekton,openshift,proxy"

RUN microdnf install -y shadow-utils && \
    groupadd -r -g 65532 nonroot && \
    useradd --no-log-init -r -u 65532 -g nonroot nonroot

USER 65532
ENTRYPOINT ["/usr/local/bin/openshift-pipelines-operator-proxy"]
