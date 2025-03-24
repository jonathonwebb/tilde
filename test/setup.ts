import path from "node:path";
import { snapshot } from "node:test";

snapshot.setResolveSnapshotPath((testFile) => {
	if (testFile === undefined)
		throw new Error("Snapshot resolution requires a test file path");
	const relPath = path.relative(import.meta.dirname, testFile);
	return path.join("test", `${relPath}.snapshot`);
});
