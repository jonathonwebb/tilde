import type { Buffer } from "node:buffer";
import { Writable, type WritableOptions } from "node:stream";
import { StringDecoder } from "node:string_decoder";

export class StringWritable extends Writable {
	data: string;
	private _decoder: StringDecoder;

	constructor(opts?: WritableOptions) {
		super(opts);
		this.data = "";
		this._decoder = new StringDecoder(opts?.defaultEncoding);
	}

	override _write(
		chunk: Buffer | string,
		encoding: BufferEncoding | "buffer",
		callback: (error?: Error | null) => void,
	): void {
		if (encoding === "buffer") {
			// biome-ignore lint/style/noParameterAssign: style
			chunk = this._decoder.write(chunk);
		}
		this.data += chunk;
		callback();
	}
}
