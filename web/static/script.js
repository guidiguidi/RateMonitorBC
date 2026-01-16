const API_BASE_URL = 'http://localhost:8080/api';
const DEMO_CURRENCIES = [
    { id: 1, name: 'Bitcoin', code: 'BTC' },
    { id: 2, name: 'Ethereum', code: 'ETH' },
    { id: 3, name: 'US Dollar', code: 'USD' },
    { id: 4, name: 'Euro', code: 'EUR' },
    { id: 5, name: 'Polish Zloty', code: 'PLN' },
    { id: 6, name: 'Ukrainian Hryvnia', code: 'UAH' }
];

const elements = {
    form: document.getElementById('rateForm'),
    fromCurrency: document.getElementById('fromCurrency'),
    toCurrency: document.getElementById('toCurrency'),
    amount: document.getElementById('amount'),
    loading: document.getElementById('loading'),
    resultsSection: document.getElementById('resultsSection'),
    resultsList: document.getElementById('resultsList'),
    error: document.getElementById('error'),
    successMessage: document.getElementById('successMessage')
};

// Initialize currency dropdowns
function initializeCurrencies() {
    const populateSelect = (selectElement) => {
        DEMO_CURRENCIES.forEach(curr => {
            const option = document.createElement('option');
            option.value = curr.id;
            option.textContent = `${curr.code} - ${curr.name}`;
            selectElement.appendChild(option);
        });
    };

    populateSelect(elements.fromCurrency);
    populateSelect(elements.toCurrency);

    // Set default values
    elements.fromCurrency.value = '1'; // BTC
    elements.toCurrency.value = '3';   // USD
}

// Get selected marks
function getSelectedMarks() {
    const marks = [];
    document.querySelectorAll('.mark-item input[type="checkbox"]:checked').forEach(checkbox => {
        marks.push(checkbox.value);
    });
    return marks;
}

// Show error message
function showError(message) {
    elements.error.textContent = message;
    elements.error.classList.add('show');
    setTimeout(() => elements.error.classList.remove('show'), 5000);
}

// Show success message
function showSuccess(message) {
    elements.successMessage.textContent = message;
    elements.successMessage.classList.add('show');
    setTimeout(() => elements.successMessage.classList.remove('show'), 4000);
}

// Format currency value
function formatCurrency(value) {
    return parseFloat(value).toFixed(8);
}

// Create result card HTML
function createResultCard(rate, fromCode, toCode, amount) {
    const toAmount = parseFloat(rate.to_amount || 0).toFixed(8);
    const marks = rate.marks || [];

    return `
        <div class="result-card">
            <div class="result-header">
                <div class="exchanger-name">${rate.exchanger_id}</div>
                <div class="exchange-rate">${formatCurrency(rate.rate)}</div>
            </div>

            <div class="result-details">
                <div class="detail-item">
                    <div class="detail-label">You Send</div>
                    <div class="detail-value">${amount} ${fromCode}</div>
                </div>
                <div class="detail-item">
                    <div class="detail-label">You Get</div>
                    <div class="detail-value">${toAmount} ${toCode}</div>
                </div>
                <div class="detail-item">
                    <div class="detail-label">Min Amount</div>
                    <div class="detail-value">${formatCurrency(rate.inmin)}</div>
                </div>
                <div class="detail-item">
                    <div class="detail-label">Max Amount</div>
                    <div class="detail-value">${formatCurrency(rate.inmax)}</div>
                </div>
            </div>

            ${marks.length > 0 ? `
                <div class="detail-item">
                    <div class="detail-label">Tags</div>
                    <div class="marks-list">
                        ${marks.map(mark => `<span class="mark-badge">${mark}</span>`).join('')}
                    </div>
                </div>
            ` : ''}
        </div>
    `;
}

// Fetch and display results
        async function searchRates(e) {
            e.preventDefault();

            const fromId = elements.fromCurrency.value;
            const toId = elements.toCurrency.value;
            const amount = elements.amount.value;
            const marks = getSelectedMarks();

            if (!fromId || !toId) {
                showError('Please select both currencies');
                return;
            }

            if (amount <= 0) {
                showError('Amount must be greater than 0');
                return;
            }

            elements.loading.classList.add('show');
            elements.resultsSection.classList.remove('show');
            elements.error.classList.remove('show');
            elements.resultsList.innerHTML = ''; // Clear previous results

            try {
                const fromCurr = DEMO_CURRENCIES.find(c => c.id == fromId);
                const toCurr = DEMO_CURRENCIES.find(c => c.id == toId);

                const response = await fetch(`${API_BASE_URL}/v1/best-rate`, { // Updated endpoint
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ from_code: fromCurr.code, to_code: toCurr.code, amount: parseFloat(amount), marks })
                });

                if (!response.ok) {
                    const errorData = await response.json();
                    throw new Error(errorData.error || `HTTP error! status: ${response.status}`);
                }

                const result = await response.json();
                console.log(result);

                if (result.best_rate && fromCurr && toCurr) {
                    elements.resultsList.innerHTML = createResultCard(result.best_rate, fromCurr.code, toCurr.code, amount);
                    elements.resultsSection.classList.add('show');
                    showSuccess(`Found best rate for ${fromCurr.code} → ${toCurr.code}`);
                } else {
                    showError('No suitable rates found.');
                }

            } catch (err) {
                showError(`Error: ${err.message}`);
            } finally {
                elements.loading.classList.remove('show');
            }
        }
// Event listeners
elements.form.addEventListener('submit', searchRates);

// Initialize on load
initializeCurrencies();
showSuccess('Ready to find the best exchange rates!');
