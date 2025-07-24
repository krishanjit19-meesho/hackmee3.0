// Simple Meesho Home Page
class MeeshoHome {
    constructor() {
        this.userData = null;
        this.products = [];
        this.init();
    }

    init() {
        this.loadUserData();
        this.setupEventListeners();
        this.loadProducts();
    }

    loadUserData() {
        const userData = localStorage.getItem('userData');
        if (!userData) {
            // Redirect to login if no user data
            window.location.href = 'index.html';
            return;
        }
        
        this.userData = JSON.parse(userData);
        this.displayUserInfo();
    }

    displayUserInfo() {
        const userNameElement = document.getElementById('userName');
        if (userNameElement && this.userData) {
            // Use the name from API response if available, otherwise use phone number
            if (this.userData.name) {
                userNameElement.textContent = this.userData.name;
            } else {
                const phoneNumber = this.userData.phone.replace('+91', '');
                userNameElement.textContent = `User ${phoneNumber.slice(-4)}`;
            }
        }
    }

    setupEventListeners() {
        // Search functionality
        const searchInput = document.querySelector('.search-bar input');
        if (searchInput) {
            searchInput.addEventListener('keypress', (e) => {
                if (e.key === 'Enter') this.handleSearch();
            });
            
            // Add touch-friendly search
            searchInput.addEventListener('focus', () => {
                searchInput.style.fontSize = '16px'; // Prevent zoom on iOS
            });
        }

        // Category clicks with touch support
        document.querySelectorAll('.category-card').forEach(card => {
            card.addEventListener('click', (e) => {
                const categoryText = e.currentTarget.querySelector('span').textContent;
                this.filterByCategory(categoryText);
            });
            
            // Add touch feedback
            card.addEventListener('touchstart', () => {
                card.style.transform = 'scale(0.95)';
            });
            
            card.addEventListener('touchend', () => {
                card.style.transform = 'scale(1)';
            });
        });

        // Filter buttons with touch support
        document.querySelectorAll('.filter-btn').forEach(btn => {
            btn.addEventListener('click', (e) => {
                this.handleFilter(e.currentTarget.textContent.trim());
            });
            
            // Add touch feedback
            btn.addEventListener('touchstart', () => {
                btn.style.transform = 'scale(0.95)';
            });
            
            btn.addEventListener('touchend', () => {
                btn.style.transform = 'scale(1)';
            });
        });

        // Wishlist and cart buttons
        const wishlistBtn = document.querySelector('.wishlist-btn');
        const cartBtn = document.querySelector('.cart-btn');
        
        if (wishlistBtn) {
            wishlistBtn.addEventListener('click', () => {
                alert('Wishlist feature coming soon!');
            });
        }
        
        if (cartBtn) {
            cartBtn.addEventListener('click', () => {
                alert('Cart feature coming soon!');
            });
        }

        // Bottom navigation with touch support
        document.querySelectorAll('.nav-item').forEach(item => {
            item.addEventListener('click', (e) => {
                e.preventDefault();
                this.handleNavigation(e.currentTarget);
            });
            
            // Add touch feedback
            item.addEventListener('touchstart', () => {
                item.style.transform = 'scale(0.95)';
            });
            
            item.addEventListener('touchend', () => {
                item.style.transform = 'scale(1)';
            });
        });
        
        // Banner click functionality
        const bannerImage = document.getElementById('bannerImage');
        if (bannerImage) {
            bannerImage.addEventListener('click', () => {
                this.handleBannerClick();
            });
            
            // Add touch feedback for banner
            bannerImage.addEventListener('touchstart', () => {
                bannerImage.style.transform = 'scale(0.98)';
            });
            
            bannerImage.addEventListener('touchend', () => {
                bannerImage.style.transform = 'scale(1)';
            });

            // Add loading indicator
            bannerImage.style.cursor = 'pointer';
            bannerImage.title = 'Click to view catalog products';
        }
        
        // Prevent zoom on double tap
        document.addEventListener('touchend', (e) => {
            if (e.target.tagName === 'INPUT' || e.target.tagName === 'BUTTON') {
                e.preventDefault();
            }
        });
    }

    async loadProducts() {
        if (!this.userData || !this.userData.userId) {
            console.error('No user data available');
            return;
        }

        try {
            // Call the homescreen API
            const response = await fetch(`http://localhost:8080/api/v1/home/?user_id=${this.userData.userId}`, {
                method: 'GET',
                headers: {
                    'Content-Type': 'application/json',
                }
            });

            const result = await response.json();

            if (response.ok && result.success) {
                // Store the API response data
                this.homescreenData = result;
                
                // Display categories from the API
                this.displayCategories(result.categories || []);
                
                // Display products from the API
                this.displayProductsFromAPI(result.products || []);
                
                // Update user info if available from API
                if (result.user_info && result.user_info.name) {
                    this.userData.name = result.user_info.name;
                    this.displayUserInfo();
                }
            } else {
                console.error('Failed to load homescreen data:', result.error);
                // Fallback to mock data
                this.loadMockProducts();
            }
        } catch (error) {
            console.error('Error loading homescreen data:', error);
            // Fallback to mock data
            this.loadMockProducts();
        }
    }

