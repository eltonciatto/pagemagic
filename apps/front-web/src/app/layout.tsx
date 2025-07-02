import type { Metadata } from 'next';
import { Inter, Geist_Mono } from 'next/font/google';
import { ThemeProvider } from 'next-themes';
import { Providers } from '@/components/providers';
import './globals.css';

const inter = Inter({
  subsets: ['latin'],
  variable: '--font-inter',
});

const geistMono = Geist_Mono({
  subsets: ['latin'],
  variable: '--font-geist-mono',
});

export const metadata: Metadata = {
  title: {
    default: 'Page Magic - Criador de Sites com IA',
    template: '%s | Page Magic',
  },
  description: 'Crie sites profissionais em minutos usando inteligência artificial. Sem código, sem complicação.',
  keywords: ['no-code', 'website builder', 'ai', 'inteligência artificial', 'criador de sites'],
  authors: [{ name: 'Page Magic Team' }],
  creator: 'Page Magic',
  metadataBase: new URL('https://pagemagic.io'),
  openGraph: {
    type: 'website',
    locale: 'pt_BR',
    url: 'https://pagemagic.io',
    title: 'Page Magic - Criador de Sites com IA',
    description: 'Crie sites profissionais em minutos usando inteligência artificial.',
    siteName: 'Page Magic',
    images: [
      {
        url: '/og-image.png',
        width: 1200,
        height: 630,
        alt: 'Page Magic',
      },
    ],
  },
  twitter: {
    card: 'summary_large_image',
    title: 'Page Magic - Criador de Sites com IA',
    description: 'Crie sites profissionais em minutos usando inteligência artificial.',
    images: ['/og-image.png'],
    creator: '@pagemagic',
  },
  robots: {
    index: true,
    follow: true,
    googleBot: {
      index: true,
      follow: true,
      'max-video-preview': -1,
      'max-image-preview': 'large',
      'max-snippet': -1,
    },
  },
  verification: {
    google: 'google-site-verification-code',
  },
  manifest: '/manifest.json',
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html 
      lang="pt-BR" 
      className={`${inter.variable} ${geistMono.variable}`}
      suppressHydrationWarning
    >
      <body className="min-h-screen bg-background font-sans antialiased">
        <Providers>
          <div className="relative flex min-h-screen flex-col">
            <div className="flex-1">{children}</div>
          </div>
        </Providers>
      </body>
    </html>
  );
}
