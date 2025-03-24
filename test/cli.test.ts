import { type TestContext, test } from "node:test";
import { parse } from "../app/cli.js";
import { StringWritable } from "./support/utils.js";

test("parse", async (t: TestContext) => {
	t.plan(1);

	await t.test("no args", (t: TestContext) => {
		t.plan(3);

		const outStream = new StringWritable();
		const errStream = new StringWritable();

		const result = parse([], { outStream, errStream });

		t.assert.deepStrictEqual(result, { command: "abort", exitCode: 0 });
		t.assert.snapshot(outStream.data);
		t.assert.snapshot(errStream.data);
	});
});
