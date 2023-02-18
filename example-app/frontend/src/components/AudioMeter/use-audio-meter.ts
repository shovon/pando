import { useCallback, useEffect, useRef } from "react";
import {
	arrayOf,
	chain,
	InferType,
	number,
	transform,
	validate,
} from "../../validator";

// TODO: I suspect a lot of this code will need to be modularized even more

const audioStreamValidator = chain(
	transform((v) => JSON.parse(v)),
	arrayOf(arrayOf(arrayOf(arrayOf(number()))))
);

type AudioStreamValue = InferType<typeof audioStreamValidator>;
type AudioStreamListener = (value: AudioStreamValue) => void;

export function useAudioStreamMeter(stream: MediaStream) {
	const listeners = useRef<AudioStreamListener[]>([]);
	const callback = useCallback((value: AudioStreamValue) => {
		for (const listener of listeners.current) {
			listener(value);
		}
	}, []);

	useEffect(() => {
		let audioContext: AudioContext | undefined;
		let workletNode: AudioWorkletNode | undefined;

		Promise.resolve().then(async () => {
			audioContext = new AudioContext();
			await audioContext.audioWorklet.addModule(
				`/event-processor.js?ts=${Date.now().toString()}`
			);

			workletNode = new AudioWorkletNode(audioContext, "event-processor");

			const streamNode = audioContext.createMediaStreamSource(stream);

			workletNode.port.onmessage = ({ data }) => {
				try {
					callback(validate(audioStreamValidator, data));
				} catch {
					// TODO: maybe spit out errors here. Perhaps report back to home-base
					//   somehow
				}
			};

			streamNode.connect(workletNode);
			workletNode.connect(audioContext.destination);
		}); // TODO: there maybe errors here. Fix them.

		return () => {
			workletNode?.disconnect();
			audioContext?.close();
		};
	}, [stream]);

	return (eventListener: AudioStreamListener) => {
		listeners.current.push(eventListener);
		return () => {
			listeners.current = listeners.current.filter(
				(listener) => eventListener !== listener
			);
		};
	};
}
