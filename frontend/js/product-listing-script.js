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
        // Constants for page titles
        this.PAGE_TITLES = {
            DEFAULT: 'Special Offers',
            CATALOG: 'Catalog Products',
            SEARCH: 'Search Results',
            CATEGORY: 'Category Products',
            FEATURED: 'Featured Products'
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
        
        try {
            // Get catalog data from localStorage (set by home page)
            const catalogDataString = localStorage.getItem('catalogData');
            
            if (catalogDataString) {
                const catalogData = JSON.parse(catalogDataString);
                
                if (catalogData.success && catalogData.data && catalogData.data.length > 0) {
                    // Convert backend catalog data to product format
                    this.products = this.convertCatalogToProducts(catalogData.data);
                    console.log('Loaded catalog data from backend:', catalogData);
                    
                    // Update page title to show catalog information
                    this.updatePageTitle(catalogData);
                } else {
                    console.warn('No catalog data available, using mock data');
                    this.loadMockProducts();
                    this.setPageTitle(this.PAGE_TITLES.DEFAULT);
                }
            } else {
                console.warn('No catalog data in localStorage, using mock data');
                this.loadMockProducts();
                this.setPageTitle(this.PAGE_TITLES.DEFAULT);
            }
        } catch (error) {
            console.error('Error loading catalog data:', error);
            this.loadMockProducts();
            this.setPageTitle(this.PAGE_TITLES.DEFAULT);
        }

        this.filteredProducts = [...this.products];
        this.displayProducts();
        this.showLoading(false);
    }

    updatePageTitle(catalogData) {
        // Always use "Special Offers" regardless of data source
        this.setPageTitle(this.PAGE_TITLES.DEFAULT);
    }

    setPageTitle(title) {
        const pageTitle = document.getElementById('pageTitle');
        if (pageTitle) {
            pageTitle.textContent = title;
        }
    }

    convertCatalogToProducts(catalogData) {
        return catalogData.map((item, index) => {
            // Parse price data from backend response
            let currentPrice = 0;
            let originalPrice = 0;
            let discountPercent = 0;
            
            // Extract numeric values from price strings (remove ₹ symbol and convert to number)
            if (item.price) {
                currentPrice = parseInt(item.price.replace('₹', '').replace(/,/g, '')) || 0;
            }
            
            if (item.original_price) {
                originalPrice = parseInt(item.original_price.replace('₹', '').replace(/,/g, '')) || 0;
            }
            
            if (item.discount_percent) {
                discountPercent = item.discount_percent;
            } else if (originalPrice > 0 && currentPrice > 0) {
                // Calculate discount percentage if not provided
                discountPercent = Math.floor(((originalPrice - currentPrice) / originalPrice) * 100);
            }
            
            // Fallback to random prices if backend data is missing
            if (currentPrice === 0 || originalPrice === 0) {
                currentPrice = Math.floor(Math.random() * 500) + 100;
                originalPrice = Math.floor(currentPrice * (1 + Math.random() * 0.5));
                discountPercent = Math.floor(((originalPrice - currentPrice) / originalPrice) * 100);
            }
            
            return {
                id: item.product_id || `catalog-${index}`,
                title: item.title || `${item.category} - ${item.sub_category}`,
                currentPrice: currentPrice,
                originalPrice: originalPrice,
                discount: `${discountPercent}% off`,
                specialOffer: `₹${Math.floor(currentPrice * 0.9)} with ${Math.floor(Math.random() * 3) + 1} Special Offers`,
                delivery: 'Free Delivery',
                rating: (Math.random() * 2 + 3).toFixed(1), // Random rating between 3-5
                reviews: Math.floor(Math.random() * 50000) + 1000,
                image: item.image_url || '../assets/icon.png',
                category: item.category || 'General',
                subCategory: item.sub_category || '',
                catalogId: item.catalog_id,
                price: item.price,
                isTrusted: Math.random() > 0.7 // 30% chance of being trusted
            };
        });
    }

    loadMockProducts() {
        // Mock product data similar to Meesho style (fallback)
        this.products = [
            {
                id: 's-375061729',
                title: 'Trendy Retro Men Shirts',
                currentPrice: 109,
                originalPrice: 199,
                discount: '45% off',
                specialOffer: '₹84 with 2 Special Offers',
                delivery: 'Free Delivery',
                rating: 4.0,
                reviews: 46686,
                image: '../assets/icon.png'
            },
            {
                id: 's-510539001',
                title: 'Selfie Sticks',
                currentPrice: 187,
                originalPrice: 299,
                discount: '37% off',
                specialOffer: '₹163 with 1 Special Offer',
                delivery: 'Free Delivery',
                rating: 4.1,
                reviews: 16237,
                image: '../assets/icon.png'
            },
            {
                id: 's-191005349',
                title: 'Chitrarekha Attractive Kurtis',
                currentPrice: 184,
                originalPrice: 399,
                discount: '54% off',
                specialOffer: '₹165 with 3 Special Offers',
                delivery: 'Free Delivery',
                rating: 4.2,
                reviews: 8923,
                image: '../assets/icon.png'
            },
            {
                id: 's-507286436',
                title: 'Urbane Graceful Women Tops & T-Shirts',
                currentPrice: 292,
                originalPrice: 599,
                discount: '51% off',
                specialOffer: '₹263 with 2 Special Offers',
                delivery: 'Free Delivery',
                rating: 4.3,
                reviews: 15432,
                image: '../assets/icon.png'
            },
            {
                id: 's-123456789',
                title: 'Stylish Men Casual Shirts',
                currentPrice: 199,
                originalPrice: 399,
                discount: '50% off',
                specialOffer: '₹179 with 1 Special Offer',
                delivery: 'Free Delivery',
                rating: 4.4,
                reviews: 23456,
                image: '../assets/icon.png'
            },
            {
                id: 's-987654321',
                title: 'Elegant Women Kurtis Collection',
                currentPrice: 249,
                originalPrice: 499,
                discount: '50% off',
                specialOffer: '₹224 with 2 Special Offers',
                delivery: 'Free Delivery',
                rating: 4.1,
                reviews: 18765,
                image: '../assets/icon.png'
            }
        ];
    }

    displayProducts() {
        const container = document.getElementById('productsGrid');
        if (!container) return;

        container.innerHTML = this.filteredProducts.map(product => `
            <div class="product-card" data-product-id="${product.id}">
                <div class="wishlist-heart">♡</div>
                <div class="product-image-container">
                    <img src="${product.image}" alt="${product.title}" class="product-image" onerror="this.src='../assets/icon.png'">
                </div>
                
                <div class="product-title">${product.title}</div>
                
                <div class="product-price">
                    <span class="current-price">₹${product.currentPrice}</span>
                    <span class="original-price">₹${product.originalPrice}</span>
                    <span class="discount">${product.discount}</span>
                </div>
                
                <div class="special-offer">${product.specialOffer}</div>
                <div class="delivery-info">${product.delivery}</div>
                
                <div class="product-rating">
                    <span class="rating-stars">★</span>
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