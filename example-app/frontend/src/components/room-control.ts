// There is a difference between muting and unavailable media device. Question
// is, should we conflate the two? Like is deleting the video from the call the
// same as simply cutting the outbound video stream?
//
// Let's get somethings straight.
//
// A participant will have a list of media devices.
//
// Possible states:
//
// - media available and streaming but disabled (e.g. a "mute" audio, or "hide"
//   video, or "disabled' media flag on the server)
// - media available but not streaming
//   - how do we handle this? How do we know this is intentional or not?
// 	  - bear in mind, it is very difficult to distinguish between "not streaming
//      because participant disabled on their end", and "failing to establish a
//      read stream")
// - media not available (the participant simply is not sharing it)
//
// On most applications, "disaled" and "unavailable" is often treated as the
// same.
//
// On Google meet, failing to send video is grounds for terminating the
// participant from the call.
//
// We're just not going to do this. We're just going to warn all other
// participants that their media is available but something is wrong

export type MediaMeta = {
	readonly media: MediaStream | null;
	readonly deviceId: [string | null] | null;
};

export class RoomControl {
	private pool = new WeakMap<MediaStreamTrack, MediaStream>();

	constructor(private _name: string) {}

	static getMediaDevicesList() {
		return navigator.mediaDevices.enumerateDevices();
	}

	static async getMedia(
		kind: "audio" | "video",
		deviceId?: string | null
	): Promise<MediaMeta | null> {
		const media = await navigator.mediaDevices.getUserMedia({
			[kind]: deviceId ? { deviceId } : true,
		});

		const fnName = (kind[0].toUpperCase() + kind.slice(1)) as "Audio" | "Video";

		const [mediaTrack] = media[`get${fnName}Tracks`]();

		return { media, deviceId: [mediaTrack?.getSettings().deviceId || null] };
	}

	static async getVideo(deviceId?: string | null): Promise<MediaMeta | null> {
		const media = await navigator.mediaDevices.getUserMedia({
			// TODO: check to see what happens if a bad device ID is inserted
			...(deviceId ? {} : { video: deviceId ? { deviceId } : true }),
		});

		const [video] = media.getVideoTracks();
		return { media, deviceId: [video?.getSettings().deviceId || null] };
	}

	static async getAudio(deviceId?: string): Promise<MediaMeta | null> {
		const media = await navigator.mediaDevices.getUserMedia({
			// TODO: check to see what happens if a bad device ID is inserted
			...(deviceId ? {} : { video: deviceId ? { deviceId } : true }),
		});

		const [audio] = media.getAudioTracks();
		return { media, deviceId: [audio?.getSettings().deviceId || null] };
	}
}
