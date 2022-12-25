import { useEffect, useState } from "react";
import { StreamPlayer } from "./StreamPlayer";

export function Room() {
	const [video, setVideo] = useState<MediaStream | null>(null);

	function getVideo() {
		return navigator.mediaDevices.getUserMedia({ video: true });
	}

	function getAudio() {
		return navigator.mediaDevices.getUserMedia({ audio: true });
	}

	useEffect(() => {
		Promise.resolve().then(async () => {
			setVideo(await getVideo());
		});
	}, []);

	return (
		<div>
			{video ? (
				<StreamPlayer style={{ transform: "scaleX(-1)" }} stream={video} />
			) : null}
		</div>
	);
}
