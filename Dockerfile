ARG NODE_VERSION=22.14.0
FROM node:${NODE_VERSION}-alpine AS base

RUN apk add --no-cache tini
ENTRYPOINT ["/sbin/tini", "--"]

ENV PNPM_HOME="/opt/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
RUN corepack enable

USER node
WORKDIR /opt/tilde
COPY --chown=node:node package.json pnpm-lock.yaml ./

FROM base AS dev

RUN --mount=type=cache,id=pnpm,target=/opt/pnpm/store,sharing=locked \
	pnpm install --frozen-lockfile

COPY --chown=node:node . .
RUN ./bin/build

FROM base AS prod

RUN --mount=type=cache,id=pnpm,target=/opt/pnpm/store,sharing=locked \
	pnpm install --prod --frozen-lockfile

COPY --chown=node:node . .
COPY --from=dev /opt/tilde/out ./out
