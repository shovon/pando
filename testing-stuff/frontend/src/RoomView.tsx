import { useEffect, useReducer, useRef, useState } from "react";
import { useParams } from "react-router-dom";
import { Room } from "./room";

function RoomWithId({ id }: { id: string }) {
	const roomRef = useRef<Room | null>(null);
	const [update] = useReducer(() => ({}), {});

	useEffect(() => {
		console.log("Connecting");
		roomRef.current = new Room(id, "test");

		return () => {
			roomRef.current?.dispose();
		};
	}, []);

	if (!roomRef.current) {
		return <div>Creating local room</div>;
	}

	return <div>Room</div>;
}

export function RoomView() {
	let { id } = useParams();

	if (!id) {
		return <div>Invalid state!</div>;
	}

	return <RoomWithId id={id} />;
}
