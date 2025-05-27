"use client";

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import axios from 'axios';

// Define the API base URL, consistent with the registration page.
const API_URL = "http://localhost:8080";

export default function LoginPage() {
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [error, setError] = useState('');
    const router = useRouter();

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setError('');

        if (!email || !password) {
            setError('Email and password are required.');
            return;
        }

        try {
            const response = await axios.post(`${API_URL}/auth/login`, {
                email,
                password,
            });

            if (response.data && response.data.token) {
                localStorage.setItem('authToken', response.data.token); // Store the token
                // Optionally, add a success message or directly redirect
                // setSuccessMessage('Login successful! Redirecting...');
                // setTimeout(() => {
                router.push('/workouts'); // Redirect to a protected route
                // }, 1000); 
            } else {
                // This case might occur if the backend responds 200 OK but without a token
                setError('Login failed. No token received.');
            }
        } catch (err: any) {
            if (axios.isAxiosError(err) && err.response) {
                // Error from backend (e.g., 400, 401)
                setError(err.response.data.error || 'Login failed. Please check your credentials.');
            } else {
                // Network or other unexpected error
                setError('An unexpected error occurred. Please try again.');
                console.error('Login error:', err);
            }
        }
    };

    return (
        <div className="container mx-auto p-4 max-w-md">
            <h1 className="text-3xl font-bold mb-6 text-center">Login</h1>
            <form onSubmit={handleSubmit} className="space-y-6 bg-white shadow-md rounded px-8 pt-6 pb-8 mb-4">
                <div>
                    <label htmlFor="email" className="block text-gray-700 text-sm font-bold mb-2">
                        Email
                    </label>
                    <input 
                        type="email" 
                        id="email" 
                        value={email} 
                        onChange={(e) => setEmail(e.target.value)} 
                        required 
                        className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
                        placeholder="Enter your email"
                    />
                </div>
                <div>
                    <label htmlFor="password" className="block text-gray-700 text-sm font-bold mb-2">
                        Password
                    </label>
                    <input 
                        type="password" 
                        id="password" 
                        value={password} 
                        onChange={(e) => setPassword(e.target.value)} 
                        required
                        className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 mb-3 leading-tight focus:outline-none focus:shadow-outline"
                        placeholder="Enter your password"
                    />
                </div>
                {error && <p className="text-red-500 text-xs italic text-center">{error}</p>}
                <div className="flex items-center justify-between">
                    <button 
                        type="submit" 
                        className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline w-full"
                    >
                        Sign In
                    </button>
                </div>
                <p className="text-center text-gray-500 text-xs">
                    Need an account? <a href="/register" className="text-blue-500 hover:text-blue-700 font-bold">Register here</a>.
                </p>
            </form>
        </div>
    );
}
