import { useEffect } from "react";
import { useNavigate } from "react-router-dom";

function randomString(length: number) {
	const characters =
		"ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";

	return Array.from({ length })
		.map(() => characters[Math.floor(Math.random() * characters.length)])
		.join("");
}

export function GetRoom() {
	const navigate = useNavigate();

	useEffect(() => {
		setTimeout(() => {
			navigate(`/room/${randomString(25)}`);
		}, 0);
	});

	// TODO: in the event that we are likely to take longer than 16 seconds,
	//   perhaps show a loader at the 16 seconds mark

	return <></>;
}
