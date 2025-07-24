// DOM Elements
const phoneInput = document.getElementById('phoneNumber');
const continueBtn = document.getElementById('continueBtn');

// Phone number validation
function validatePhoneNumber(phone) {
    // Remove any spaces or special characters
    const cleanPhone = phone.replace(/\D/g, '');
    
    // Check if it's a valid 10-digit Indian mobile number
    const phoneRegex = /^[6-9]\d{9}$/;
    return phoneRegex.test(cleanPhone);
}

// Enable/disable continue button based on phone validation
function updateContinueButton() {
    const phoneValue = phoneInput.value.trim();
    const isValid = validatePhoneNumber(phoneValue);
    
    continueBtn.disabled = !isValid;
    
    if (isValid) {
        continueBtn.style.background = '#e91e63';
        continueBtn.style.cursor = 'pointer';
    } else {
        continueBtn.style.background = '#ccc';
        continueBtn.style.cursor = 'not-allowed';
    }
}

// Format phone number as user types (optional)
function formatPhoneNumber(value) {
    // Remove all non-digits
    const digits = value.replace(/\D/g, '');
    
    // Limit to 10 digits
    return digits.substring(0, 10);
}

// Event listeners
phoneInput.addEventListener('input', function(e) {
    // Format the input value
    const formattedValue = formatPhoneNumber(e.target.value);
    e.target.value = formattedValue;
    
    // Update button state
    updateContinueButton();
});

// Add touch-friendly improvements
phoneInput.addEventListener('focus', function() {
    this.style.fontSize = '16px'; // Prevent zoom on iOS
});

phoneInput.addEventListener('blur', function() {
    this.style.fontSize = '16px';
});

// Handle continue button click
continueBtn.addEventListener('click', async function(e) {
    e.preventDefault();
    
    const phoneNumber = phoneInput.value.trim();
    
    if (!validatePhoneNumber(phoneNumber)) {
        alert('Please enter a valid 10-digit mobile number');
        return;
    }
    
    // Show loading state
    continueBtn.innerHTML = 'Processing...';
    continueBtn.disabled = true;
    
    try {
        // Call the login API
        const response = await fetch('http://localhost:8080/api/v1/auth/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                phone_number: phoneNumber
            })
        });
        
        const result = await response.json();
        
        if (response.ok && result.success) {
            // Store user data from API response
            const userData = {
                phone: '+91' + phoneNumber,
                loginTime: new Date().toISOString(),
                userId: result.data.user_id,
                name: result.data.name
            };
            
            // Store in localStorage
            localStorage.setItem('userData', JSON.stringify(userData));
            
            console.log('User logged in:', userData);
            
            // Redirect to home page
            window.location.href = 'home.html';
        } else {
            // Show error message
            alert(result.error || 'Login failed. Please try again.');
            continueBtn.innerHTML = 'Continue';
            continueBtn.disabled = false;
        }
    } catch (error) {
        console.error('Login error:', error);
        alert('Network error. Please check your connection and try again.');
        continueBtn.innerHTML = 'Continue';
        continueBtn.disabled = false;
    }
});

// Initialize button state
document.addEventListener('DOMContentLoaded', function() {
    updateContinueButton();
    
    // Search functionality (basic)
    const searchInput = document.querySelector('.search-bar input');
    const searchBtn = document.querySelector('.search-btn');
    
    searchBtn.addEventListener('click', function() {
        const searchTerm = searchInput.value.trim();
        if (searchTerm) {
            alert(`Searching for: ${searchTerm}`);
            // In real app, this would trigger search functionality
        }
    });
    
    searchInput.addEventListener('keypress', function(e) {
        if (e.key === 'Enter') {
            searchBtn.click();
        }
    });
});

// Additional utility functions for future RTO functionality
function getUserData() {
    const userData = localStorage.getItem('userData');
    return userData ? JSON.parse(userData) : null;
}

function isUserLoggedIn() {
    return getUserData() !== null;
}

// Mock function to simulate RTO data fetching (for future integration)
function mockFetchRTOOrders(userId) {
    // This will be replaced with actual API calls later
    const mockRTOOrders = [
        {
            orderId: 'ORD_001',
            productName: 'Blue Cotton Saree',
            originalPrice: 1299,
            discountedPrice: 999,
            rtoRisk: 'high',
            locality: 'Bangalore',
            branch: 'Electronics',
            category: 'Fashion'
        },
        {
            orderId: 'ORD_002',
            productName: 'Red Designer Kurti',
            originalPrice: 899,
            discountedPrice: 699,
            rtoRisk: 'medium',
            locality: 'Mumbai',
            branch: 'Fashion',
            category: 'Women Clothing'
        }
    ];
    
    return Promise.resolve(mockRTOOrders);
}

// Export functions for potential module usage
window.MeeshoApp = {
    validatePhoneNumber,
    getUserData,
    isUserLoggedIn,
    mockFetchRTOOrders
}; 