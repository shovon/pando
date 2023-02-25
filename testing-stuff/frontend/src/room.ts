import { generateKeys } from "@sparkscience/wskeyid-browser/src/utils";
import { Session } from "./session";
import PubSub, {
	toAsyncIterable,
} from "@sparkscience/wskeyid-browser/src/pub-sub";
import { ROOM_WEBSOCKET_SERVER_ORIGIN } from "./constants";
import { InferType, object, string } from "./validator";
import { kvToReadOnlyMap } from "./custom-validators";
import { createSubject, Subscribable } from "./events";
import { Participant } from "./participant";

const participantSchema = object({
	name: string(),
	// hasVideo: boolean(),
	// hasAudio: boolean(),
});

const participantsListSchema = kvToReadOnlyMap(string(), participantSchema);

const roomStateSchema = object({ participants: participantsListSchema });

type Subject<T> = ReturnType<typeof createSubject<T>>;

export class Room {
	private session: Session | null = null;
	private cancel: boolean = false;
	private _participants: ReadOnlyMap<string, any> = new Map();
	private _roomStateChangeEvents: Subject<InferType<typeof roomStateSchema>> =
		createSubject();

	constructor(private _roomId: string, private _name: string) {
		this.connect();
	}

	private connect() {
		Promise.resolve()
			.then(async () => {
				const keys = await generateKeys();

				if (this.cancel) {
					return;
				}
				this.session = new Session(
					`${ROOM_WEBSOCKET_SERVER_ORIGIN}/room/${this._roomId}}`,
					keys
				);

				let session = this.session;

				session.sessionStatusChangeEvents.addEventListener((status) => {
					if (status.type === "CONNECTING") {
						console.log(
							"status type is %s and sub status is %s",
							status.type,
							status.status
						);
					} else if (status.type === "CONNECTED") {
						session.send(
							JSON.stringify({ type: "SET_NAME", data: this._name })
						);
						console.log(
							"Connected to room %s with name %s",
							this._roomId,
							this._name
						);
					} else {
						console.log("Status type is %s", status.type);
					}
				});

				for await (const { data: buffer } of toAsyncIterable(
					session.messageEvents
				)) {
					try {
						const { type, data } = JSON.parse(buffer);
						switch (type) {
							case "ROOM_STATE":
								const roomState = roomStateSchema.validate(data);
								if (roomState.isValid) {
									this._participants = roomState.value.participants;
									this._roomStateChangeEvents.emit(roomState.value);
								} else {
									console.error("Got something invalid from the server!");
								}
						}
					} catch (e) {
						console.error(e);
					}
				}
			})
			.catch(console.error);
	}

	dispose() {
		this.cancel = true;
		this.session?.endSession();
	}

	/**
	 * Gets an event that fires when the room state changes.
	 */
	get roomStateChangeEvents(): Subscribable<InferType<typeof roomStateSchema>> {
		return this._roomStateChangeEvents;
	}

	/**
	 * Gets the current room's ID.
	 *
	 * This class will never be used for
	 */
	get roomId() {
		return this._roomId;
	}

	/**
	 * Gets our current name in the room.
	 */
	get name() {
		return this._name;
	}

	/**
	 * Gets the current list of participants in the room.
	 */
	get participants(): ReadOnlyMap<string, InferType<typeof participantSchema>> {
		return this._participants;
	}
}
