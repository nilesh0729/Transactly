import { useState, useEffect } from 'react';
import api from '../api/axios';
import './Dashboard.css';

const Dashboard = () => {
    const [accounts, setAccounts] = useState([]);
    const [loading, setLoading] = useState(true);
    const [showCreateModal, setShowCreateModal] = useState(false);
    const [currency, setCurrency] = useState('USD');

    const fetchAccounts = async () => {
        try {
            const response = await api.get('/accounts?page_id=1&page_size=10');
            setAccounts(response.data || []);
        } catch (error) {
            console.error("Failed to fetch accounts", error);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchAccounts();
    }, []);

    const handleCreateAccount = async () => {
        try {
            await api.post('/accounts', { currency });
            setShowCreateModal(false);
            fetchAccounts();
        } catch (error) {
            alert("Failed to create account: " + (error.response?.data?.error || error.message));
        }
    };

    if (loading) return <div className="loading">Loading accounts...</div>;

    return (
        <div className="dashboard-container">
            <div className="dashboard-header">
                <h2>Your Accounts</h2>
                <button onClick={() => setShowCreateModal(true)}>+ New Account</button>
            </div>

            <div className="accounts-grid">
                {accounts.length === 0 ? (
                    <p className="no-accounts">No accounts found. Create one to get started!</p>
                ) : (
                    accounts.map(account => (
                        <div key={account.id} className="card account-card">
                            <div className="account-header">
                                <span className="currency-badge">{account.currency}</span>
                                <span className="account-id">#{account.id}</span>
                            </div>
                            <div className="account-balance">
                                <h3>{new Intl.NumberFormat('en-US', { style: 'currency', currency: account.currency }).format(account.balance)}</h3>
                                <p>Available Balance</p>
                            </div>
                            <div className="account-footer">
                                <small>Created: {new Date(account.created_at).toLocaleDateString()}</small>
                                <button className="btn-link" onClick={() => window.location.href = `/transactions/${account.id}`}>View History</button>
                            </div>
                        </div>
                    ))
                )}
            </div>

            {showCreateModal && (
                <div className="modal-overlay">
                    <div className="modal">
                        <h3>Create New Account</h3>
                        <div className="form-group">
                            <label>Currency</label>
                            <select value={currency} onChange={(e) => setCurrency(e.target.value)}>
                                <option value="USD">USD</option>
                                <option value="EUR">EUR</option>
                                <option value="CAD">CAD</option>
                                <option value="INR">INR</option>
                                <option value="YEN">YEN</option>
                                <option value="BDT">BDT</option>
                                <option value="BRL">BRL</option>
                                <option value="FJD">FJD</option>
                            </select>
                        </div>
                        <div className="modal-actions">
                            <button className="btn-secondary" onClick={() => setShowCreateModal(false)}>Cancel</button>
                            <button onClick={handleCreateAccount}>Create</button>
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
};

export default Dashboard;
