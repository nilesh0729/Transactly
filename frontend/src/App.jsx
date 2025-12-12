import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { AuthProvider, useAuth } from './context/AuthContext';
import Navbar from './components/Navbar';
import Login from './pages/Login';
import Register from './pages/Register';
import Dashboard from './pages/Dashboard';
import Transfer from './pages/Transfer';
import Transactions from './pages/Transactions';
import './index.css';

// Protected Route Component
const ProtectedRoute = ({ children }) => {
    const { user, loading } = useAuth();

    if (loading) return <div>Loading...</div>;
    if (!user) return <Navigate to="/login" />;

    return children;
};

function App() {
    return (
        <AuthProvider>
            <Router>
                <div className="app-main">
                    <Navbar />
                    <div className="app-container">
                        <Routes>
                            {/* Public Routes */}
                            <Route path="/login" element={<Login />} />
                            <Route path="/register" element={<Register />} />

                            {/* Protected Routes */}
                            <Route path="/dashboard" element={
                                <ProtectedRoute>
                                    <Dashboard />
                                </ProtectedRoute>
                            } />
                            <Route path="/transfer" element={
                                <ProtectedRoute>
                                    <Transfer />
                                </ProtectedRoute>
                            } />
                            <Route path="/transactions/:accountId" element={
                                <ProtectedRoute>
                                    <Transactions />
                                </ProtectedRoute>
                            } />

                            {/* Default Redirect */}
                            <Route path="/" element={<Navigate to="/dashboard" />} />
                        </Routes>
                    </div>
                </div>
            </Router>
        </AuthProvider>
    );
}

export default App;
