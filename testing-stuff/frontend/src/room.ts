import { generateKeys } from "@sparkscience/wskeyid-browser/src/utils";
import { Session } from "./session";
import { toAsyncIterable } from "@sparkscience/wskeyid-browser/src/pub-sub";
import { ROOM_WEBSOCKET_SERVER_ORIGIN } from "./constants";

export class Room {
	private session: Session | null = null;

	constructor(private _roomId: string, private _name: string) {
		this.connect();
	}

	private connect() {
		Promise.resolve()
			.then(async () => {
				const keys = await generateKeys();
				console.log(ROOM_WEBSOCKET_SERVER_ORIGIN);
				this.session = new Session(
					`${ROOM_WEBSOCKET_SERVER_ORIGIN}/room/some_room`,
					keys
				);

				this.session.sessionStatusChangeEvents.addEventListener((status) => {
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

				this.session.send(
					JSON.stringify({ type: "SET_NAME", data: this._name })
				);

				for await (const { data: buffer } of toAsyncIterable(
					this.session.messageEvents
				)) {
					try {
						const { type, data } = JSON.parse(buffer);
						switch (type) {
							case "ROOM_STATE":
								console.log("Got room state", data);
						}
					} catch (e) {
						console.error(e);
					}
				}
			})
			.catch(console.error);
	}

	dispose() {
		this.session?.endSession();
	}

	/**
	 * Gets the current room's ID.
	 *
	 * This class will never be used for
	 */
	get roomId() {
		return this._roomId;
	}

	get name() {
		return this._name;
	}
}
