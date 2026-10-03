'use client'

import { useState } from 'react'
import { useCartStore } from '@/store/useCartStore'
import { useStore } from '@/hooks/useStore'
import Link from 'next/link'

// Динамический URL API с фолбэком на внешний IP
const API_URL =
	process.env.NEXT_PUBLIC_API_URL || 'http://135.106.193.11:8080/api/v1'

export default function CartPage() {
	const [isSubmitting, setIsSubmitting] = useState(false)
	const [orderStatus, setOrderStatus] = useState<'idle' | 'success' | 'error'>(
		'idle',
	)

	// Безопасно получаем состояние из localStorage без ошибок hydration mismatch
	const items = useStore(useCartStore, state => state.items) ?? []
	const totalPrice = useStore(useCartStore, state => state.getTotalPrice()) ?? 0

	const updateQuantity = useCartStore(state => state.updateQuantity)
	const removeFromCart = useCartStore(state => state.removeFromCart)
	const clearCart = useCartStore(state => state.clearCart)

	const handleCheckout = async () => {
		if (items.length === 0) return

		setIsSubmitting(true)
		setOrderStatus('idle')

		try {
			const response = await fetch(`${API_URL}/orders`, {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
				},
				body: JSON.stringify({
					user_id: 1,
					items: items.map(item => ({
						product_id: Number(item.id),
						quantity: item.quantity,
					})),
				}),
			})

			if (!response.ok) {
				throw new Error('Ошибка при оформлении заказа')
			}

			setOrderStatus('success')
			clearCart()
		} catch (err) {
			console.error(err)
			setOrderStatus('error')
		} finally {
			setIsSubmitting(false)
		}
	}

	if (orderStatus === 'success') {
		return (
			<div className='max-w-2xl mx-auto my-12 p-8 text-center border rounded-xl shadow-sm bg-white'>
				<h2 className='text-2xl font-bold text-green-600 mb-2'>
					🎉 Заказ успешно оформлен!
				</h2>
				<p className='text-gray-600 mb-6'>
					Ваш заказ отправлен в обработку. Воркер Asynq подхватил задачу из
					Redis.
				</p>
				<Link
					href='/'
					className='inline-block bg-blue-600 hover:bg-blue-700 text-white font-medium px-6 py-2.5 rounded-lg transition'
				>
					Вернуться на главную
				</Link>
			</div>
		)
	}

	if (items.length === 0) {
		return (
			<div className='max-w-2xl mx-auto my-12 p-8 text-center border rounded-xl shadow-sm bg-white'>
				<h2 className='text-2xl font-bold mb-2'>Ваша корзина пуста</h2>
				<p className='text-gray-600 mb-6'>
					Выберите товары из каталога, чтобы продолжить.
				</p>
				<Link
					href='/'
					className='inline-block bg-blue-600 hover:bg-blue-700 text-white font-medium px-6 py-2.5 rounded-lg transition'
				>
					Перейти к товарам
				</Link>
			</div>
		)
	}

	return (
		<div className='max-w-4xl mx-auto py-8'>
			<h1 className='text-2xl font-bold mb-6'>Корзина покупок</h1>

			<div className='grid grid-cols-1 md:grid-cols-3 gap-8'>
				<div className='md:col-span-2 space-y-4'>
					{items.map(item => (
						<div
							key={item.id}
							className='flex items-center justify-between border p-4 rounded-lg bg-white shadow-sm'
						>
							<div>
								<h3 className='font-semibold text-lg'>{item.name}</h3>
								<p className='text-gray-600'>${item.price} за шт.</p>
							</div>

							<div className='flex items-center gap-4'>
								<div className='flex items-center border rounded-lg overflow-hidden'>
									<button
										onClick={() => updateQuantity(item.id, item.quantity - 1)}
										className='px-3 py-1 bg-gray-100 hover:bg-gray-200 transition'
									>
										-
									</button>
									<span className='px-4 font-medium'>{item.quantity}</span>
									<button
										onClick={() => updateQuantity(item.id, item.quantity + 1)}
										className='px-3 py-1 bg-gray-100 hover:bg-gray-200 transition'
									>
										+
									</button>
								</div>

								<button
									onClick={() => removeFromCart(item.id)}
									className='text-red-500 hover:text-red-700 text-sm font-medium'
								>
									Удалить
								</button>
							</div>
						</div>
					))}
				</div>

				<div className='border p-6 rounded-lg bg-white shadow-sm h-fit'>
					<h2 className='text-lg font-bold mb-4'>Итого</h2>
					<div className='flex justify-between text-xl font-bold mb-6'>
						<span>Сумма:</span>
						<span>${totalPrice.toFixed(2)}</span>
					</div>

					{orderStatus === 'error' && (
						<p className='text-red-600 text-sm mb-4'>
							Произошла ошибка при отправке. Проверьте подключение к API.
						</p>
					)}

					<button
						onClick={handleCheckout}
						disabled={isSubmitting}
						className='w-full bg-blue-600 hover:bg-blue-700 disabled:bg-blue-300 text-white font-bold py-3 px-4 rounded-lg transition'
					>
						{isSubmitting ? 'Оформить заказ...' : 'Оформить заказ'}
					</button>
				</div>
			</div>
		</div>
	)
}
