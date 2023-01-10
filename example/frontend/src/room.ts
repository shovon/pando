import { generateKeys } from "@sparkscience/wskeyid-browser/src/utils";
import { Session } from "./session";
import { toAsyncIterable } from "@sparkscience/wskeyid-browser/src/pub-sub";

export class Room {
	private video: MediaStream | null = null;
	private audio: MediaStream | null = null;

	constructor(roomId: string) {}

	connect() {
		Promise.resolve()
			.then(async function () {
				const keys = await generateKeys();
				const session = await new Session(
					"ws://localhost:8080/room/some_room",
					keys
				);

				session.sessionStatusChangeEvents.addEventListener((status) => {
					if (status.type === "CONNECTING") {
						console.log(
							"status type is %s and sub status is %s",
							status.type,
							status.status
						);
					} else {
						console.log("Status type is %s", status.type);
					}
				});

				for await (const event of toAsyncIterable(session.messageEvents)) {
					console.log(event);
				}
			})
			.catch(console.error);
	}
}
