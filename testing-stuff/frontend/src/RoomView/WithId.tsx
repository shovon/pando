import { useEffect, useReducer, useRef, useState } from "react";
import { Room } from "../room";
import * as it from "../iterables";

// This is the actual room view where all the jazz will happen.
//
// Go wild
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

async function createRoom(id: string, initialNameValue: string) {
	const room = new Room(id, initialNameValue);
	return room;
}

/**
 * This is the actual room view containing the ID for the room
 * @param param0 The props for the component
 * @returns JSX.Element
 */
export function WithId({
	id,
	initialNameValue,
}: {
	id: string;
	initialNameValue: string;
}) {
	const [room, setRoom] = useState<Room | null>(null);
	const [error, setError] = useState<object | null>(null);

	useEffect(() => {
		createRoom(id, initialNameValue)
			.then(setRoom)
			.catch((e) => {
				setError(e);
			});

		return () => {
			room?.dispose();
		};
	}, []);

	if (error) {
		return (
			<div>
				An unhandled error occurred. This is no fault of your own. This is the
				mistake of the developers. If you are seeing this, please report this to
				the developers so this can be fixed.
			</div>
		);
	}

	if (!room) {
		return <div>Creating local room</div>;
	}

	return <WithRoom room={room} />;
}
