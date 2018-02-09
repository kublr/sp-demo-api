FROM alpine
MAINTAINER Brian Redmond <brianisrunning@gmail.com>

ARG VCS_REF
ARG BUILD_DATE
ARG IMAGE_TAG_REF

# Metadata
LABEL org.label-schema.vcs-ref=$VCS_REF \
      org.label-schema.name="Microsmack API app" \
      org.label-schema.description="Simple golang web api for use in Kubernetes demos" \
      org.label-schema.vcs-url="https://github.com/alex-egorov/smachapi" \
      org.label-schema.build-date=$BUILD_DATE \
      org.label-schema.docker.dockerfile="Dockerfile"

ENV GIT_SHA $VCS_REF
ENV IMAGE_BUILD_DATE $BUILD_DATE
ENV IMAGE_TAG $IMAGE_TAG_REF

WORKDIR /app
ADD target/smackapi /app/

RUN chmod +x /app/smackapi

ENTRYPOINT /app/smackapi
EXPOSE 8020