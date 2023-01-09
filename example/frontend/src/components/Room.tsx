import { useEffect, useState } from "react";
import { AudioMeter } from "./AudioMeter/AudioMeter";
import { StreamPlayer } from "./StreamPlayer";
import { RoomControl } from "./room-control";
import { CallMediaSelector } from "./CallMediaSelector";

/**
 * Represents the room in a call
 * @returns JSX component that represents the DOM layout of the room
 */
export function Room() {
	return <CallMediaSelector />;
}
