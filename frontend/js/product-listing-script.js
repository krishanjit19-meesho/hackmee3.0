// Product Listing Page
class ProductListing {
    constructor() {
        this.products = [];
        this.filteredProducts = [];
        this.currentFilters = {
            sort: 'default',
            category: 'all',
            price: 'all'
        };
        this.init();
    }

    init() {
        this.setupEventListeners();
        this.loadProducts();
    }

    setupEventListeners() {
        // Back button
        const backBtn = document.getElementById('backBtn');
        if (backBtn) {
            backBtn.addEventListener('click', () => {
                window.history.back();
            });
        }

        // Filter buttons
        document.querySelectorAll('.filter-btn').forEach(btn => {
            btn.addEventListener('click', (e) => {
                this.handleFilterClick(e.currentTarget);
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
    }

    async loadProducts() {
        this.showLoading(true);
        
        // Simulate API call delay
        await new Promise(resolve => setTimeout(resolve, 1000));
        
        // Mock product data similar to Meesho style
        this.products = [
            {
                id: 's-532528444',
                title: 'Graceful Artificial Plant, Flower and Decorative Item for Home',
                currentPrice: 163,
                originalPrice: 175,
                discount: '7% off',
                specialOffer: '₹147 with 1 Special Offer',
                delivery: 'Free Delivery',
                rating: 3.8,
                reviews: 444,
                image: '../assets/icon.png',
                category: 'Home & Garden',
                isTrusted: false
            },
            {
                id: 's-541804123',
                title: 'New Collections Of Planter Stand',
                currentPrice: 488,
                originalPrice: 541,
                discount: '10% off',
                specialOffer: '₹408 with 2 Special Offers',
                delivery: 'Free Delivery',
                rating: 3.9,
                reviews: 586,
                image: '../assets/icon.png',
                category: 'Home & Garden',
                isTrusted: true
            },
            {
                id: 's-394453979',
                title: 'Unique Artificial Plant, Flower and Decorative Item',
                currentPrice: 539,
                originalPrice: 579,
                discount: '7% off',
                specialOffer: '₹484 with 2 Special Offers',
                delivery: 'Free Delivery',
                rating: 4.3,
                reviews: 6317,
                image: '../assets/icon.png',
                category: 'Home & Garden',
                isTrusted: false
            },
            {
                id: 's-388338864',
                title: 'Attractive Artificial Plant, Flower and Decorative Item',
                currentPrice: 281,
                originalPrice: 301,
                discount: '7% off',
                specialOffer: '₹260 with 2 Special Offers',
                delivery: 'Free Delivery',
                rating: 4.4,
                reviews: 3485,
                image: '../assets/icon.png',
                category: 'Home & Garden',
                isTrusted: false
            },
            {
                id: 's-123456789',
                title: 'Premium Artificial Bonsai Tree for Office Decoration',
                currentPrice: 899,
                originalPrice: 1299,
                discount: '31% off',
                specialOffer: '₹799 with 3 Special Offers',
                delivery: 'Free Delivery',
                rating: 4.6,
                reviews: 892,
                image: '../assets/icon.png',
                category: 'Home & Garden',
                isTrusted: true
            },
            {
                id: 's-987654321',
                title: 'Modern Hanging Plant Pot with Chain',
                currentPrice: 245,
                originalPrice: 345,
                discount: '29% off',
                specialOffer: '₹220 with 1 Special Offer',
                delivery: 'Free Delivery',
                rating: 4.1,
                reviews: 1234,
                image: '../assets/icon.png',
                category: 'Home & Garden',
                isTrusted: false
            }
        ];

        this.filteredProducts = [...this.products];
        this.displayProducts();
        this.showLoading(false);
    }

    displayProducts() {
        const container = document.getElementById('productsGrid');
        if (!container) return;

        container.innerHTML = this.filteredProducts.map(product => `
            <div class="product-card" data-product-id="${product.id}">
                <div class="wishlist-heart">♡</div>
                <img src="${product.image}" alt="${product.title}" class="product-image" onerror="this.src='../assets/icon.png'">
                
                <div class="product-header">
                    <div class="delivery-tag">${product.delivery}</div>
                    <div class="rating-section">
                        <span class="rating-stars">★</span>
                        <span>${product.rating}</span>
                        ${product.isTrusted ? '<span class="trusted-badge">m Trusted</span>' : ''}
                    </div>
                </div>
                
                <div class="product-id">${product.id}</div>
                <div class="product-title">${product.title}</div>
                
                <div class="product-price">
                    <span class="current-price">₹${product.currentPrice}</span>
                    <span class="original-price">₹${product.originalPrice}</span>
                    <span class="discount">${product.discount}</span>
                </div>
                
                <div class="special-offer">${product.specialOffer}</div>
                <div class="delivery-info">${product.delivery}</div>
                
                <div class="product-rating">
                    <span class="rating-stars">★★★★☆</span>
                    <span>(${product.reviews.toLocaleString()})</span>
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

    handleFilterClick(button) {
        const buttonText = button.textContent.trim();
        
        if (buttonText.includes('Sort')) {
            this.showSortOptions();
        } else if (buttonText.includes('Category')) {
            this.showCategoryOptions();
        } else if (buttonText.includes('Price')) {
            this.showPriceOptions();
        } else if (buttonText.includes('Filters')) {
            this.showFilterOptions();
        }
    }

    showSortOptions() {
        const options = ['Price: Low to High', 'Price: High to Low', 'Rating: High to Low', 'Newest First'];
        const selected = prompt('Choose sorting option:\n' + options.map((opt, i) => `${i + 1}. ${opt}`).join('\n'));
        
        if (selected && !isNaN(selected) && selected >= 1 && selected <= options.length) {
            this.sortProducts(selected);
        }
    }

    showCategoryOptions() {
        const categories = ['All Categories', 'Home & Garden', 'Artificial Plants', 'Planters', 'Decorative Items'];
        const selected = prompt('Choose category:\n' + categories.map((cat, i) => `${i + 1}. ${cat}`).join('\n'));
        
        if (selected && !isNaN(selected) && selected >= 1 && selected <= categories.length) {
            this.filterByCategory(selected - 1);
        }
    }

    showPriceOptions() {
        const ranges = ['All Prices', 'Under ₹200', '₹200 - ₹500', '₹500 - ₹1000', 'Above ₹1000'];
        const selected = prompt('Choose price range:\n' + ranges.map((range, i) => `${i + 1}. ${range}`).join('\n'));
        
        if (selected && !isNaN(selected) && selected >= 1 && selected <= ranges.length) {
            this.filterByPrice(selected - 1);
        }
    }

    showFilterOptions() {
        alert('Advanced filters coming soon!\n\nYou can filter by:\n• Brand\n• Rating\n• Delivery options\n• Offers');
    }

    sortProducts(sortType) {
        switch(sortType) {
            case 1: // Price: Low to High
                this.filteredProducts.sort((a, b) => a.currentPrice - b.currentPrice);
                break;
            case 2: // Price: High to Low
                this.filteredProducts.sort((a, b) => b.currentPrice - a.currentPrice);
                break;
            case 3: // Rating: High to Low
                this.filteredProducts.sort((a, b) => b.rating - a.rating);
                break;
            case 4: // Newest First (by ID)
                this.filteredProducts.sort((a, b) => b.id.localeCompare(a.id));
                break;
        }
        this.displayProducts();
    }

    filterByCategory(categoryIndex) {
        const categories = ['all', 'Home & Garden', 'Artificial Plants', 'Planters', 'Decorative Items'];
        const selectedCategory = categories[categoryIndex];
        
        if (selectedCategory === 'all') {
            this.filteredProducts = [...this.products];
        } else {
            this.filteredProducts = this.products.filter(product => 
                product.category === selectedCategory
            );
        }
        this.displayProducts();
    }

    filterByPrice(priceIndex) {
        const ranges = [
            { min: 0, max: Infinity },
            { min: 0, max: 200 },
            { min: 200, max: 500 },
            { min: 500, max: 1000 },
            { min: 1000, max: Infinity }
        ];
        
        const selectedRange = ranges[priceIndex];
        this.filteredProducts = this.products.filter(product => 
            product.currentPrice >= selectedRange.min && product.currentPrice <= selectedRange.max
        );
        this.displayProducts();
    }

    showProductDetails(productId) {
        // Navigate to product details page
        window.location.href = `product-details.html?id=${productId}`;
    }

    showLoading(show) {
        const loadingIndicator = document.getElementById('loadingIndicator');
        const productsSection = document.querySelector('.products-section');
        
        if (loadingIndicator && productsSection) {
            if (show) {
                loadingIndicator.style.display = 'flex';
                productsSection.style.display = 'none';
            } else {
                loadingIndicator.style.display = 'none';
                productsSection.style.display = 'block';
            }
        }
    }
}

// Initialize the product listing page when DOM is loaded
document.addEventListener('DOMContentLoaded', function() {
    window.productListing = new ProductListing();
});

// Export for use in other files
window.ProductListing = ProductListing; 