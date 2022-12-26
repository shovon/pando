import { useEffect, useState } from "react";
import { AudioMeter } from "./AudioMeter/AudioMeter";
import { StreamPlayer } from "./StreamPlayer";
import { RoomControl } from "./room-control";

export function Room() {
	const [video, setVideo] = useState<MediaStream | null>(null);
	const [audio, setAudio] = useState<MediaStream | null>(null);
	const [mediaDevices, setMediaDevices] = useState<MediaDeviceInfo[] | null>(
		null
	);

	function getVideoDevices(): MediaDeviceInfo[] | null {
		if (!mediaDevices) {
			return null;
		}

		return mediaDevices.filter((info) => info.kind === "videoinput");
	}

	function getCurrentlySelectedVideoId() {
		return video?.getTracks()[0].getSettings().deviceId ?? "";
	}

	useEffect(() => {
		Promise.resolve().then(async () => {
			setVideo(await RoomControl.getVideo());
			setAudio(await RoomControl.getAudio());
			setMediaDevices(await RoomControl.getMediaDevicesList());
		});
	}, []);

	const videoDevices = getVideoDevices();

	return (
		<div>
			{videoDevices ? (
				<div>
					<label>Video devices</label>
					<select
						onChange={(event) => {
							Promise.resolve(event.target.value).then(async (deviceId) => {
								setVideo(await RoomControl.getVideo(deviceId));
							}); // TODO: catch the error here, and do something
						}}
						value={getCurrentlySelectedVideoId()}
					>
						<option value="" />
						{videoDevices.map((device) => (
							<option value={device.deviceId} key={device.deviceId}>
								{device.label}
							</option>
						))}
					</select>
				</div>
			) : null}

			{video ? (
				<StreamPlayer style={{ transform: "scaleX(-1)" }} stream={video} />
			) : null}
			{audio ? <AudioMeter audioStream={audio} /> : null}
		</div>
	);
}
