import { useEffect, useState } from "react";
import { AudioMeter } from "./AudioMeter/AudioMeter";
import { StreamPlayer } from "./StreamPlayer";
import { RoomControl } from "./room-control";

type CallMediaSelectorProps = {
	mediaSet: (props: {
		video: MediaStream | null;
		audio: MediaStream | null;
	}) => void;
};

/**
 * Represents the room in a call
 * @returns JSX component that represents the DOM layout of the room
 */
export function CallMediaSelector({ mediaSet }: CallMediaSelectorProps) {
	// TODO: find a way to prevent resetting the video every time this modal is
	//   invoked

	// This is the video that we are previewing
	const [video, setVideo] = useState<MediaStream | null>(null);

	// This is the audio that we are previewing
	const [audio, setAudio] = useState<MediaStream | null>(null);

	// The list of media devices (both audio and video) that we want to select
	// from
	const [mediaDevices, setMediaDevices] = useState<MediaDeviceInfo[] | null>(
		null
	);

	function getVideoDevicesList(): MediaDeviceInfo[] | null {
		if (!mediaDevices) {
			return null;
		}

		return mediaDevices.filter((info) => info.kind === "videoinput");
	}

	function getAudioDevices(): MediaDeviceInfo[] | null {
		if (!mediaDevices) {
			return null;
		}

		return mediaDevices.filter((info) => info.kind === "audioinput");
	}

	function getCurrentlySelectedVideoId() {
		return video?.getTracks()[0].getSettings().deviceId ?? "";
	}

	function getCurrentlySelectedAudioId() {
		return audio?.getTracks()[0].getSettings().deviceId ?? "";
	}

	useEffect(() => {
		Promise.resolve().then(async () => {
			await Promise.allSettled([
				Promise.resolve().then(async () => {
					setVideo(await RoomControl.getVideo());
				}),
				Promise.resolve().then(async () => {
					setAudio(await RoomControl.getAudio());
				}),
				Promise.resolve().then(async () => {
					setMediaDevices(await RoomControl.getMediaDevicesList());
				}),
			]);
		});
	}, []);

	const videoDevices = getVideoDevicesList();
	const audioDevices = getAudioDevices();

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

			{audioDevices ? (
				<div>
					<label>Audio devices</label>
					<select
						onChange={(event) => {
							Promise.resolve(event.target.value).then(async (deviceId) => {
								setAudio(await RoomControl.getAudio(deviceId));
							}); // TODO: catch the error here, and do something
						}}
						value={getCurrentlySelectedAudioId()}
					>
						{audioDevices.map((device) => (
							<option value={device.deviceId} key={device.deviceId}>
								{device.label}
							</option>
						))}
					</select>
				</div>
			) : null}

			{audio ? <AudioMeter audioStream={audio} /> : null}

			<button
				onClick={() => {
					mediaSet({ audio, video });
				}}
			>
				Accept
			</button>
		</div>
	);
}
