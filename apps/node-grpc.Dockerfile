# Copyright 2025 Simon Emms <simon@simonemms.com>
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

FROM node:lts AS dev
ARG APP
ARG GIT_COMMIT
ARG VERSION
ARG GRPC_HEALTH_PROBE_VERSION=v0.4.37
RUN wget -qO /bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 \
  && chmod +x /bin/grpc_health_probe
USER node
WORKDIR /home/node/root
ENV GIT_COMMIT="${GIT_COMMIT}"
ENV VERSION="${VERSION}"
COPY --chown=node:node . .
WORKDIR /home/node/root/apps/$APP
ENV PORT=3000
ENV GRPC_HEALTH_PROBE_VERSION="${GRPC_HEALTH_PROBE_VERSION}"
ENV PROTO_PATH=/home/node/root/proto
ENV LOGGER_COLORS_ENABLED=true
ENV LOGGER_JSON_ENABLED=false
ENV LOGGER_LEVEL=verbose
EXPOSE 3000
CMD [ "npm", "run", "start:dev" ]

FROM node:lts-alpine AS builder
ARG APP
COPY --from=dev /bin/grpc_health_probe /bin/grpc_health_probe
USER node
ENV PROTO_PATH=/home/node/root/proto
WORKDIR /home/node/root
COPY --from=dev /home/node/root .
WORKDIR /home/node/root/apps/$APP
RUN npm ci \
  && npm run build

FROM node:lts-alpine
ARG APP
ARG GIT_COMMIT
ARG VERSION
WORKDIR /opt/app
ENV GIT_COMMIT="${GIT_COMMIT}"
ENV PROTO_PATH=/opt/app/proto
ENV VERSION="${VERSION}"
ENV SERVER_PORT=3000
COPY --from=builder /home/node/root/apps/$APP/dist dist
COPY --from=builder /home/node/root/proto /opt/app/proto
COPY --from=builder /home/node/root/apps/$APP/node_modules node_modules
COPY --from=builder /home/node/root/apps/$APP/package.json package.json
COPY --from=builder /home/node/root/apps/$APP/package-lock.json package-lock.json
COPY --from=dev /bin/grpc_health_probe /bin/grpc_health_probe
RUN npm prune --omit=dev \
  && npm rebuild \
  && npm dedupe \
  && npm version ${VERSION} --no-git-tag-version --allow-same-version || true
USER node
EXPOSE 3000
CMD [ "npm", "run", "start:prod" ]
