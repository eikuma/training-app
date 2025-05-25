import Link from 'next/link';
import React from 'react';

const Home: React.FC = () => {
  return (
    <div style={{ textAlign: 'center', marginTop: '2rem' }}>
      <h1>ジムのトレーニング記録アプリ</h1>
      <div style={{ margin: '1rem' }}>
        <Link href="/workouts" passHref>
          <button style={{ padding: '0.5rem 1rem', fontSize: '1rem' }}>記録ページへ</button>
        </Link>
      </div>
      <div style={{ margin: '1rem' }}>
        <Link href="/recommendations" passHref>
          <button style={{ padding: '0.5rem 1rem', fontSize: '1rem' }}>提案ページへ</button>
        </Link>
      </div>
    </div>
  );
};

export default Home;
