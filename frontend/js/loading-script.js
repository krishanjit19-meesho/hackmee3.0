class LoadingPage {
    constructor() {
        this.userData = null;
        this.init();
    }

    init() {
        this.loadUserData();
        this.setupEventListeners();
        this.callCatalogAPI();
    }

    loadUserData() {
        // Get user data from localStorage
        const userDataString = localStorage.getItem('userData');
        if (userDataString) {
            try {
                this.userData = JSON.parse(userDataString);
            } catch (error) {
                console.error('Error parsing user data:', error);
                this.showError('Invalid user data. Please login again.');
            }
        } else {
            this.showError('No user data found. Please login first.');
        }
    }

    setupEventListeners() {
        // Back button
        const backBtn = document.getElementById('backBtn');
        if (backBtn) {
            backBtn.addEventListener('click', () => {
                this.goBack();
            });
        }

        // Add touch feedback for mobile
        if (backBtn) {
            backBtn.addEventListener('touchstart', () => {
                backBtn.style.transform = 'scale(0.95)';
            });
            
            backBtn.addEventListener('touchend', () => {
                backBtn.style.transform = 'scale(1)';
            });
        }
    }



    async callCatalogAPI() {
        if (!this.userData || !this.userData.userId) {
            this.showError('No user data available');
            return;
        }

        try {
            // Call the catalog API with user ID
            const response = await fetch(`http://localhost:8080/api/v1/catalog/?user_id=${this.userData.userId}`, {
                method: 'GET',
                headers: {
                    'Content-Type': 'application/json',
                }
            });

            const result = await response.json();

            if (response.ok && result.success) {
                // Store the catalog data in localStorage for the product listing page
                localStorage.setItem('catalogData', JSON.stringify(result));
                
                // Redirect to product listing page immediately
                window.location.href = 'product-listing.html';
            } else {
                console.error('Failed to load catalog data:', result.error);
                this.showError('Failed to load catalog data. Please try again.');
            }
        } catch (error) {
            console.error('Error calling catalog API:', error);
            this.showError('Network error. Please check your connection and try again.');
        }
    }





    showError(message) {
        // Update the loading title and subtitle
        const title = document.querySelector('.loading-title');
        const subtitle = document.querySelector('.loading-subtitle');
        
        if (title) title.textContent = 'Error Loading Catalog';
        if (subtitle) subtitle.textContent = message;
        
        // Change spinner to error icon
        const spinner = document.querySelector('.spinner');
        if (spinner) {
            spinner.style.border = '3px solid #f44336';
            spinner.style.borderTop = '3px solid #f44336';
            spinner.style.animation = 'none';
            spinner.innerHTML = '✗';
            spinner.style.fontSize = '20px';
            spinner.style.color = '#f44336';
            spinner.style.display = 'flex';
            spinner.style.alignItems = 'center';
            spinner.style.justifyContent = 'center';
        }

        // Show retry button
        this.showRetryButton();
    }

    showRetryButton() {
        const loadingCard = document.querySelector('.loading-card');
        if (loadingCard) {
            const retryButton = document.createElement('button');
            retryButton.textContent = 'Retry';
            retryButton.className = 'retry-btn';
            retryButton.style.cssText = `
                background: linear-gradient(135deg, #9f2089, #667eea);
                color: white;
                border: none;
                padding: 12px 24px;
                border-radius: 8px;
                font-size: 16px;
                font-weight: 600;
                cursor: pointer;
                margin-top: 20px;
                transition: all 0.3s ease;
            `;
            
            retryButton.addEventListener('click', () => {
                window.location.reload();
            });
            
            loadingCard.appendChild(retryButton);
        }
    }

    goBack() {
        // Go back to home page
        window.location.href = 'home.html';
    }
}

// Initialize the loading page when DOM is loaded
document.addEventListener('DOMContentLoaded', () => {
    new LoadingPage();
}); 