import { useCallback, useEffect, useRef, useState } from "react";
import { useAnimationFrame } from "./use-animation-frame";
import { useAudioStreamMeter } from "./use-audio-meter";

type AudioMeterProps = {
	audioStream: MediaStream;
};

export function AudioMeter({ audioStream }: AudioMeterProps) {
	const subscribeToAudio = useAudioStreamMeter(audioStream);
	const amplitudeRef = useRef(0);
	const [amplitude, setAmplitude] = useState(amplitudeRef.current);
	const animationFrameCallback = useCallback(() => {
		setAmplitude(amplitudeRef.current);
	}, [audioStream]);

	useEffect(() => {
		const unsubscribe = subscribeToAudio((event) => {
			const result = event.flatMap((a) =>
				a.flatMap((b) => b.flatMap((c) => c))
			);

			const mean = Math.abs(
				result.reduce((prev, next) => prev + next) / result.length
			);

			amplitudeRef.current = mean;
		});

		return () => {
			unsubscribe();
		};
	}, [audioStream]);

	useAnimationFrame(animationFrameCallback);

	return (
		<div>
			<div
				style={{
					background: "green",
					height: 10,
					width: `${amplitude * 100}%`,
				}}
			></div>
		</div>
	);
}
