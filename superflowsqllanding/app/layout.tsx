import React from "react"
import type { Metadata } from 'next'
import { Inter, JetBrains_Mono } from 'next/font/google'
import { Analytics } from '@vercel/analytics/next'
import './globals.css'

const inter = Inter({ 
  subsets: ['latin'],
  weight: ['300', '400', '500', '600', '700'],
  variable: '--font-inter',
});

const jetbrainsMono = JetBrains_Mono({ 
  subsets: ['latin'],
  weight: ['400', '500'],
  variable: '--font-jetbrains',
});

export const metadata: Metadata = {
  title: 'SuperFlowSQL | Complete Data Orchestration Platform',
  description: 'The complete platform integrating Airflow, PostgreSQL, PgAdmin, and Superset. Automate pipelines, manage data, and visualize results in one unified environment.',
  icons: {
    icon: [
      {
        url: '/icon-light-32x32.ico',
        media: '(prefers-color-scheme: light)',
      },
      {
        url: '/icon-dark-32x32.ico',
        media: '(prefers-color-scheme: dark)',
      },
      {
        url: '/icon.svg',
        type: 'image/svg+xml',
      },
    ],
    apple: '/apple-icon.png',
  },
}

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode
}>) {
  return (
    <html lang="en" className={`scroll-smooth ${inter.variable} ${jetbrainsMono.variable}`}>
      <body className="bg-white text-slate-900 font-sans antialiased overflow-x-hidden">
        {children}
        <Analytics />
      </body>
    </html>
  )
}
