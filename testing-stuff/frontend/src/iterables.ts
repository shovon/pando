export function map<T, V>(it: Iterable<T>, mapping: (v: T) => V): Iterable<V> {
	return {
		*[Symbol.iterator]() {
			for (const item of it) {
				yield mapping(item);
			}
		},
	};
}
