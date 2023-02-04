export const dispose: unique symbol = Symbol('dispose');

export interface Disposable {
	[dispose](): void
}
