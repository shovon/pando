import { RemoteParticipant } from "../components/remote-participant";
import { Sub } from "@sparkscience/wskeyid-browser/src/pub-sub";

/**
 * This represents a remote participant.
 *
 * This class is just a fancy tool to help independently manage incoming
 * streams, without having some central entity managing it all. Too much
 * headache, so very little reward.
 *
 * Perhaps, in the future, the single participant will also be used to handle
 * other things such as incoming DMs from the specified participant
 */
export class Participant {
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
