export class RoomControl {
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
