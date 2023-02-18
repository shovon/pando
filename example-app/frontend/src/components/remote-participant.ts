import {
	string,
	object,
	transform,
	chain,
	tuple,
	either,
	exact,
	Validator,
	boolean,
	ValidationError,
	InferType,
	any,
} from "../validator";

class NotIterableError extends ValidationError {
	constructor(value: any) {
		super(
			"Not an iterable error",
			"The supplied value is not an iterable",
			value
		);
	}
}

class IterableValueIsNotValid extends ValidationError {
	constructor(it: Iterable<any>, private _value: any) {
		super(
			"Bad iterable value error",
			"A value in the given iterable is not valid",
			it
		);
	}

	get theValue() {
		return this._value;
	}
}

const iterableOf = <V>(validator: Validator<V>): Validator<Iterable<V>> => ({
	validate: (value: any) => {
		if (!value[Symbol.iterator]) {
			return { isValid: false, error: new NotIterableError(value) };
		}
		for (const v of value) {
			const validation = validator.validate(v);
			if (!validation.isValid) {
				return { isValid: false, error: new IterableValueIsNotValid(value, v) };
			}
		}

		return { isValid: true, value };
	},
});

const iterableToReadOnlyMap = <K, V>(
	key: Validator<K>,
	value: Validator<V>
): Validator<ReadOnlyMap<K, V>> =>
	chain(
		iterableOf(tuple([key, value])),
		transform((v) => new Map(v))
	);

const readOnly = <T>(v: Validator<T>): Validator<Readonly<T>> => v;

const nullable = <T>(validator: Validator<T>) => either(validator, exact(null));

const mediaObject = object({
	disabled: boolean(),
	sources: iterableToReadOnlyMap(string(), any()),
});

export const remoteParticipant = readOnly(
	object({
		name: nullable(string()),
		audio: mediaObject,
		video: mediaObject,
		screenShare: mediaObject,
	})
);

export type RemoteParticipant = InferType<typeof remoteParticipant>;
