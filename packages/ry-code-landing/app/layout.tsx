import type { Metadata } from 'next'
import './globals.css'

export const metadata: Metadata = {
  metadataBase: new URL('https://ry-code.com'),
  title: {
    default: 'RyCode - World\'s Most Advanced Open Source Coding Agent | 5 SOTA AI Models',
    template: '%s | RyCode'
  },
  description: 'The world\'s most advanced open source coding agent. Switch between 5 state-of-the-art AI models (Claude 4.5 Sonnet, Gemini 3.0 Pro, Codex, Grok Code Fast, Qwen 3 Coder) with a single keystroke. Zero context loss. 60 FPS terminal UI. Production-ready with 31/31 tests passing.',
  keywords: [
    'AI CLI',
    'multi-agent AI',
    'Claude CLI',
    'Gemini CLI',
    'AI code assistant',
    'developer tools',
    'AI terminal',
    'Claude 4.5 Sonnet',
    'Gemini 3.0 Pro',
    'Codex',
    'Grok Code Fast',
    'Qwen 3 Coder',
    'AI model switching',
    'command line AI',
    'toolkit-cli',
    'AI powered development',
    'multi-model AI assistant'
  ],
  authors: [
    { name: 'RyCode Team', url: 'https://ry-code.com' },
    { name: 'Claude AI' }
  ],
  creator: 'toolkit-cli',
  publisher: 'RyCode',
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
  openGraph: {
    type: 'website',
    locale: 'en_US',
    url: 'https://ry-code.com',
    siteName: 'RyCode',
    title: 'RyCode - World\'s Most Advanced Open Source Coding Agent',
    description: 'Switch between 5 SOTA AI models (Claude 4.5 Sonnet, Gemini 3.0 Pro, Codex, Grok Code Fast, Qwen 3 Coder) with a single keystroke. Zero context loss. Production-ready.',
    images: [
      {
        url: '/og-image.png',
        width: 1200,
        height: 630,
        alt: 'RyCode - World\'s Most Advanced Open Source Coding Agent with 5 SOTA Models',
        type: 'image/png',
      }
    ],
  },
  twitter: {
    card: 'summary_large_image',
    site: '@rycode',
    creator: '@rycode',
    title: 'RyCode - World\'s Most Advanced Open Source Coding Agent',
    description: 'Switch between 5 SOTA AI models with one keystroke. Zero context loss. 60 FPS terminal UI. Production-ready with 31/31 tests passing.',
    images: ['/twitter-image.png'],
  },
  alternates: {
    canonical: 'https://ry-code.com',
  },
  category: 'technology',
  classification: 'Developer Tools',
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="en">
      <head>
        {/* Schema.org JSON-LD */}
        <script
          type="application/ld+json"
          dangerouslySetInnerHTML={{
            __html: JSON.stringify({
              '@context': 'https://schema.org',
              '@type': 'SoftwareApplication',
              name: 'RyCode',
              applicationCategory: 'DeveloperApplication',
              operatingSystem: 'macOS, Linux, Windows',
              offers: {
                '@type': 'Offer',
                price: '0',
                priceCurrency: 'USD',
              },
              description: 'The world\'s most advanced open source coding agent. Multi-agent AI CLI that allows developers to switch between 5 state-of-the-art AI models (Claude 4.5 Sonnet, Gemini 3.0 Pro, Codex, Grok Code Fast, Qwen 3 Coder) with a single keystroke. Zero context loss, 60 FPS terminal UI, production-ready with 31/31 tests passing.',
              screenshot: 'https://ry-code.com/assets/splash_to_selector.gif',
              featureList: [
                '5 State-of-the-Art AI Models',
                'Instant Model Switching with Tab Key',
                'Context Preservation Across Models',
                'Professional Terminal UI',
                'Built with toolkit-cli'
              ],
              aggregateRating: {
                '@type': 'AggregateRating',
                ratingValue: '5.0',
                ratingCount: '1',
                bestRating: '5',
                worstRating: '1'
              },
            })
          }}
        />
        {/* Additional Schema for Organization */}
        <script
          type="application/ld+json"
          dangerouslySetInnerHTML={{
            __html: JSON.stringify({
              '@context': 'https://schema.org',
              '@type': 'Organization',
              name: 'RyCode',
              url: 'https://ry-code.com',
              logo: 'https://ry-code.com/logo.png',
              sameAs: [
                'https://github.com/aaronmrosenthal/RyCode',
              ],
              contactPoint: {
                '@type': 'ContactPoint',
                contactType: 'Technical Support',
                url: 'https://ry-code.com'
              }
            })
          }}
        />
      </head>
      <body>{children}</body>
    </html>
  )
}
