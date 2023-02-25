import { useEffect, useReducer, useRef, useState } from "react";
import { useParams } from "react-router-dom";
import { Room } from "./room";
import * as it from "./iterables";

function RoomWithId({
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

export function NamePicker({
	onNamePicked,
}: {
	onNamePicked: (name: string) => void;
}) {
	const [name, setName] = useState<string>("");

	return (
		<input
			type="text"
			value={name || ""}
			onChange={(e) => setName(e.target.value)}
			onKeyDown={(e) => {
				if (e.key === "Enter") {
					onNamePicked(name);
				}
			}}
		/>
	);
}

export function RoomView() {
	let { id } = useParams();
	const [name, setName] = useState<string | null>(null);

	if (!id) {
		return <div>Invalid state!</div>;
	}

	if (!name) {
		return <NamePicker onNamePicked={(name) => setName(name)} />;
	}

	return <RoomWithId id={id} initialNameValue={name} />;
}
