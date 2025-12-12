import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import api from '../api/axios';
import './Transactions.css';

const Transactions = () => {
    const { accountId } = useParams();
    const [transfers, setTransfers] = useState([]);
    const [entries, setEntries] = useState([]);
    const [loading, setLoading] = useState(true);
    const navigate = useNavigate();

    useEffect(() => {
        const fetchData = async () => {
            try {
                const [transfersRes, entriesRes] = await Promise.all([
                    api.get(`/transfers?account_id=${accountId}&page_id=1&page_size=20`),
                    api.get(`/accounts/${accountId}/entries?page_id=1&page_size=20`)
                ]);
                setTransfers(transfersRes.data || []);
                setEntries(entriesRes.data || []);
            } catch (error) {
                console.error("Failed to fetch history", error);
            } finally {
                setLoading(false);
            }
        };
        fetchData();
    }, [accountId]);

    if (loading) return <div className="loading">Loading history...</div>;

    return (
        <div className="transactions-container">
            <div className="header-actions">
                <button className="btn-secondary" onClick={() => navigate('/dashboard')}>&larr; Back to Dashboard</button>
                <h2>History for Account #{accountId}</h2>
            </div>

            <div className="history-grid">
                <div className="card history-card">
                    <h3>Recent Transfers</h3>
                    <div className="table-responsive">
                        <table>
                            <thead>
                                <tr>
                                    <th>ID</th>
                                    <th>Amount</th>
                                    <th>Type</th>
                                    <th>Date</th>
                                </tr>
                            </thead>
                            <tbody>
                                {transfers.length === 0 ? (
                                    <tr><td colSpan="4" className="empty-state">No transfers found</td></tr>
                                ) : (
                                    transfers.map(t => (
                                        <tr key={t.id}>
                                            <td>#{t.id}</td>
                                            <td className={t.to_account_id === parseInt(accountId) ? 'text-success' : 'text-danger'}>
                                                {t.to_account_id === parseInt(accountId) ? '+' : '-'} {t.amount}
                                            </td>
                                            <td>
                                                {t.to_account_id === parseInt(accountId) ? `From #${t.from_account_id}` : `To #${t.to_account_id}`}
                                            </td>
                                            <td>{new Date(t.created_at).toLocaleDateString()}</td>
                                        </tr>
                                    ))
                                )}
                            </tbody>
                        </table>
                    </div>
                </div>

                <div className="card history-card">
                    <h3>Balance History (Audit)</h3>
                    <div className="table-responsive">
                        <table>
                            <thead>
                                <tr>
                                    <th>ID</th>
                                    <th>Amount Change</th>
                                    <th>Date</th>
                                </tr>
                            </thead>
                            <tbody>
                                {entries.length === 0 ? (
                                    <tr><td colSpan="3" className="empty-state">No entries found</td></tr>
                                ) : (
                                    entries.map(e => (
                                        <tr key={e.id}>
                                            <td>#{e.id}</td>
                                            <td className={e.amount > 0 ? 'text-success' : 'text-danger'}>
                                                {e.amount > 0 ? '+' : ''} {e.amount}
                                            </td>
                                            <td>{new Date(e.created_at).toLocaleDateString()}</td>
                                        </tr>
                                    ))
                                )}
                            </tbody>
                        </table>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default Transactions;
