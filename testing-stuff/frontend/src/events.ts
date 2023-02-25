/*
Copyright 2022 Salehen Shovon Rahman
Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies
of the Software, and to permit persons to whom the Software is furnished to do
so, subject to the following conditions:
The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.
THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

/**
 * A generic type that represents a function whose sole purpose is to receive
 * a value.
 *
 * Useful when listening for an event
 */
export type Listener<T> = (value: T) => void;

/**
 * A subscribable event listener
 */
export type Subscribable<T> = {
	subscribe(listener: Listener<T>): () => void;
};

export type Subscriber<T> = {
	emit(value: T): void;
};

export type OperatorFunction<T, V> = (
	subscribable: Subscribable<T>
) => Subscribable<V>;

export const createSubject = <T>(): Subscribable<T> & Subscriber<T> => {
	let listeners: Listener<T>[] = [];

	return {
		subscribe: (listener: Listener<T>) => {
			listeners.push(listener);
			return () => {
				listeners = listeners.filter((l) => l !== listener);
			};
		},
		emit: (value: T) => {
			for (const listener of listeners) {
				listener(value);
			}
		},
	};
};

export const createSubscribable = <T>(
	fn: (subscriber: Subscriber<T>) => void
): Subscribable<T> => {
	const { subscribe, emit } = createSubject<T>();
	fn({ emit });
	return { subscribe };
};