    loadMockProducts() {
        // Mock product data as fallback
        const mockProducts = [
            {
                id: 1,
                title: "Useful Manual Choppers & Chippers",
                currentPrice: 107,
                originalPrice: 124,
                discount: "14% off",
                rating: 4.1,
                reviews: 36275,
                image: "data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='140' height='140' viewBox='0 0 140 140'%3E%3Crect width='140' height='140' fill='%23e8f8f5'/%3E%3Ctext x='70' y='75' text-anchor='middle' font-size='12' fill='%2327ae60'%3EKitchen%3C/text%3E%3C/svg%3E",
                deliveryTag: "Free Delivery"
            },
            {
                id: 2,
                title: "Denver Black Code,Caliber,Hutle...",
                currentPrice: 247,
                originalPrice: 330,
                discount: "₹73 OFF",
                rating: 4.8,
                reviews: 17000,
                image: "data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='140' height='140' viewBox='0 0 140 140'%3E%3Crect width='140' height='140' fill='%23f0f0f0'/%3E%3Ctext x='70' y='75' text-anchor='middle' font-size='12' fill='%23333'%3EDeodorant%3C/text%3E%3C/svg%3E",
                discountTag: "₹73 OFF"
            },
            {
                id: 3,
                title: "Cotton Saree with Blouse",
                currentPrice: 599,
                originalPrice: 1299,
                discount: "54% OFF",
                rating: 4.2,
                reviews: 1234,
                image: "data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='140' height='140' viewBox='0 0 140 140'%3E%3Crect width='140' height='140' fill='%23ffeef2'/%3E%3Ctext x='70' y='75' text-anchor='middle' font-size='12' fill='%23e91e63'%3ESaree%3C/text%3E%3C/svg%3E",
                deliveryTag: "Free Delivery"
            },
            {
                id: 4,
                title: "Designer Kurti for Women",
                currentPrice: 399,
                originalPrice: 899,
                discount: "56% OFF",
                rating: 4.0,
                reviews: 856,
                image: "data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='140' height='140' viewBox='0 0 140 140'%3E%3Crect width='140' height='140' fill='%23f8e8ff'/%3E%3Ctext x='70' y='75' text-anchor='middle' font-size='12' fill='%23e91e63'%3EKurti%3C/text%3E%3C/svg%3E",
                discountTag: "56% OFF"
            },
            {
                id: 5,
                title: "Men's Cotton Casual Shirt",
                currentPrice: 299,
                originalPrice: 699,
                discount: "57% OFF",
                rating: 4.1,
                reviews: 642,
                image: "data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='140' height='140' viewBox='0 0 140 140'%3E%3Crect width='140' height='140' fill='%23e8f2ff'/%3E%3Ctext x='70' y='75' text-anchor='middle' font-size='12' fill='%23333'%3EShirt%3C/text%3E%3C/svg%3E",
                deliveryTag: "Free Delivery"
            },
            {
                id: 6,
                title: "Premium Beauty Face Cream",
                currentPrice: 149,
                originalPrice: 299,
                discount: "50% OFF",
                rating: 4.0,
                reviews: 789,
                image: "data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='140' height='140' viewBox='0 0 140 140'%3E%3Crect width='140' height='140' fill='%23fff0e8'/%3E%3Ctext x='70' y='75' text-anchor='middle' font-size='12' fill='%23e91e63'%3EBeauty%3C/text%3E%3C/svg%3E",
                discountTag: "50% OFF"
            }
        ];

        this.products = mockProducts;
        this.displayProducts(mockProducts);
    }

    displayProducts(products) {
        const container = document.getElementById('productsGrid');
        if (!container) return;

        container.innerHTML = products.map(product => `
            <div class="product-card" data-product-id="${product.id}">
                ${product.deliveryTag ? `<div class="delivery-tag">${product.deliveryTag}</div>` : ''}
                ${product.discountTag ? `<div class="discount-tag">${product.discountTag}</div>` : ''}
                <div class="wishlist-heart">♡</div>
                <img src="${product.image}" alt="${product.title}" class="product-image">
                <div class="product-title">${product.title}</div>
                <div class="product-price">
                    <span class="current-price">₹${product.currentPrice}</span>
                    <span class="original-price">₹${product.originalPrice}</span>
                    <span class="discount">${product.discount}</span>
                </div>
                <div class="product-rating">
                    <span class="rating-stars">★★★★☆</span>
                    <span>${product.rating} (${product.reviews.toLocaleString()})</span>
                </div>
            </div>
        `).join('');

        // Add click listeners to product cards
        container.querySelectorAll('.product-card').forEach(card => {
            card.addEventListener('click', (e) => {
                if (!e.target.classList.contains('wishlist-heart')) {
                    const productId = e.currentTarget.dataset.productId;
                    this.showProductDetails(productId);
                }
            });
        });

        // Add wishlist functionality
        container.querySelectorAll('.wishlist-heart').forEach(heart => {
            heart.addEventListener('click', (e) => {
                e.stopPropagation();
                if (heart.textContent === '♡') {
                    heart.textContent = '♥';
                    heart.style.color = '#e91e63';
                } else {
                    heart.textContent = '♡';
                    heart.style.color = '#666';
                }
            });
        });
    }

