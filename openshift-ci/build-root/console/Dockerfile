FROM quay.io/coreos/tectonic-console-builder:v24 AS build

RUN mkdir -p /go/src/github.com/openshift/
RUN cd /go/src/github.com/openshift && git clone https://github.com/openshift/console.git
WORKDIR /go/src/github.com/openshift/console/
RUN yarn config set ignore-engines true

RUN ./build.sh

FROM openshift/origin-base

COPY --from=build /go/src/github.com/openshift/console/frontend/public/dist /opt/bridge/static
COPY --from=build /go/src/github.com/openshift/console/bin/bridge /opt/bridge/bin/bridge
COPY --from=build /go/src/github.com/openshift/console/pkg/graphql/schema.graphql /pkg/graphql/schema.graphql

LABEL io.k8s.display-name="OpenShift Console" \
      io.k8s.description="This is a component of OpenShift Container Platform and provides a web console." \
      io.openshift.tags="openshift" \
      maintainer="Samuel Padgett <spadgett@redhat.com>"

# doesn't require a root user.
USER 1001

CMD [ "/opt/bridge/bin/bridge", "--public-dir=/opt/bridge/static" ]
