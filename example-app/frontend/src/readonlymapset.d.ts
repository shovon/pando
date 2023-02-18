interface ReadOnlyMapNoGetter<K, V> {
	has(key: K): boolean;
	entries(): IterableIterator<[K, V]>;
	forEach(cb: (value: V, key: K, map: ReadOnlyMapNoGetter<K, V>) => void): void;
	keys(): IterableIterator<K>;
	values(): IterableIterator<V>;
	readonly size: number;
}

interface ReadOnlyMap<K, V> extends ReadOnlyMapNoGetter<K, V> {
	get(key: K): V | undefined;
	[Symbol.iterator](): IterableIterator<[K, V]>;
}

interface ReadOnlySet<V> extends ReadOnlyMapNoGetter<V, V> {
	[Symbol.iterator](): IterableIterator<V>;
}
