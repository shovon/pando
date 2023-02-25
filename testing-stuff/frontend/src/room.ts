import { generateKeys } from "@sparkscience/wskeyid-browser/src/utils";
import { Session } from "./session";
import PubSub, {
	toAsyncIterable,
} from "@sparkscience/wskeyid-browser/src/pub-sub";
import { ROOM_WEBSOCKET_SERVER_ORIGIN } from "./constants";
import {
	InferType,
	object,
	string,
	any,
	ValidationResult,
	ValidationSuccess,
} from "./validator";
import { kvToReadOnlyMap } from "./custom-validators";
import { createSubject, Subscribable } from "./events";
import * as subUtils from "./pub-sub-utils";
import { start as _ } from "./pipe";

const participantSchema = object({
	name: string(),
	// hasVideo: boolean(),
	// hasAudio: boolean(),
});

const participantsListSchema = kvToReadOnlyMap(string(), participantSchema);

const roomStateSchema = object({ participants: participantsListSchema });

type Subject<T> = ReturnType<typeof createSubject<T>>;

const messageFromParticipantSchema = object({ from: string(), data: any() });

export class Participant {
	private unsubscribeFromMessages: () => void;

	constructor(
		messageStream: Subscribable<any>,
		private _sendMessage: (message: string) => void
	) {
		this.unsubscribeFromMessages = messageStream.subscribe((message) => {
			// Handle messages here
		});
	}

	sendMessage(message: string) {
		this._sendMessage(message);
	}

	dispose() {}
}

// TODO: Room#session could potentially be set to null. Refactor this code

export class Room {
	private session: Session | null = null;
	private cancel: boolean = false;
	private _participants: ReadOnlyMap<string, any> = new Map();
	private _roomStateChangeEvents: Subject<void> = createSubject();

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
								this.handleRoomState(data);
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

	private handleRoomState(data: any) {
		const roomState = roomStateSchema.validate(data);
		if (roomState.isValid) {
			this._participants = roomState.value.participants;
			this._roomStateChangeEvents.emit();
		} else {
			console.error("Got something invalid from the server!");
		}
	}

	private refreshParticipants() {
		if (!this.session) {
			console.error(
				"Fatal error! The sesssion object is not defined, for some reason!"
			);
			return;
		}
		const newParticipants = new Map<string, Participant>();
		for (const [id, participant] of this._participants) {
			const cool = _(this.session.messageEvents)
				._((m) =>
					subUtils.map(m, (message) => {
						try {
							const { type, data } = JSON.parse(message.data);
							return { type, data };
						} catch (e) {
							return { type: string };
						}
					})
				)
				._((m) =>
					subUtils.filter(
						m,
						(message) => message.type === "MESSAGE_FROM_PARTICIPANT"
					)
				)
				._((m) =>
					subUtils.map(m, (message) =>
						messageFromParticipantSchema.validate(message.data)
					)
				)
				._((m) =>
					subUtils.filter<
						ValidationResult<InferType<typeof messageFromParticipantSchema>>,
						ValidationSuccess<InferType<typeof messageFromParticipantSchema>>
					>(m, (message) => message.isValid)
				)
				._((m) => subUtils.map(m, (m) => m.value))
				._((m) => subUtils.filter(m, (m) => m.from === id));
			newParticipants.set(
				id,
				new Participant(cool, (message) => {
					this.session!.send(
						JSON.stringify({
							type: "SEND_MESSAGE",
							data: message,
						})
					);
				})
			);
		}
	}

	/**
	 * Gets an event that fires when the room state changes.
	 */
	get roomStateChangeEvents(): Subscribable<void> {
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
