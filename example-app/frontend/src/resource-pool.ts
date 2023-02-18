/**
 * The object that will aid in overseeing the allocation and deallocation of
 * resources.
 */
type Allocator<K, V> = {
	allocate: (key: K) => V;
	deallocate: (key: K, value: V) => void;
};

/**
 * Creates a resource pool for allocating values associated with a key.
 *
 * Usage:
 *
 *     // An example associating a MediaStream to a pool of `HTMLVideoElement`s
 *     const pool = createWeakResourcePool<MediaStream, HTMLVideoElement>(
 *       () => document.createElement('video')
 *     );
 *
 *     // For some `stream` that we were given access to.
 *     const video = pool.allocate(stream);
 *
 *     // Later…
 *
 *     pool.deallocate(stream, video);
 *
 * The benefit of the resource pool is that if you lose access to the key, the
 * resource associated with the key will garbage collect itself, and thus
 * avoiding memory and resource leaks.
 *
 * However, this is not a guarantee against memory leaks. If the key is still
 * being referenced, all allocated (but unused) resources will remain in memory.
 * And thus, care needs to be taken when invoking the `allocate` method.
 * Additionally, the `deallocate` method does not delete a resource; it merely
 * "flags" it as "unused", so that the `allocate` method will opt for returning
 * the "unused" resource, rather than to initialize an entirely new resource.
 *
 * Only non-primitive reference types allowed as keys (i.e. objects, functions,
 * regular expressions or anything initialized using the `new` keyword are
 * perfectly allowed). What is NOT allowed are primitive types, which are
 * booleans, numbers, and strings.
 * @param create The function that will be used for creating the item that will
 *   be placed in the pool
 */
export function createWeakResourcePool<K extends object, V>(
	create: () => V
): Allocator<K, V> {
	// The WeakMap is used for clearing out all resources for keys that are no
	// longer needed.
	const map = new WeakMap<K, Set<V>>();

	return {
		allocate: (key: K): V => {
			// No resource associated with the key.
			if (!map.has(key)) {
				map.set(key, new Set());
			}

			// Get the list of resources associated with the key.
			const freeList = map.get(key);
			if (!freeList) {
				throw new Error("Should not be here.");
			}

			// If there are no free resources, allocate a new resource, and return it.
			if (freeList.size <= 0) {
				const value = create();
				return value;
			}

			// Get some random item from the list, and delete it.
			const value = freeList.values().next().value;
			freeList.delete(value);

			// Return the retrieved random item.
			return value;
		},
		deallocate: (key: K, value: V) => {
			const freeList = map.get(key);
			if (!freeList) {
				console.error("Something is not quite right.");
				return;
			}
			freeList.add(value);
		},
	};
}
