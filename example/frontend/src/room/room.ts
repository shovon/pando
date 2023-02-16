import { generateKeys } from "@sparkscience/wskeyid-browser/src/utils";
import { Session } from "../session";
import { toAsyncIterable } from "@sparkscience/wskeyid-browser/src/pub-sub";
import { InferType } from "../validator";
import { Disposable, dispose } from "../disposable";
import { Participant } from "./participant";

class ParticipantsManager {
	private _participants: Map<string, Participant> = new Map();

	setRoomState() {}
}

export class Room {
	private _primaryVideo: MediaStream | null = null;
	private _screenshareVideo: MediaStream | null = null;
	private _audio: MediaStream | null = null;
	private participants: ParticipantsManager = new ParticipantsManager();
	private session: Session | null = null;

	constructor(private _roomId: string) {
		this.connect();
	}

	private connect() {
		Promise.resolve()
			.then(async () => {
				const keys = await generateKeys();
				this.session = new Session("ws://localhost:8080/room/some_room", keys);

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
	 * Gets our current video
	 */
	get video(): MediaStream | null {
		return this._primaryVideo;
	}

	/**
	 * Sets our current video
	 */
	set video(value: MediaStream | null) {
		this._primaryVideo = value;
	}

	/**
	 * Gets our current audio
	 */
	get audio(): MediaStream | null {
		return this._audio;
	}

	/**
	 * Sets our current audio
	 */
	set audio(value: MediaStream | null) {
		this._audio = value;
	}

	/**
	 * Gets our current screenshare
	 */
	get screenshareVideo(): MediaStream | null {
		return this._screenshareVideo;
	}

	/**
	 * Sets our current screenshare
	 */
	set screenshareVideo(value: MediaStream | null) {
		this._screenshareVideo = value;
	}

	/**
	 * Gets the current room's ID.
	 *
	 * This class will never be used for
	 */
	get roomId() {
		return this._roomId;
	}
}
