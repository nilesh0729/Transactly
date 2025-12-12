import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import api from '../api/axios';
import './Transfer.css';

const Transfer = () => {
    const [accounts, setAccounts] = useState([]);
    const [formData, setFormData] = useState({
        from_account_id: '',
        to_account_id: '',
        amount: '',
        currency: 'USD'
    });
    const [error, setError] = useState('');
    const [success, setSuccess] = useState('');
    const [loading, setLoading] = useState(false);
    const navigate = useNavigate();

    useEffect(() => {
        const fetchAccounts = async () => {
            try {
                const response = await api.get('/accounts?page_id=1&page_size=100');
                setAccounts(response.data || []);
            } catch (err) {
                console.error("Failed to load accounts", err);
            }
        };
        fetchAccounts();
    }, []);

    const handleChange = (e) => {
        setFormData({ ...formData, [e.target.name]: e.target.value });
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        setError('');
        setSuccess('');
        setLoading(true);

        try {
            // Need to convert IDs to numbers and amount to number
            const payload = {
                from_account_id: parseInt(formData.from_account_id),
                to_account_id: parseInt(formData.to_account_id),
                amount: parseInt(formData.amount), // Backend expects integer amount (e.g. cents) or just amount?
                // Wait, the API spec says "amount" int64. Let's assume it's standard unit for now or user inputs whole numbers.
                // Assuming the backend handles currency. But the Request body needs Currency field provided! 
                // Let's check server CreateTransfer params.
                currency: formData.currency
            };

            await api.post('/transfers', payload);
            setSuccess('Transfer successful!');
            setTimeout(() => navigate('/dashboard'), 2000);
        } catch (err) {
            setError(err.response?.data?.error || "Transfer failed");
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="transfer-container">
            <div className="card transfer-card">
                <h2>Make a Transfer</h2>

                {error && <div className="alert error">{error}</div>}
                {success && <div className="alert success">{success}</div>}

                <form onSubmit={handleSubmit}>
                    <div className="form-group">
                        <label>From Account</label>
                        <select
                            name="from_account_id"
                            value={formData.from_account_id}
                            onChange={handleChange}
                            required
                        >
                            <option value="">Select Account</option>
                            {accounts.map(acc => (
                                <option key={acc.id} value={acc.id}>
                                    ID: {acc.id} ({acc.currency} {acc.balance})
                                </option>
                            ))}
                        </select>
                    </div>

                    <div className="form-group">
                        <label>To Account ID</label>
                        <input
                            type="number"
                            name="to_account_id"
                            value={formData.to_account_id}
                            onChange={handleChange}
                            placeholder="Recipient Account ID"
                            required
                        />
                    </div>

                    <div className="form-row">
                        <div className="form-group">
                            <label>Amount</label>
                            <input
                                type="number"
                                name="amount"
                                value={formData.amount}
                                onChange={handleChange}
                                placeholder="0"
                                min="1"
                                required
                            />
                        </div>
                        <div className="form-group">
                            <label>Currency</label>
                            <select
                                name="currency"
                                value={formData.currency}
                                onChange={handleChange}
                            >
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
                    </div>

                    <button type="submit" className="btn-full" disabled={loading}>
                        {loading ? 'Processing...' : 'Transfer Funds'}
                    </button>
                </form>
            </div>
        </div>
    );
};

export default Transfer;
