ARG MARIADB_VERSION=10.8.2

# Use the official NGINX image as the base image
FROM mariadb:${MARIADB_VERSION}

RUN apt update \
    && apt install --no-install-recommends -y tzdata \
    && apt clean
