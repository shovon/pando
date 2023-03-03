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
	Validator,
	exact,
} from "./validator";
import { json, kvToReadOnlyMap } from "./custom-validators";
import { createSubject, Subject, Subscribable } from "./events";
import * as subUtils from "./pub-sub-utils";
import { start as _ } from "./pipe";
import { Sub } from "@sparkscience/wskeyid-browser/src/pub-sub";
import { Participant } from "./participant";

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
 * The main room instance for the call
 */
export class Room {
	private _participants: ReadOnlyMap<string, Participant> = new Map();

	private _roomStateChangeEvents: Subject<void> = createSubject();

	// This is just os odd. Oh well, it'll do for now
	private _failedEvents: Subject<void> = createSubject();

	/**
	 * Initializes a new Room instance
	 * @param session The session instance
	 * @param _name The participant's name to connect to the room with
	 */
	constructor(private session: Session, private _name: string) {
		this.connect();
	}

	private connect() {
		Promise.resolve()
			.then(async () => {
				const session = this.session;

				session.sessionStatusChangeEvents.addEventListener((status) => {
					if (status.type === "CONNECTING") {
						console.log(
							"status type is %s and sub status is %s",
							status.type,
							status.status
						);
					} else if (status.type === "CONNECTED") {
						console.log(
							"Connected. Now sending name and awaiting name and room status"
						);
						session.send(
							JSON.stringify({ type: "SET_NAME", data: this._name })
						);
					} else {
						console.log("Status type is %s", status.type);
					}
				});

				// And here we just go loop-de-loop and listen for new messages coming
				// in from the server
				for await (const { data: buffer } of toAsyncIterable(
					session.messageEvents
				)) {
					try {
						const { type, data } = JSON.parse(buffer);

						// Might get ugly. We might need to move this code elsewhere
						// eventually
						switch (type) {
							case "ROOM_STATE":
								this.handleRoomState(data);
						}
					} catch (e) {
						console.error(e);
					}
				}
			})
			.catch((e) => {
				console.error(e);
				this._failedEvents.emit();
			});
	}

	/**
	 * Gets the participants in the room
	 */
	dispose() {
		this._participants.forEach((participant) => {
			participant.dispose();
		});

		// TOOD: maybe it's a bad idea to do this here?
		this.session.endSession();
	}

	// Handles the changes to the rooms' state
	private handleRoomState(data: any) {
		const roomState = roomStateSchema.validate(data);
		if (roomState.isValid) {
			console.log(roomState.value);
			roomState.value.participants;
			this._roomStateChangeEvents.emit();
		} else {
			// TODO: do something more serious so that the user is notified that
			//   something really bad happened, to no fault of their own
			console.error("Got something invalid from the server!");
		}
	}

	// Filters out all messages that don't match the supplied schema
	private getMessageOfSchema<T>(validator: Validator<T>): Sub<T> {
		return _(this.session.messageEvents)
			._((m) =>
				subUtils.map(m, (message) => {
					return validator.validate(message.data);
				})
			)
			._((m) =>
				subUtils.filter<ValidationResult<T>, ValidationSuccess<T>>(
					m,
					(message) => message.isValid
				)
			)
			._((m) => subUtils.map(m, (message) => message.value)).value;
	}

	// Gets the message stream that is destined for a single participant
	private getMessageEventsFromParticipant(id: string) {
		return _(
			this.getMessageOfSchema(
				json(
					object({
						type: exact("MESSAGE_FROM_PARTICIPANT"),
						data: object({
							from: string(),
							data: any(),
						}),
					})
				)
			)
		)
			._((m) => subUtils.map(m, (message) => message.data))
			._((m) => subUtils.filter(m, (m) => m.from === id))
			._((m) => subUtils.map(m, (m) => m.data)).value;
	}

	private refreshParticipants(
		participants: InferType<typeof participantsListSchema>
	) {
		const newParticipants = new Map<string, Participant>();
		for (const [id, participant] of participants) {
			newParticipants.set(
				id,
				new Participant(
					this.getMessageEventsFromParticipant(id),
					(message: any) => {
						this.session.send(
							JSON.stringify({
								type: "MESSAGE_TO_PARTICIPANT",
								data: {
									to: id,
									data: message,
								},
							})
						);
					}
				)
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
	 * Gets our current name in the room.
	 */
	get name() {
		return this._name;
	}

	/**
	 * Sets the name of the participant that is in the room
	 */
	set name(name: string) {
		// throw new Error("Not yet implemented");
		this._name = name;
		this.session.send(JSON.stringify({ type: "SET_NAME", data: name }));
	}

	/**
	 * Gets the current list of participants in the room.
	 */
	get participants(): ReadOnlyMap<string, InferType<typeof participantSchema>> {
		return this._participants;
	}
}
