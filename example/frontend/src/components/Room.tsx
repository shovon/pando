import { useEffect, useReducer, useRef, useState } from "react";
import { AudioMeter } from "./AudioMeter/AudioMeter";
import { StreamPlayer } from "./StreamPlayer";
import { RoomControl } from "./room-control";
import { CallMediaSelector } from "./CallMediaSelector";
import { useParams } from "react-router-dom";
import { Room } from "../room";

/**
 * Represents the room in a call
 * @returns JSX component that represents the DOM layout of the room
 */
export function RoomView() {
	let { id } = useParams();

	const roomRef = useRef<Room>(new Room());

	const [video, setVideo] = useState<MediaStream | null>(null);
	const [audio, setAudio] = useState<MediaStream | null>(null);

	// Reminder of the day: this will change to a more elaborate state management
	// system.
	//
	// Start small, then expand.
	//
	// That said DO NOT just add more state! It's time to start moving things out
	// to maintain synchronicity.
	//
	// That said, it's also important to
	const [isSelectingVideo, setIsSelectingVideo] = useState(true);

	if (isSelectingVideo) {
		return (
			<CallMediaSelector
				mediaSet={({ audio, video }) => {
					setAudio(audio);
					setVideo(video);
					setIsSelectingVideo(false);
				}}
			/>
		);
	}

	return <></>;
}
