# cobbled together from
# https://github.com/felixbuenemann/vips-alpine/blob/master/Dockerfile
# https://github.com/mikestead/docker-imaginary-alpine/blob/master/Dockerfile

FROM golang:alpine as builder

ARG VIPS_VERSION=8.6.4

ADD . /go-iiif

ENV VIPS_DIR=/vips
ENV PKG_CONFIG_PATH=${VIPS_DIR}/lib/pkgconfig:$PKG_CONFIG_PATH

RUN wget -O- https://github.com/jcupitt/libvips/releases/download/v${VIPS_VERSION}/vips-${VIPS_VERSION}.tar.gz | tar xzC /tmp \
    && apk update \
    && apk upgrade \

    && apk add \
    make libc-dev gcc \
    zlib libxml2 glib gobject-introspection \
    libjpeg-turbo libexif lcms2 fftw giflib libpng \
    libwebp orc tiff poppler-glib librsvg libgsf openexr \

    && apk add --virtual vips-dependencies build-base \
    zlib-dev libxml2-dev glib-dev gobject-introspection-dev \
    libjpeg-turbo-dev libexif-dev lcms2-dev fftw-dev giflib-dev libpng-dev \
    libwebp-dev orc-dev tiff-dev poppler-dev librsvg-dev libgsf-dev openexr-dev \

    && cd /tmp/vips-${VIPS_VERSION} \
    && ./configure --prefix=${VIPS_DIR} \
                   --disable-static \
		   --without-python \
                   --disable-dependency-tracking \
                   --enable-silent-rules \
    && make -s install-strip \

    && cd /go-iiif \
    && make bin

FROM alpine

COPY --from=builder /vips/lib/ /usr/local/lib
COPY --from=builder /go-iiif/bin/iiif-server /bin/iiif-server

EXPOSE 8080

# RUN ME...
