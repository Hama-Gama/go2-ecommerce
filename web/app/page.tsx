'use client'

import { useState, useEffect } from 'react'
import { Product, OrderResponse } from '@/types'
import { useCartStore, CartItem } from '@/store/useCartStore'

// Динамический URL API с фолбэком на внешний IP
const API_URL =
	process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1'

export default function Home() {
	const [products, setProducts] = useState<Product[]>([])
	const [search, setSearch] = useState('')
	const [loading, setLoading] = useState(false)
	const [orderSuccess, setOrderSuccess] = useState<OrderResponse | null>(null)

	// Zustand state and actions
	const items = useCartStore(state => state.items)
	const addToCart = useCartStore(state => state.addToCart)
	const removeFromCart = useCartStore(state => state.removeFromCart)
	const clearCart = useCartStore(state => state.clearCart)
	const getTotalPrice = useCartStore(state => state.getTotalPrice)

	useEffect(() => {
		const fetchProducts = async () => {
			setLoading(true)
			try {
				const queryParams = search.trim()
					? `?query=${encodeURIComponent(search)}`
					: ''

				const res = await fetch(`${API_URL}/products${queryParams}`)
				if (res.ok) {
					const data = await res.json()
					const list = Array.isArray(data) ? data : data.products || []
					setProducts(list)
				} else {
					setProducts([])
				}
			} catch (err) {
				console.error('Failed to fetch products:', err)
				setProducts([])
			} finally {
				setLoading(false)
			}
		}

		const timer = setTimeout(fetchProducts, 300)
		return () => clearTimeout(timer)
	}, [search])

	const handleCheckout = async () => {
		if (items.length === 0) return

		const payload = {
			user_id: 1,
			items: items.map((item: CartItem) => ({
				product_id: Number(item.id),
				quantity: item.quantity,
			})),
		}

		try {
			const res = await fetch(`${API_URL}/orders`, {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify(payload),
			})

			if (res.ok) {
				const orderData: OrderResponse = await res.json()
				setOrderSuccess(orderData)
				clearCart()
			} else {
				alert('Ошибка при создании заказа')
			}
		} catch (err) {
			console.error('Checkout failed:', err)
		}
	}

	const handleAddToCart = (p: Product) => {
		addToCart({
			id: String(p.id),
			name: p.title,
			price: Number(p.price),
		})
	}

	return (
		<main className='min-h-screen bg-gray-50 p-8 font-sans text-gray-900'>
			<div className='max-w-6xl mx-auto'>
				<header className='flex justify-between items-center mb-8 border-b pb-4'>
					<h1 className='text-3xl font-bold text-slate-800'>
						Go E-Commerce Store
					</h1>
					<div className='text-sm bg-blue-50 text-blue-700 px-3 py-1 rounded-full border border-blue-200'>
						Backend Status:{' '}
						<span className='font-semibold'>Connected (:8080)</span>
					</div>
				</header>

				{orderSuccess && (
					<div className='mb-6 p-4 bg-green-100 border border-green-400 text-green-700 rounded-lg'>
						🎉 <strong>Заказ #{orderSuccess.id} успешно создан!</strong> Сумма:
						${orderSuccess.total_amount}. Фоновая задача отправлена в Asynq.
					</div>
				)}

				<div className='grid grid-cols-1 lg:grid-cols-3 gap-8'>
					<div className='lg:col-span-2'>
						<div className='mb-6'>
							<input
								type='text'
								placeholder='🔍 Поиск товаров...'
								value={search}
								onChange={e => setSearch(e.target.value)}
								className='w-full p-3 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 bg-white'
							/>
						</div>

						{loading ? (
							<div className='text-gray-500 py-8 text-center'>
								Загрузка товаров...
							</div>
						) : (
							<div className='grid grid-cols-1 md:grid-cols-2 gap-4'>
								{!products || products.length === 0 ? (
									<div className='col-span-2 text-gray-400 py-8 text-center'>
										Товары не найдены
									</div>
								) : (
									products.map(p => (
										<div
											key={p.id}
											className='bg-white p-5 rounded-xl shadow-sm border hover:shadow-md transition'
										>
											<h3 className='font-semibold text-lg text-slate-800'>
												{p.title}
											</h3>
											<p className='text-sm text-gray-500 mt-1 line-clamp-2'>
												{p.description}
											</p>
											<div className='mt-4 flex justify-between items-center'>
												<div>
													<span className='text-xl font-bold text-slate-900'>
														${p.price}
													</span>
												</div>
												<button
													onClick={() => handleAddToCart(p)}
													className='bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg text-sm transition'
												>
													В корзину
												</button>
											</div>
										</div>
									))
								)}
							</div>
						)}
					</div>

					<div className='bg-white p-6 rounded-xl shadow-sm border h-fit sticky top-8'>
						<h2 className='text-xl font-bold mb-4 text-slate-800 border-b pb-2'>
							Корзина
						</h2>
						{items.length === 0 ? (
							<p className='text-gray-400 text-sm py-4'>Корзина пуста</p>
						) : (
							<div className='space-y-4'>
								{items.map((item: CartItem) => (
									<div
										key={item.id}
										className='flex justify-between items-center border-b pb-3'
									>
										<div>
											<h4 className='font-medium text-sm text-slate-800'>
												{item.name}
											</h4>
											<div className='text-xs text-gray-500'>
												${item.price} × {item.quantity}
											</div>
										</div>
										<button
											onClick={() => removeFromCart(item.id)}
											className='text-red-500 text-xs hover:underline'
										>
											Удалить
										</button>
									</div>
								))}

								<div className='pt-2 border-t flex justify-between font-bold text-lg text-slate-900'>
									<span>Итого:</span>
									<span>${getTotalPrice().toFixed(2)}</span>
								</div>

								<button
									onClick={handleCheckout}
									className='w-full bg-emerald-600 hover:bg-emerald-700 text-white font-medium py-3 rounded-lg transition'
								>
									Оформить заказ
								</button>
							</div>
						)}
					</div>
				</div>
			</div>
		</main>
	)
}
