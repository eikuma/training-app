// app/layout.tsx
import type { Metadata } from 'next';

export const metadata: Metadata = {
  title: 'ジムのトレーニング記録・提案アプリ',
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="ja">
      <body>{children}</body>
    </html>
  );
}