    displayCategories(categories) {
        // Update the categories section with data from API
        const categoriesContainer = document.querySelector('.categories-section .categories-grid');
        if (!categoriesContainer || !categories.length) return;

        // Clear existing categories and add new ones from API
        categoriesContainer.innerHTML = categories.map(category => `
            <div class="category-card" data-category-id="${category.id}">
                <div class="category-icon">
                    <img src="${category.image}" alt="${category.title}" style="width: 40px; height: 40px; object-fit: cover; border-radius: 8px;">
                </div>
                <span>${category.title}</span>
            </div>
        `).join('');

        // Add click listeners to new category cards
        categoriesContainer.querySelectorAll('.category-card').forEach(card => {
            card.addEventListener('click', (e) => {
                const categoryId = e.currentTarget.dataset.categoryId;
                const categoryTitle = e.currentTarget.querySelector('span').textContent;
                this.filterByCategory(categoryTitle);
            });
        });
    }

    displayProductsFromAPI(products) {
        const container = document.getElementById('productsGrid');
        if (!container) return;

        // Convert API products to display format
        const displayProducts = products.map(product => ({
            id: product.id,
            title: product.title,
            image: product.image,
            currentPrice: Math.floor(Math.random() * 500) + 100, // Mock price for demo
            originalPrice: Math.floor(Math.random() * 1000) + 500,
            discount: `${Math.floor(Math.random() * 50) + 10}% OFF`,
            rating: (Math.random() * 2 + 3).toFixed(1), // Random rating between 3-5
            reviews: Math.floor(Math.random() * 50000) + 1000,
            deliveryTag: Math.random() > 0.5 ? "Free Delivery" : null,
            discountTag: Math.random() > 0.5 ? `${Math.floor(Math.random() * 100) + 50} OFF` : null
        }));

        this.products = displayProducts;
        this.displayProducts(displayProducts);
    }

    handleSearch() {
        const searchInput = document.querySelector('.search-bar input');
        const searchTerm = searchInput.value.trim();
        
        if (searchTerm) {
            const filteredProducts = this.products.filter(product =>
                product.title.toLowerCase().includes(searchTerm.toLowerCase())
            );
            this.displayProducts(filteredProducts);
        } else {
            this.displayProducts(this.products);
        }
    }

    filterByCategory(category) {
        // Simple category filtering
        if (category === 'All Categories') {
            this.displayProducts(this.products);
        } else {
            // For demo, just show a message
            alert(`Filtering by ${category} - Feature coming soon!`);
        }
    }

    handleFilter(filterType) {
        // Handle filter buttons
        alert(`${filterType} filter - Feature coming soon!`);
    }

    showProductDetails(productId) {
        const product = this.products.find(p => p.id == productId);
        if (product) {
            alert(`Product: ${product.title}\nPrice: ₹${product.currentPrice}\nRating: ${product.rating}\n\nProduct details page coming soon!`);
        }
    }

    async handleBannerClick() {
        if (!this.userData || !this.userData.userId) {
            console.error('No user data available for banner click');
            alert('Please login first');
            return;
        }

        try {
            // Show loading state
            const bannerImage = document.getElementById('bannerImage');
            if (bannerImage) {
                bannerImage.style.opacity = '0.7';
                bannerImage.style.cursor = 'wait';
            }

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
                
                // Redirect to product listing page
                window.location.href = 'product-listing.html';
            } else {
                console.error('Failed to load catalog data:', result.error);
                alert('Failed to load catalog data. Please try again.');
            }
        } catch (error) {
            console.error('Error calling catalog API:', error);
            alert('Network error. Please check your connection and try again.');
        } finally {
            // Reset banner state
            const bannerImage = document.getElementById('bannerImage');
            if (bannerImage) {
                bannerImage.style.opacity = '1';
                bannerImage.style.cursor = 'pointer';
            }
        }
    }

    handleNavigation(navItem) {
        // Update active nav item
        document.querySelectorAll('.nav-item').forEach(item => {
            item.classList.remove('active');
        });
        navItem.classList.add('active');

        // Handle navigation
        const navText = navItem.querySelector('.nav-text').textContent;
        
        switch(navText) {
            case 'Home':
                this.loadProducts();
                break;
            case 'Categories':
                alert('Categories page - Coming soon!');
                break;
            case 'Mall':
                alert('Mall page - Coming soon!');
                break;
            case 'Live Shop':
                alert('Live Shop - Coming soon!');
                break;
            case 'My Orders':
                alert('My Orders - Coming soon!');
                break;
        }
    }
}

// Initialize the home page when DOM is loaded
document.addEventListener('DOMContentLoaded', function() {
    window.meeshoHome = new MeeshoHome();
});

// Export for use in other files
window.MeeshoHome = MeeshoHome; 