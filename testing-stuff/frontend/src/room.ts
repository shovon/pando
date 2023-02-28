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
import { createSubject, Subject, Subscribable } from "./events";
import * as subUtils from "./pub-sub-utils";
import { start as _ } from "./pipe";

const participantSchema = object({
	name: string(),
	// hasVideo: boolean(),
	// hasAudio: boolean(),
});

// The participants list
const participantsListSchema = kvToReadOnlyMap(string(), participantSchema);

// The room state that the server will be giving us.
//
// Note: we're not going to be storing the room state that the server gives us
// as-is. We're going to be holding our own local instance of each participant
// so that we can send and receive messages to/from them.
const roomStateSchema = object({ participants: participantsListSchema });

// The message that a remote participant will be relaying to us through the
// server
const messageFromParticipantSchema = object({ from: string(), data: any() });

/**
 * A participant in a room
 */
export class Participant {
	// TODO: consider including a flag to indicate that the participant is
	//   unresponsive

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

// TODO: Room#session could potentially be set to null. Refactor this code.
//   Perhaps use dependency injection to inject the session into the room

/**
 *
 */
export class Room {
	private session: Session | null = null;
	private cancel: boolean = false;

	// TODO: when the room state changes, not always will a missing participant
	//   imply that the participant has left the room.
	//
	//   So we are going to have to handle the edge case that when the connection
	//   with the room server has been severed, then wait an n number of time
	//   before killing the participant.
	//
	//   Perhaps maintain a flag somewhere.
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

	/**
	 * Handles the event that we got a ROOM_STATE message from the server
	 * @param data The data to process
	 */
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
