/**
 * Used for emitting audio samples to listeners.
 */
class EventProcessor extends AudioWorkletProcessor {
	process(inputs, outputs) {
		const data = [
			inputs.map((channels) => channels.map((channel) => Array.from(channel))),
		];
		this.port.postMessage(JSON.stringify(data));

		// Mute the audio, so that the audio processor does not needlessly output it
		// to the speakers.
		//
		// TODO: not sure if this is needed. We need to test this.
		//
		//    Thing to test: does an audio context not connected to the output
		//    actually end up producing any sound? And even if it does not produce
		//    anything audible, does it get picked up by WebRTC?
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
