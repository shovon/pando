export type Pipe<T> = {
	_<V>(fn: (value: T) => V): Pipe<V>;
	readonly value: T;
};

export const start = <T>(initial: T): Pipe<T> => ({
	_: <V>(fn: (value: T) => V): Pipe<V> => start(fn(initial)),
	value: initial,
});
