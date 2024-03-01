# Define function directory
ARG FUNCTION_DIR="/function"
ARG PYTHON_VERSION=3.11
ARG ISET=avx

FROM python:${PYTHON_VERSION}-slim-bullseye as build-image
# Install aws-worker-cpp build dependencies
RUN apt-get update; apt-get upgrade; apt-get install -y \
  g++ \
  make \
  cmake \
  unzip \
  libcurl4-openssl-dev

# Include global arg in this stage of the build
ARG FUNCTION_DIR
# Create function directory
RUN mkdir -p ${FUNCTION_DIR}
COPY ./poetry.lock ./pyproject.toml  ${FUNCTION_DIR}

# Install the runtime interface client
RUN pip install \
        --target ${FUNCTION_DIR} \
        awslambdaric


# Multi-stage build: grab a fresh copy of the base image
FROM python:${PYTHON_VERSION}-slim-bullseye as lambda-image
RUN pip install poetry==1.3.2 poetry-core==1.4.0 --quiet --disable-pip-version-check --no-cache-dir

# Include global arg in this stage of the build
ARG FUNCTION_DIR

# Copy in the build image dependencies
COPY --from=build-image ${FUNCTION_DIR} ${FUNCTION_DIR}

# Set working directory to function root directory
WORKDIR ${FUNCTION_DIR}

ARG ARCHITECTURE
ADD bin/aws-lambda-rie-${ARCHITECTURE} /usr/local/bin/aws-lambda-rie
RUN chmod -R 765 /usr/local/bin/aws-lambda-rie

FROM lambda-image
# Copy function code
COPY ./s  ${FUNCTION_DIR}/s
COPY ./object_storage  ${FUNCTION_DIR}/object_storage

ENTRYPOINT ["s/entry-script"]
