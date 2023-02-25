import { useEffect, useReducer, useRef, useState } from "react";
import { useParams } from "react-router-dom";
import { Room } from "./room";
import * as it from "./iterables";

function RoomWithId({ id }: { id: string }) {
	const roomRef = useRef<Room | null>(null);
	const [, update] = useReducer(() => ({}), {});

	useEffect(() => {
		console.log("Connecting");
		roomRef.current = new Room(id, "test");
		update();

		return () => {
			roomRef.current?.dispose();
		};
	}, []);

	if (!roomRef.current) {
		return <div>Creating local room</div>;
	}

	return <WithRoom room={roomRef.current} />;
}

function WithRoom({ room }: { room: Room }) {
	const [participants, setParticipantsSet] = useState(room.participants);

	useEffect(() => {
		const unsubscribe = room.roomStateChangeEvents.subscribe(() => {
			setParticipantsSet(room.participants);
		});

		return () => {
			unsubscribe();
		};
	});
	return (
		<div>{it.map(participants, ([k, participant]) => participant.name)}</div>
	);
}

export function RoomView() {
	let { id } = useParams();

	if (!id) {
		return <div>Invalid state!</div>;
	}

	return <RoomWithId id={id} />;
}
