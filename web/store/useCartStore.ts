import { create } from 'zustand'
import { persist, createJSONStorage } from 'zustand/middleware'

export interface CartItem {
	id: string
	name: string
	price: number
	quantity: number
	image_url?: string
}

interface CartStore {
	items: CartItem[]
	addToCart: (product: {
		id: string
		name: string
		price: number
		image_url?: string
	}) => void
	removeFromCart: (productId: string) => void
	updateQuantity: (productId: string, quantity: number) => void
	clearCart: () => void
	getTotalPrice: () => number
	getTotalItems: () => number
}

export const useCartStore = create<CartStore>()(
	persist(
		(set, get) => ({
			items: [],

			addToCart: product => {
				const currentItems = get().items
				const existingItem = currentItems.find(item => item.id === product.id)

				if (existingItem) {
					set({
						items: currentItems.map(item =>
							item.id === product.id
								? { ...item, quantity: item.quantity + 1 }
								: item,
						),
					})
				} else {
					set({
						items: [...currentItems, { ...product, quantity: 1 }],
					})
				}
			},

			removeFromCart: productId => {
				set({
					items: get().items.filter(item => item.id !== productId),
				})
			},

			updateQuantity: (productId, quantity) => {
				if (quantity <= 0) {
					get().removeFromCart(productId)
					return
				}
				set({
					items: get().items.map(item =>
						item.id === productId ? { ...item, quantity } : item,
					),
				})
			},

			clearCart: () => set({ items: [] }),

			getTotalPrice: () => {
				return get().items.reduce(
					(total, item) => total + item.price * item.quantity,
					0,
				)
			},

			getTotalItems: () => {
				return get().items.reduce((total, item) => total + item.quantity, 0)
			},
		}),
		{
			name: 'shopping-cart-storage',
			storage: createJSONStorage(() => localStorage),
		},
	),
)
