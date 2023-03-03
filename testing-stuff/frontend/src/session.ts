import AuthenticatedConnection, {
	ConnectionStatus,
} from "@sparkscience/wskeyid-browser/src/authenticated-connection";
import PubSub, {
	getNext,
	Sub,
} from "@sparkscience/wskeyid-browser/src/pub-sub";
import { createEmitOnce, OnceEmitter, SubjectOnce } from "./emit-once";

const backoffMSIncrement = 120;
const maxBackoffExponent = 9;

type SessionStatus =
	| ConnectionStatus
	| {
			type: "SLEEPING";
			sleepUntil: number;
			restartReason: { type: "CLOSED" } | { type: "ERROR"; data: any };
	  };

// An arbitrary amount of time in milliseconds to wait before considering that
// the connection is steady after a restart
export const restartTimeout = 30000;

/**
 * A wrapper object around AuthenticatedConnection. Not only this class performs
 * the wskey-id handshake, but unlike AuthenticatedConnection, if the underlying
 * WebSocket connection is severed, a reconnection attempt will occur.
 */
export class Session {
	private connection: AuthenticatedConnection | null = null;
	private readonly _messageEvents: PubSub<MessageEvent> = new PubSub();
	private _sessionStatus: Readonly<SessionStatus> = { type: "PENDING" };
	private readonly _sessionStatusChangeEvents: PubSub<SessionStatus> =
		new PubSub();
	private messagesBuffer: string[] = [];
	private sessionEndedOnceEmit: SubjectOnce<void> = createEmitOnce();

	private backoffExponent = 0;
	private sleeping = false;

	private _recentlyRestartedStartTimeAndTimeout: {
		startTime: Date;
		timeout: number;
	} | null = null;

	private async restart(error: { data: any } | null) {
		if (this._recentlyRestartedStartTimeAndTimeout) {
			clearTimeout(this._recentlyRestartedStartTimeAndTimeout.timeout);
			this._recentlyRestartedStartTimeAndTimeout = null;
		}

		if (this.sleeping) {
			return;
		}

		const backoffExponent = this.backoffExponent;

		this.backoffExponent =
			this.backoffExponent >= maxBackoffExponent
				? maxBackoffExponent
				: this.backoffExponent + 1;

		const timeoutTime = backoffMSIncrement * 2 ** backoffExponent;
		const sleepUntil = Date.now() + timeoutTime;

		if (error) {
			const errorStatus: Readonly<ConnectionStatus> =
				Object.freeze<ConnectionStatus>({
					type: "CLOSED",
					reason: { type: "CONNECTION_ERROR", data: error },
				});

			this.setSessionStatus(errorStatus);
		}

		this.setSessionStatus({
			type: "SLEEPING",
			sleepUntil: sleepUntil,
			restartReason: error
				? {
						type: "ERROR",
						data: error.data,
				  }
				: {
						type: "CLOSED",
				  },
		});

		await new Promise((resolve) => {
			this.sleeping = true;
			setTimeout(() => {
				this.sleeping = false;
				this.connect()
					.catch((e) => {
						this.restart({ data: e });
					})
					.then(resolve);
			}, timeoutTime);
		});
	}

	/**
	 * Initalizes a new Session object, and also ocnnects to the WebSocket server
	 * @param url The URL to connect to the WebSocket server
	 * @param key The keys to use for the wskey-id handshake
	 */
	constructor(
		private readonly url: string,
		private readonly key: CryptoKeyPair
	) {
		this.connect().catch((e) => {
			this.restart(e);
		});
	}

	private setSessionStatus(status: SessionStatus) {
		this._sessionStatus = status;
		setTimeout(() => {
			this._sessionStatusChangeEvents.emit(status);
		});
	}

	private async connect() {
		this.connection = await AuthenticatedConnection.connect(this.url, this.key);

		this.connection.sessionStatusChangeEvents.addEventListener((status) => {
			if (status.type === "CONNECTED") {
				this.backoffExponent = 0;

				if (this.connection) {
					for (const message of this.messagesBuffer) {
						this.connection.send(message);
					}
					this.messagesBuffer = [];
				}
			}
			this.setSessionStatus(status);
			if (status.type === "CLOSED") {
				if (!this.isClosed) {
					this.restart(null);
				}
			}
		});
		this.connection.messageEvents.addEventListener((message) => {
			this._messageEvents.emit(message);
		});
	}

	/**
	 * Closes the connection to the WebSocket server, and stops any reconnection
	 */
	endSession() {
		// TODO: send a message to the server to tell it to close the connection
		this.connection?.close();
		this.sessionEndedOnceEmit.emit();
	}

	/**
	 * Sends a message to the WebSocket server
	 * @param message The message to send to the WebSocket server. Does not need
	 *   to be a JSON string.
	 */
	send(message: string) {
		if (this.connection) {
			this.connection.send(message);
		} else {
			this.messagesBuffer.push(message);
		}
	}

	/**
	 * From the message events, gets the next message event from the WebSocket
	 * @returns The next message event from the WebSocket server
	 */
	getNextMessage() {
		return getNext(this._messageEvents);
	}

	/**
	 * Gets the message event emitter, containing an event stream of messages from
	 * the WebSocket server
	 */
	get messageEvents(): Sub<MessageEvent> {
		return this._messageEvents;
	}

	/**
	 * Gets the current status of the connection to the WebSocket server
	 */
	get sessionStatus(): SessionStatus {
		return this._sessionStatus;
	}

	/**
	 * Gets the session status change event emitter, containing an event stream of
	 * changes to the connection status
	 */
	get sessionStatusChangeEvents(): Sub<SessionStatus> {
		return this._sessionStatusChangeEvents;
	}

	/**
	 * Gets whether the session has ended. If this is true, the session will not
	 * attempt to reconnect to the WebSocket server
	 */
	get isClosed() {
		return this.sessionEndedOnceEmit.value.isEmitted;
	}

	/**
	 * Gets the single emitting event emitter for the event when the session ends
	 */
	get sessionEnded(): OnceEmitter<void> {
		return this.sessionEndedOnceEmit;
	}
}
