import type { Metadata } from 'next'
import { Inter, Geist } from 'next/font/google'
import './globals.css'
import { CartProvider } from '@/context/CartContext'
import { cn } from '@/lib/utils'
import { Header } from '@/components/Header'

const geist = Geist({ subsets: ['latin'], variable: '--font-sans' })
const inter = Inter({ subsets: ['latin'] })

export const metadata: Metadata = {
	title: 'Go E-Commerce Store',
	description: 'Next.js Frontend for Go Backend',
}

export default function RootLayout({
	children,
}: {
	children: React.ReactNode
}) {
	return (
		<html lang='en' className={cn('font-sans', geist.variable)}>
			<body className={inter.className}>
				<CartProvider>
					<Header /> {/* <-- 2. Вставили Header над основным контентом */}
					<main className='container mx-auto p-4'>{children}</main>
				</CartProvider>
			</body>
		</html>
	)
}
