import React from 'react';
import TrainingForm from "@/components/training-form"
import Link from 'next/link';


const Recommendations = () => {
  return (
    <main className="min-h-screen py-12 px-4">
      <div className="container mx-auto">
        <div style={{ margin: '1rem' }}>
          <Link href="/" passHref>
            <button style={{ padding: '0.5rem 1rem', fontSize: '1rem' }}>ホームへ</button>
          </Link>
        </div>
        <h1 className="text-3xl font-bold text-center mb-8">パーソナライズドトレーニング提案</h1>
        <TrainingForm />
      </div>
    </main>
  )
};

export default Recommendations;

