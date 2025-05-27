"use client"; // Add this

import Link from 'next/link';
import React, { useEffect } from 'react'; // Import useEffect
import { useRouter } from 'next/navigation'; // Import useRouter

const Home: React.FC = () => {
  const router = useRouter();

  useEffect(() => {
    const token = localStorage.getItem('authToken');
    if (!token) {
      router.replace('/login');
    }
    // If token exists, user stays on the page.
    // Optionally, you could redirect to '/workouts' here if a token exists
    // else {
    //   router.replace('/workouts');
    // }
  }, [router]); // Add router to dependency array

  // The content will be rendered if the user is authenticated (token exists)
  // and not redirected away. If no token, there might be a brief flash
  // of this content before redirection.
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
