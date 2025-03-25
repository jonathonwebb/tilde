import fastifyView from "@fastify/view";
import Fastify from "fastify";
import nunjucks from "nunjucks";
import type { BaseConfig } from "./cli.js";

export type ServeConfig = BaseConfig & {
	server: {
		host: string;
		port: number;
	};
};

export async function serveCmd(config: ServeConfig) {
	const app = Fastify({
		logger: {
			level: config.level,
			...(config.format === "pretty"
				? {
						transport: { target: "pino-pretty", options: { singleLine: true } },
					}
				: {}),
		},
	});

	await app.register(fastifyView, {
		engine: { nunjucks },
		production: config.env === "production",
		defaultContext: {},
		root: "templates",
		options: {
			onConfigure: (_env: nunjucks.Environment) => {},
		},
	});

	app.get("/", async (_req, reply) => {
		return reply.viewAsync("home.html", { message: "hello, world!" });
	});

	app.get("/health/ready", async (_req, reply) => {
		return reply.send({ status: "ready" });
	});

	const { host, port } = config.server;
	await app.listen({ host, port });
}
