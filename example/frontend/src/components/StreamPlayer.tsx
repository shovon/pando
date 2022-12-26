import React, { useEffect, useRef } from "react";
import { createWeakResourcePool } from "../resource-pool";

const pool = createWeakResourcePool<
	MediaStream | { src: string },
	HTMLVideoElement
>((): HTMLVideoElement => {
	const video = document.createElement("video");
	return video;
});

export type StreamPlayerProps = {
	/**
	 * The MediaStream that the stream player will play
	 */
	stream: MediaStream | { src: string };

	/**
	 * The CSS properties for styling the stream player.
	 */
	style?: React.CSSProperties;

	/**
	 * Is the video muted?
	 */
	muted?: boolean;

	/**
	 * This is for determining how the video will be rendered, if it doesn't fit
	 * the square.
	 */
	objectFit?: "fill" | "contain" | "cover" | "none" | "scale-down";

	/**
	 * Does the video loop?
	 */
	loop?: boolean;
};

/**
 * A component to play video/audio streams.
 * @param props The props for the stream player.
 */
export function StreamPlayer({
	stream,
	style,
	muted,
	objectFit,
	loop,
}: StreamPlayerProps) {
	const streamRef = useRef<MediaStream | { src: string }>(stream);
	const videoRef = useRef<HTMLVideoElement | null>(null);

	useEffect(() => {
		const previousStream = streamRef.current;
		const previousVideo = videoRef.current;
		return () => {
			if (previousVideo) {
				pool.deallocate(previousStream, previousVideo);
			}
		};
	});

	return (
		<div
			style={style}
			ref={(ref) => {
				const previousStream = streamRef.current;
				streamRef.current = stream;
				if (videoRef.current) {
					videoRef.current.remove();
					pool.deallocate(previousStream, videoRef.current);
				}
				if (ref) {
					const video = pool.allocate(stream);
					video.style.objectFit = objectFit ? objectFit : "cover";
					video.style.width = "100%";
					video.style.height = "100%";
					if (muted) {
						video.muted = true;
					}
					video.onclick = (e) => {
						e.preventDefault();
						e.stopPropagation();
					};
					video.setAttribute("playsinline", "playsinline");
					video.addEventListener("contextmenu", (e) => {
						e.preventDefault();
					});
					videoRef.current = video;
					if (!video.srcObject) {
						if (stream instanceof MediaStream) {
							video.srcObject = stream;
						} else {
							video.src = stream.src;
						}
					}
					if (loop) {
						video.loop = true;
					}
					video.play();

					ref.appendChild(video);
				}
			}}
		/>
	);
}
