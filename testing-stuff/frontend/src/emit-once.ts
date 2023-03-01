export type OnceEmitter<T> = {
	readonly value: { isEmitted: false } | { isEmitted: true; value: T };
	addEventListener(listener: (value: T) => void): () => void;
};

export type Emittable<T> = {
	emit(value: T): void;
};

export type SubjectOnce<T> = OnceEmitter<T> & Emittable<T>;

/**
 * Creates an event emitter that will only ever be emitted once
 * @returns A subject that can only be emitted once
 */
export function createEmitOnce<T>(): SubjectOnce<T> {
	const listeners = new Set<(value: T) => void>();
	let emittedValue: [T] | null = null;

	return {
		get value(): { isEmitted: false } | { isEmitted: true; value: T } {
			if (emittedValue) {
				return { isEmitted: true, value: emittedValue[0] };
			}
			return { isEmitted: false };
		},
		addEventListener(listener: (value: T) => void) {
			if (emittedValue) {
				listener(emittedValue[0]);
				return () => {};
			}
			listeners.add(listener);
			return () => {
				listeners.delete(listener);
			};
		},
		emit(value: T) {
			if (value) {
				console.error("Already emitted!");
			}
			emittedValue = [value];
			for (const listener of listeners) {
				listener(value);
			}
		},
	};
}
