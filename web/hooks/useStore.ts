import { useSyncExternalStore } from 'react'

export function useStore<T, F>(
	store: {
		(callback: (state: T) => unknown): unknown
		getState: () => T
		subscribe: (listener: () => void) => () => void
	},
	callback: (state: T) => F,
): F | undefined {
	return useSyncExternalStore(
		listener => store.subscribe(listener),
		() => callback(store.getState()),
		() => undefined, // На сервере (SSR) возвращаем undefined до гидратации
	)
}
