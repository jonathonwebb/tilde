import {
	Command,
	CommanderError,
	InvalidOptionArgumentError,
	Option,
} from "commander";
import { type Level, levels } from "pino";
import { type ServeConfig, serveCmd } from "./serve.js";

type ParseOptions = {
	outStream?: NodeJS.WritableStream;
	errStream?: NodeJS.WritableStream;
	version?: string;
};

export type BaseConfig = {
	env: string;
	level: Level;
};

type ParseResult =
	| { command: "serve"; config: ServeConfig }
	| { command: "abort"; exitCode: number };

export async function execute(result: ParseResult) {
	switch (result.command) {
		case "serve":
			serveCmd(result.config);
			break;
		case "abort":
			process.exit(result.exitCode);
	}
}

export function parse(
	args: string[],
	{ outStream, errStream, version }: ParseOptions,
): ParseResult {
	let result: ParseResult | null = null;

	try {
		const program = new Command("tilde")
			.configureOutput({
				writeOut: (s) => outStream?.write(s),
				writeErr: (s) => errStream?.write(s),
			})
			.configureHelp({
				styleTitle: (s) => s.toLowerCase(),
				showGlobalOptions: true,
			})
			.addOption(
				new Option("-e, --env <env>", "application environment")
					.default("production")
					.env("MDW_ENV"),
			)
			.addOption(
				new Option("-l, --level <level>", "log level")
					.argParser(parseLevelOptions)
					.default("info")
					.env("MDW_LEVEL"),
			)
			.version(
				`tilde v${version || "(unknown)"}`,
				"-V, --version",
				"show version info and exit",
			)
			.exitOverride();

		const serveCmd = new Command("serve")
			.copyInheritedSettings(program)
			.addOption(
				new Option("-a, --host <host>", "listener host")
					.default("localhost")
					.env("MDW_HOST"),
			)
			.addOption(
				new Option("-p, --port <number>", "listener port")
					.default(0, "ephemeral")
					.env("MDW_PORT"),
			)
			.action(function (this: Command) {
				const opts = this.optsWithGlobals();
				result = {
					command: "serve",
					config: {
						env: opts["env"],
						level: opts["level"],
						server: {
							host: opts["host"],
							port: opts["port"],
						},
					},
				};
			});

		program.addCommand(serveCmd);

		program.parse(args, { from: "user" });

		if (result === null)
			throw new Error("Expected parse result not to be null");

		return result;
	} catch (err) {
		if (err instanceof CommanderError) {
			switch (err.code) {
				case "commander.help":
				case "commander.helpDisplayed":
				case "commander.version":
					return { command: "abort", exitCode: 0 };
				default:
					return { command: "abort", exitCode: 2 };
			}
		}
		throw err;
	}
}
function parseLevelOptions(val: string, _prev: Level): Level {
	const lowerValue = val.toLowerCase();
	if (!isLevel(lowerValue))
		throw new InvalidOptionArgumentError(
			`Allowed choices are ${Object.values(levels.labels).join(", ")}.`,
		);

	return lowerValue;
}

function isLevel(s: string): s is Level {
	return Object.values(levels.labels).includes(s);
}
