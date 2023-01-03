import AuthenticatedConnection, {
	ConnectionStatus,
} from "@sparkscience/wskeyid-browser/src/authenticated-connection";
import PubSub, {
	getNext,
	Sub,
} from "@sparkscience/wskeyid-browser/src/pub-sub";

const backoffMSIncrement = 120;
const maxBackoffExponent = 9;

type SessionStatus =
	| ConnectionStatus
	| {
			type: "SLEEPING";
			sleepUntil: number;
			restartReason: { type: "CLOSED" } | { type: "ERROR"; data: any };
	  };

export class Session {
	private connection: AuthenticatedConnection | null = null;
	private readonly _messageEvents: PubSub<MessageEvent> = new PubSub();
	private _sessionStatus: Readonly<SessionStatus> = { type: "PENDING" };
	private readonly _sessionStatusChangeEvents: PubSub<SessionStatus> =
		new PubSub();
	private messagesBuffer: string[] = [];

	private backoffExponent = 0;
	private sleeping = false;

	private async restart(error: { data: any } | null) {
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
				this.restart(null);
			}
		});
		this.connection.messageEvents.addEventListener((message) => {
			this._messageEvents.emit(message);
		});
	}

	send(message: string) {
		if (this.connection) {
			this.connection.send(message);
		} else {
			this.messagesBuffer.push(message);
		}
	}

	getNextMessage() {
		return getNext(this._messageEvents);
	}

	get messageEvents(): Sub<MessageEvent> {
		return this._messageEvents;
	}

	get sessionStatus(): SessionStatus {
		return this._sessionStatus;
	}

	get sessionStatusChangeEvents(): Sub<SessionStatus> {
		return this._sessionStatusChangeEvents;
	}
}
