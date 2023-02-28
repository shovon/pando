import { Sub } from "@sparkscience/wskeyid-browser/src/pub-sub";

export type ExtractSubParams<T> = T extends Sub<infer P> ? P : never;

export function map<T, V>(sub: Sub<T>, mapping: (v: T) => V): Sub<V> {
	return {
		addEventListener(listener: (v: V) => void) {
			return sub.addEventListener((v) => listener(mapping(v)));
		},
	};
}

export function filter<T, V extends T = T>(
	sub: Sub<T>,
	predicate: (v: T) => boolean
): Sub<V> {
	return {
		addEventListener(listener: (v: V) => void) {
			return sub.addEventListener((v) => {
				if (predicate(v)) {
					listener(v as V);
				}
			});
		},
	};
}

export function toSubscribable(sub: Sub<any>): {
	subscribe: (listener: (v: any) => void) => () => void;
} {
	return {
		subscribe(listener: (v: any) => void) {
			return sub.addEventListener(listener);
		},
	};
}
