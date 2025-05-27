"use client"; // Ensure this is at the top for client-side hooks

import Link from 'next/link';
import React, { useEffect } from 'react'; // Import useEffect
import { useRouter } from 'next/navigation'; // Import useRouter

const Home: React.FC = () => {
  const router = useRouter();

  useEffect(() => {
    const token = localStorage.getItem('authToken');
    if (!token) {
      router.push('/login');
    }
  }, [router]); // Add router to dependency array

  // It might be good to return null or a loading spinner if redirecting,
  // to avoid briefly flashing the page content.
  const token = typeof window !== 'undefined' ? localStorage.getItem('authToken') : null;
  if (!token) {
    return null; // Or a loading component
  }

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
