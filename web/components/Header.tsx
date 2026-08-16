'use client'

import { CartBadge } from './CartBadge'

export function Header() {
	return (
		<header className='border-b p-4 flex justify-between items-center bg-white'>
			<h1 className='text-xl font-bold'>Go E-Commerce</h1>
			<CartBadge />
		</header>
	)
}
