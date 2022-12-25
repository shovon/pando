import { useEffect } from "react";
import { useAudioStreamMeter } from "../hooks/use-audio-meter";

type AudioMeterProps = {
	audioStream: MediaStream;
};

export function AudioMeter({ audioStream }: AudioMeterProps) {
	const subscribeToAudio = useAudioStreamMeter(audioStream);

	useEffect(() => {
		const unsubscribe = subscribeToAudio((event) => {});
	});

	return <div></div>;
}
