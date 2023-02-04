type Listener<T> = (event: T) => void;

export interface Sub<T> {
	addEventListener(listener: Listener<T>): () => void;
}

export const filter =
	<T>(pred: (v: T) => boolean): ((s: Sub<T>) => Sub<T>) =>
	(s: Sub<T>): Sub<T> => ({
		addEventListener(listener: Listener<T>) {
			return s.addEventListener((e) => {
				if (pred(e)) {
					listener(e);
				}
			});
		},
	});
