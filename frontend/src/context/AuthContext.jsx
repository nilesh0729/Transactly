import { createContext, useState, useEffect, useContext } from 'react';
import api from '../api/axios';

const AuthContext = createContext();

export const AuthProvider = ({ children }) => {
    const [user, setUser] = useState(null);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        // Check if token exists
        const token = localStorage.getItem('access_token');
        const storedUser = localStorage.getItem('user');
        if (token && storedUser) {
            setUser(JSON.parse(storedUser));
        }
        setLoading(false);
    }, []);

    const login = async (username, password) => {
        try {
            const response = await api.post('/user/login', { username, password });
            const { access_token, user: userData } = response.data;

            // The API returns access_token and user info
            localStorage.setItem('access_token', access_token);
            localStorage.setItem('user', JSON.stringify(userData));
            setUser(userData);
            return { success: true };
        } catch (error) {
            console.error("Login failed", error);
            return {
                success: false,
                error: error.response?.data?.error || "Login failed"
            };
        }
    };

    const register = async (username, password, email, fullName) => {
        try {
            // POST /users based on API Reference, but verification showed POST /user
            // Let's check server.go again to be sure: router.POST("/user", server.CreateUser)
            await api.post('/user', {
                username,
                password,
                full_name: fullName,
                email
            });
            // Auto login after register or just redirect? Let's just return success
            return { success: true };
        } catch (error) {
            console.error("Registration failed", error);
            return {
                success: false,
                error: error.response?.data?.error || "Registration failed"
            };
        }
    };

    const logout = () => {
        localStorage.removeItem('access_token');
        localStorage.removeItem('user');
        setUser(null);
    };

    return (
        <AuthContext.Provider value={{ user, login, register, logout, loading }}>
            {children}
        </AuthContext.Provider>
    );
};

export const useAuth = () => useContext(AuthContext);
