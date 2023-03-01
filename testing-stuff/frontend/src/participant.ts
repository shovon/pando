import { Sub } from "@sparkscience/wskeyid-browser/src/pub-sub";
import { start as _ } from "./pipe";
import { filter } from "./pub-sub-utils";

// Right now participant-to-participant messages are messages associated with a
// specific participant.
//
// But eventually, we will also want to receive meta-information regarding a
// particular participant.
//
// Question is, should the room manager handle these meta information or should
// the participant handle it?

/**
 * This represents a remote participant.
 *
 * This class is just a fancy tool to help independently manage incoming
 * streams, without having some central entity managing it all. Too much
 * headache, so very little reward
 */
export class Participant {
	private _videoStream: MediaStream | null = null;
	private _audioStream: MediaStream | null = null;
	private _screenshareStream: MediaStream | null = null;
	private unsubscribeFromMessages: () => void;
	private _receivedMessages: string[] = [];

	constructor(
		messageEvents: Sub<MessageEvent<any>>,
		private _sendMessage: (message: any) => void
	) {
		this.unsubscribeFromMessages = filter(
			messageEvents,
			(m) => !!m
		).addEventListener((message) => {
			// TOOD: handle messages
		});
	}

	sendMessage(message: any) {
		this._sendMessage(message);
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
