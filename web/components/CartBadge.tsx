'use client'

import { useCartStore } from '../store/useCartStore'
import { useStore } from '../hooks/useStore'
import Link from 'next/link'

export function CartBadge() {
	const totalItems = useStore(useCartStore, state => state.getTotalItems()) ?? 0

	return (
		<Link
			href='/cart'
			className='flex items-center gap-2 cursor-pointer font-medium hover:text-blue-600 transition'
		>
			<span>🛒 Корзина</span>
			{totalItems > 0 && (
				<span className='bg-blue-600 text-white text-xs font-bold px-2 py-0.5 rounded-full'>
					{totalItems}
				</span>
			)}
		</Link>
	)
}
