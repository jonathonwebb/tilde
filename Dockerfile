ARG NODE_VERSION=22.14.0
FROM node:${NODE_VERSION}-alpine AS base

# handle signals & reap processes when running node via sh exec:
RUN apk add --no-cache tini && \
		ln -s /sbin/tini /bin/tini
ENTRYPOINT ["/bin/tini", "--"]

ENV PNPM_HOME="/opt/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
RUN corepack enable

USER node
WORKDIR /opt/tilde
COPY --chown=node:node package.json pnpm-lock.yaml ./

FROM base AS dev

# npm dependency cache:
RUN --mount=type=cache,id=pnpm,target=/opt/pnpm/store,sharing=locked \
	pnpm install --frozen-lockfile

COPY --chown=node:node . .
RUN ./bin/build

FROM base AS prod

# install only production deps via cache:
RUN --mount=type=cache,id=pnpm,target=/opt/pnpm/store,sharing=locked \
	pnpm install --prod --frozen-lockfile

COPY --chown=node:node . .
# ensure changes to out/ aren't skipped in prod stage:
RUN --mount=from=dev,source=/opt/tilde/out,target=/tmp/out cp -r /tmp/out ./
