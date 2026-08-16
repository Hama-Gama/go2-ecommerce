'use client'

import { useCartStore } from '../store/useCartStore'

interface Product {
	id: string
	name: string
	price: number
	image_url?: string
}

export function ProductCard({ product }: { product: Product }) {
	const addToCart = useCartStore(state => state.addToCart)

	return (
		<div className='border p-4 rounded-lg shadow flex flex-col justify-between bg-white'>
			<div>
				<h3 className='font-bold text-lg'>{product.name}</h3>
				<p className='text-gray-600'>${product.price}</p>
			</div>
			<button
				onClick={() => addToCart(product)}
				className='mt-4 bg-blue-600 hover:bg-blue-700 text-white font-medium py-2 px-4 rounded transition'
			>
				В корзину
			</button>
		</div>
	)
}
