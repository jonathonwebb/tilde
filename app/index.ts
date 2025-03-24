// @ts-expect-error
import pkgJson from "../../package.json" with { type: "json" };
import { execute, parse } from "./cli.js";

await execute(
	parse(process.argv.slice(2), {
		outStream: process.stdout,
		errStream: process.stderr,
		version: pkgJson.version,
	}),
);
