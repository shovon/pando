export class RoomControl {
	private _mediaStreamTracks = new Map<string, Map<string, MediaStreamTrack>>();
	private pool = new WeakMap<MediaStreamTrack, MediaStream>();

	setStream(subject: string, stream: MediaStream) {
		for (const track of stream.getTracks()) {
			this.setTrack(subject, track);
			this.pool.set(track, stream);
		}
	}

	setTrack(subject: string, track: MediaStreamTrack) {
		let tracks = this._mediaStreamTracks.get(track.kind);
		if (!tracks) {
			tracks = new Map();
			this._mediaStreamTracks.set(track.kind, tracks);
		}

		tracks.set(subject, track);
	}

	getStream(subject: string, kind: string): MediaStream | null {
		const tracks = this._mediaStreamTracks.get(kind);
		if (!tracks) {
			return null;
		}

		const track = tracks.get(subject);
		if (!track) {
			return null;
		}

		if (track.kind !== kind) {
			throw new Error(
				`Fatal error! Expected a track of kind ${kind} but got ${track.kind}`
			);
		}

		let stream = this.pool.get(track);
		if (!stream) {
			stream = new MediaStream([track]);
			this.pool.set(track, stream);
		}

		return stream;
	}

	get mediaStreamTracks(): ReadOnlyMap<
		string,
		ReadOnlyMap<string, MediaStreamTrack>
	> {
		return this._mediaStreamTracks;
	}

	static getMediaDevicesList() {
		return navigator.mediaDevices.enumerateDevices();
	}

	static getVideo(deviceId?: string) {
		return navigator.mediaDevices.getUserMedia({
			// TODO: check to see what happens if a bad device ID is inserted
			video: deviceId ? { deviceId } : true,
		});
	}

	static getAudio(deviceId?: string) {
		return navigator.mediaDevices.getUserMedia({
			// TODO: check to see what happens if a bad device ID is inserted
			audio: deviceId ? { deviceId } : true,
		});
	}
}
