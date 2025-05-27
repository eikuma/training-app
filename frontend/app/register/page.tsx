"use client";

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import axios from 'axios'; // Using axios as seen in other parts of the project

// Define the API base URL. In a real app, this would likely come from an environment variable.
const API_URL = "http://localhost:8080";

export default function RegisterPage() {
    const [username, setUsername] = useState('');
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [error, setError] = useState('');
    const [successMessage, setSuccessMessage] = useState('');
    const router = useRouter();

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setError('');
        setSuccessMessage('');

        if (!username || !email || !password) {
            setError('All fields are required.');
            return;
        }

        // Basic email validation (can be more sophisticated)
        if (!/\S+@\S+\.\S+/.test(email)) {
            setError('Please enter a valid email address.');
            return;
        }

        // Basic password length validation
        if (password.length < 8) {
            setError('Password must be at least 8 characters long.');
            return;
        }

        try {
            const response = await axios.post(`${API_URL}/auth/register`, {
                username,
                email,
                password,
            });

            // Assuming backend returns a message on success
            setSuccessMessage(response.data.message || 'Registration successful! Redirecting to login...');
            setTimeout(() => {
                router.push('/login');
            }, 2000);

        } catch (err: any) {
            if (axios.isAxiosError(err) && err.response) {
                // Backend error (e.g., 400, 409, 500)
                setError(err.response.data.error || 'Registration failed. Please try again.');
            } else {
                // Network or other unexpected error
                setError('An unexpected error occurred. Please try again.');
                console.error(err);
            }
        }
    };

    return (
        <div className="container mx-auto p-4 max-w-md">
            <h1 className="text-3xl font-bold mb-6 text-center">Register</h1>
            <form onSubmit={handleSubmit} className="space-y-6 bg-white shadow-md rounded px-8 pt-6 pb-8 mb-4">
                <div>
                    <label htmlFor="username" className="block text-gray-700 text-sm font-bold mb-2">
                        Username
                    </label>
                    <input 
                        type="text" 
                        id="username" 
                        value={username} 
                        onChange={(e) => setUsername(e.target.value)} 
                        required 
                        className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
                        placeholder="Choose a username"
                    />
                </div>
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
                        minLength={8}
                        className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 mb-3 leading-tight focus:outline-none focus:shadow-outline"
                        placeholder="Create a password (min. 8 characters)"
                    />
                </div>
                {error && <p className="text-red-500 text-xs italic">{error}</p>}
                {successMessage && <p className="text-green-500 text-xs italic">{successMessage}</p>}
                <div className="flex items-center justify-between">
                    <button 
                        type="submit" 
                        className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline w-full"
                    >
                        Register
                    </button>
                </div>
                 <p className="text-center text-gray-500 text-xs">
                    Already have an account? <a href="/login" className="text-blue-500 hover:text-blue-700">Login here</a>.
                </p>
            </form>
        </div>
    );
}
