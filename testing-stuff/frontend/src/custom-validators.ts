import { arrayOf, transform, tuple, chain, Validator } from "./validator";

export const iterableOf = <T>(validator: Validator<T>) =>
	chain(
		transform((it) => [...it]),
		arrayOf(validator)
	);

export const kvToReadOnlyMap = <K, V>(key: Validator<K>, value: Validator<V>) =>
	chain(
		iterableOf(tuple([key, value])),
		transform<ReadOnlyMap<K, V>>((arr) => new Map<K, V>(arr))
	);
