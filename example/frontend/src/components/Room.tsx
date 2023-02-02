import { useEffect, useReducer, useRef, useState } from "react";
import { AudioMeter } from "./AudioMeter/AudioMeter";
import { StreamPlayer } from "./StreamPlayer";
import { RoomControl } from "./room-control";
import { CallMediaSelector } from "./CallMediaSelector";
import { Navigate, useParams } from "react-router-dom";
import { Room } from "../room";

function InRoom({ roomId }: { roomId: string }) {
	const roomRef = useRef<Room>(new Room(roomId));

	function connect({
		video,
		audio,
	}: {
		video: MediaStream;
		audio: MediaStream;
	}) {}

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
	const [isSelectingStreams, setIsSelectingStreams] = useState(true);

	if (isSelectingStreams) {
		return (
			<CallMediaSelector
				onMediaSet={({ audio, video }) => {
					setIsSelectingStreams(false);
				}}
			/>
		);
	}

	return (
		<div>
			{video ? (
				<StreamPlayer style={{ transform: "scaleX(-1)" }} stream={video} />
			) : null}
		</div>
	);
}

/**
 * Represents the room in a call
 * @returns JSX component that represents the DOM layout of the room
 */
export function RoomView() {
	let { id } = useParams();

	if (!id) {
		return <Navigate to="/" />;
	}

	return <InRoom roomId={id} />;
}
