import Fastify from "fastify";
import type { BaseConfig } from "./cli.js";

export type ServeConfig = BaseConfig & {
	server: {
		host: string;
		port: number;
	};
};

export async function serveCmd(config: ServeConfig) {
	const app = Fastify({ logger: { level: config.level } });

	app.get("/", async (_req, _reply) => {
		return { message: "hello, world!" };
	});

	const { host, port } = config.server;
	await app.listen({ host, port });
}
