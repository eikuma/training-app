"use client";

import Link from 'next/link';
import React, { useState, useEffect } from 'react'; // Import useState and useEffect
import { useRouter } from 'next/navigation';

const Home: React.FC = () => {
  const router = useRouter();
  const [isLoading, setIsLoading] = useState(true); // Start in loading state
  const [isAuthenticated, setIsAuthenticated] = useState(false); // New state for auth status

  useEffect(() => {
    const token = localStorage.getItem('authToken');
    if (!token) {
      router.push('/login');
    } else {
      setIsAuthenticated(true); // User has a token
    }
    setIsLoading(false); // Finished checking auth
  }, [router]);

  if (isLoading) {
    return null; // Or a loading spinner, same on server and initial client
  }

  // If not authenticated and loading is finished, router.push already called.
  // Rendering null here will prevent flashing of content before redirect completes.
  if (!isAuthenticated) {
    return null;
  }

  // Only render content if authenticated and not loading
  // This is the original content of the Home page
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
