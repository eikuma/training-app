// app/layout.tsx
"use client"; // Required for ThemeProvider and CssBaseline if theme is not serializable or for context
import type { Metadata } from 'next'; // Still can be used for metadata
import Layout from '../components/Layout';
import theme from './theme'; // Import the created theme
import { ThemeProvider } from '@mui/material/styles';
import CssBaseline from '@mui/material/CssBaseline';

// export const metadata: Metadata = { // Metadata can be exported from client components in Next.js 13+ App Router
//   title: 'ジムのトレーニング記録・提案アプリ',
// };
// If you need dynamic metadata or metadata that depends on client-side state,
// you might need to handle it differently, but for static metadata, this is fine.

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="ja">
      <head>
        {/* Next.js handles title and other head elements, 
            but if you were setting static metadata here, it would go here.
            The metadata object above is the preferred way for static metadata. */}
        <title>ジムのトレーニング記録・提案アプリ</title>
        <meta name="viewport" content="initial-scale=1, width=device-width" />
      </head>
      <body>
        <ThemeProvider theme={theme}>
          <CssBaseline /> {/* Normalizes styles and applies background color from theme */}
          <Layout>{children}</Layout>
        </ThemeProvider>
      </body>
    </html>
  );
}
