import { useEffect, useReducer, useRef, useState } from "react";
import { Room } from "../room";
import * as it from "../iterables";

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
		<ul>
			{it.map(participants, ([k, participant]) => (
				<li key={k}>{participant.name}</li>
			))}
		</ul>
	);
}

export function WithId({
	id,
	initialNameValue,
}: {
	id: string;
	initialNameValue: string;
}) {
	const roomRef = useRef<Room | null>(null);
	const [, update] = useReducer(() => ({}), {});

	useEffect(() => {
		console.log("Connecting");
		roomRef.current = new Room(id, initialNameValue);
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
