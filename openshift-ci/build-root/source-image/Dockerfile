FROM src

COPY oc /usr/bin/oc
RUN yum -y install make wget gcc git httpd-tools dnf
RUN dnf -y install npm
RUN npm --version
RUN npm i -g xunit-viewer
RUN xunit-viewer --help
