import path from "node:path";
import { snapshot } from "node:test";

const __dirname = import.meta.dirname;
const SNAPSHOTS_DIR = "test/support/snapshots";

snapshot.setResolveSnapshotPath((testFile) => {
	if (testFile === undefined)
		throw new Error("Snapshot resolution requires a test file path");
	const testDir = path.resolve(__dirname, "..");
	const relPath = path.relative(testDir, testFile);
	return path.join(SNAPSHOTS_DIR, `${relPath}.snapshot`);
});
