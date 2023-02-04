import { generateKeys } from "@sparkscience/wskeyid-browser/src/utils";
import { Session } from "./session";
import {
	Sub,
	toAsyncIterable,
} from "@sparkscience/wskeyid-browser/src/pub-sub";
import { InferType } from "./validator";
import { RemoteParticipant } from "./components/remote-participant";
import { Disposable, dispose } from "./disposable";

class Participant {
	private _videoStream: MediaStream | null = null;
	private _audioStream: MediaStream | null = null;
	private _screenshareStream: MediaStream | null = null;
	private _remoteParticipant: RemoteParticipant;
	private unsubscribeFromMessages: () => void;

	constructor(
		remoteParticipant: RemoteParticipant,
		messageEvents: Sub<MessageEvent<any>>
	) {
		this._remoteParticipant = remoteParticipant;

		this.unsubscribeFromMessages = messageEvents.addEventListener(
			(message) => {}
		);
	}

	sendMessage() {}

	setRemoteParticipant(p: RemoteParticipant) {
		this._remoteParticipant = p;
	}

	dispose() {
		this.unsubscribeFromMessages();
	}

	get videoStream(): MediaStream | null {
		return this._videoStream;
	}

	get audioStream(): MediaStream | null {
		return this._audioStream;
	}

	get screenshareStream(): MediaStream | null {
		return this._screenshareStream;
	}
}

class ParticipantsManager {
	private _participants: Map<string, Participant> = new Map();

	setRoomState() {}
}

export class Room {
	private _primaryVideo: MediaStream | null = null;
	private _screenshareVideo: MediaStream | null = null;
	private _audio: MediaStream | null = null;
	private _participants: ReadOnlyMap<
		string,
		{
			participant: RemoteParticipant;
			media: {
				audio: MediaStream;
				video: MediaStream;
				screenshare: MediaStream;
			};
		}
	> = new Map();

	constructor(private _roomId: string) {}

	connect() {
		Promise.resolve()
			.then(async function () {
				const keys = await generateKeys();
				const session = new Session("ws://localhost:8080/room/some_room", keys);

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

				for await (const { data: buffer } of toAsyncIterable(
					session.messageEvents
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

	dispose() {}

	get video(): MediaStream | null {
		return this._primaryVideo;
	}

	set video(value: MediaStream | null) {
		this._primaryVideo = value;
	}

	get audio(): MediaStream | null {
		return this._audio;
	}

	set audio(value: MediaStream | null) {
		this._audio = value;
	}

	get screenshareVideo(): MediaStream | null {
		return this._screenshareVideo;
	}

	set screenshareVideo(value: MediaStream | null) {
		this._screenshareVideo = value;
	}

	get roomId() {
		return this._roomId;
	}
}
