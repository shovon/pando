import { useState } from "react";
import { useParams } from "react-router-dom";
import { WithId } from "./WithId";

function NamePicker({
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

/**
 * The room view, when a room ID is provided in the URL
 * @returns JSX.Element
 */
export function RoomView() {
	let { roomId } = useParams();

	const [name, setName] = useState<string | null>(null);

	if (!roomId) {
		return (
			<div>
				Unable to get room ID from the URL. This is a mistake done by the
				software developers who wrote this code, and (likely) nothing to do with
				you. Please do report this to the developers, so they can fix it.
			</div>
		);
	}

	if (!name) {
		return <NamePicker onNamePicked={(name) => setName(name)} />;
	}

	return <WithId id={roomId} initialNameValue={name} />;
}
