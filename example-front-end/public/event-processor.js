/**
 * Used for emitting audio samples to listeners.
 */
class EventProcessor extends AudioWorkletProcessor {
	process(inputs, outputs) {
		const data = [
			inputs.map((channels) => channels.map((channel) => Array.from(channel))),
		];
		this.port.postMessage(JSON.stringify(data));
		for (const output of outputs) {
			for (const channel of output) {
				channel.forEach((_, i) => {
					channel[i] = 0;
				});
			}
		}
		return true;
	}
}

registerProcessor("event-processor", EventProcessor);
