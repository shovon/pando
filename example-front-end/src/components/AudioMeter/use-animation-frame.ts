import { useEffect } from "react";

export function useAnimationFrame(callback: () => void) {
	useEffect(() => {
		let frame: number | undefined;

		const run = () => {
			frame = requestAnimationFrame(() => {
				callback();
				run();
			});
		};

		run();

		return () => {
			if (typeof frame !== "undefined") {
				cancelAnimationFrame(frame);
			}
		};
	}, [callback]);
}
