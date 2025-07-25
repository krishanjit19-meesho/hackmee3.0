// Product Details Page
class ProductDetails {
    constructor() {
        this.currentImageIndex = 0;
        this.productImages = [
            '../assets/icon.png',
            '../assets/icon.png',
            '../assets/icon.png',
            '../assets/icon.png'
        ];
        this.init();
    }

    init() {
        this.setupEventListeners();
        this.loadProductData();
    }

    setupEventListeners() {
                // Back button
        const backBtn = document.getElementById('backBtn');
        console.log('Back button found:', !!backBtn);
        if (backBtn) {
            backBtn.addEventListener('click', (e) => {
                console.log('Back button clicked!');
                e.preventDefault();
                this.goBack();
            });
            
            // Also add touch event for mobile
            backBtn.addEventListener('touchend', (e) => {
                console.log('Back button touched!');
                e.preventDefault();
                this.goBack();
            });
        }

        // Cart button
        const cartBtn = document.querySelector('.cart-btn');
        console.log('Cart button found:', !!cartBtn);
        if (cartBtn) {
            cartBtn.addEventListener('click', (e) => {
                console.log('Cart button clicked!');
                e.preventDefault();
                this.openCart();
            });
            
            // Also add touch event for mobile
            cartBtn.addEventListener('touchend', (e) => {
                console.log('Cart button touched!');
                e.preventDefault();
                this.openCart();
            });
        }

        // Image navigation dots
        document.querySelectorAll('.dot').forEach((dot, index) => {
            dot.addEventListener('click', () => {
                this.changeImage(index);
            });
        });



        // Product action buttons
        const wishlistAction = document.querySelector('.wishlist-action');
        const shareAction = document.querySelector('.share-action');
        
        if (wishlistAction) {
            wishlistAction.addEventListener('click', () => {
                this.toggleWishlist();
            });
        }
        
        if (shareAction) {
            shareAction.addEventListener('click', () => {
                this.shareProduct();
            });
        }

        // Special offer click
        const specialOffer = document.querySelector('.special-offer');
        if (specialOffer) {
            specialOffer.addEventListener('click', () => {
                this.showSpecialOffers();
            });
        }

        // Action buttons
        const addToCartBtn = document.querySelector('.add-to-cart-btn');
        const buyNowBtn = document.querySelector('.buy-now-btn');
        
        console.log('Add to Cart button found:', !!addToCartBtn);
        console.log('Buy Now button found:', !!buyNowBtn);
        
        if (addToCartBtn) {
            addToCartBtn.addEventListener('click', () => {
                console.log('Add to Cart clicked');
                this.addToCart();
            });
        }
        
        if (buyNowBtn) {
            buyNowBtn.addEventListener('click', () => {
                console.log('Buy Now clicked');
                this.buyNow();
            });
            
            // Add a simple test click handler
            buyNowBtn.onclick = () => {
                console.log('Buy Now onclick triggered');
                this.buyNow();
            };
        }

        // Location banner
        const locationBanner = document.querySelector('.location-banner');
        if (locationBanner) {
            locationBanner.addEventListener('click', () => {
                this.addDeliveryLocation();
            });
        }

        // Breadcrumb navigation
        document.querySelectorAll('.breadcrumb-item').forEach((item, index) => {
            item.addEventListener('click', () => {
                this.navigateBreadcrumb(index);
            });
        });



        // Add touch feedback for mobile
        this.addTouchFeedback();
    }

    async loadProductData() {
        // Get product ID from URL or localStorage
        const urlParams = new URLSearchParams(window.location.search);
        const productId = urlParams.get('id') || 's-532528444';
        
        // Get user data from localStorage
        const userDataString = localStorage.getItem('userData');
        let userID = 'user123'; // Default fallback
        
        if (userDataString) {
            try {
                const userData = JSON.parse(userDataString);
                userID = userData.userId || userID;
            } catch (error) {
                console.error('Error parsing user data:', error);
            }
        }

        try {
            // Show loading state
            this.showLoading(true);
            
            // Call the product details API
            const response = await fetch(`http://localhost:8080/api/v1/product/details?product_id=${productId}&user_id=${userID}`, {
                method: 'GET',
                headers: {
                    'Content-Type': 'application/json',
                }
            });

            const result = await response.json();

            if (response.ok && result.success) {
                // Use the API response data
                this.productData = this.convertApiToProductData(result.data);
                console.log('Loaded product data from backend:', result);
            } else {
                console.error('Failed to load product data:', result.error);
                // Fall back to mock data
                this.loadMockProductData(productId);
            }
        } catch (error) {
            console.error('Error loading product data:', error);
            // Fall back to mock data
            this.loadMockProductData(productId);
        } finally {
            this.showLoading(false);
        }

        this.updateProductDisplay();
    }

    convertApiToProductData(apiData) {
        return {
            id: apiData.product_id,
            catalogId: apiData.catalog_id,
            title: apiData.title,
            description: apiData.description,
            category: apiData.category,
            subcategory: apiData.sub_category,
            currentPrice: this.extractPriceValue(apiData.price),
            originalPrice: this.extractPriceValue(apiData.original_price),
            discount: apiData.discount,
            discountPercent: apiData.discount_percent,
            specialOffer: `₹${Math.floor(this.extractPriceValue(apiData.price) * 0.9)} with ${Math.floor(Math.random() * 3) + 1} Special Offers`,
            delivery: apiData.delivery_info || 'Free Delivery',
            rating: apiData.rating,
            reviews: apiData.reviews,
            images: apiData.images && apiData.images.length > 0 ? apiData.images : this.productImages,
            mainImage: apiData.main_image || this.productImages[0],
            brand: apiData.brand,
            seller: apiData.seller,
            stock: apiData.stock,
            specifications: apiData.specifications || {},
            variants: apiData.variants || [],
            reviewsList: apiData.reviews_list || [],
            similarProducts: apiData.similar_products || [],
            returnPolicy: apiData.return_policy,
            warranty: apiData.warranty
        };
    }

    extractPriceValue(priceString) {
        if (!priceString) return 0;
        const match = priceString.match(/₹?(\d+)/);
        return match ? parseInt(match[1]) : 0;
    }

    loadMockProductData(productId) {
        // Fallback mock data
        this.productData = {
            id: productId,
            title: 'Artificial Natural Looking 26 Leaves Snack Rabbur Plant with Pot for Home Office Decoration',
            currentPrice: 163,
            originalPrice: 175,
            discount: '7% off',
            specialOffer: '₹147 with 1 Special Offer',
            delivery: 'Free Delivery',
            rating: 4.2,
            reviews: 1234,
            images: this.productImages,
            category: 'Home & Kitchen',
            subcategory: 'Artificial Flora',
            specifications: {
                material: 'Plastic',
                height: '45 cm',
                potSize: '15 cm',
                leaves: '26'
            },
            description: 'High-quality artificial plant perfect for home and office decoration. This natural-looking plant features 26 leaves and comes with a stylish pot. No maintenance required, looks fresh all year round.'
        };
    }

    showLoading(show) {
        // Add loading indicator to the page
        let loadingIndicator = document.getElementById('loadingIndicator');
        
        if (!loadingIndicator) {
            loadingIndicator = document.createElement('div');
            loadingIndicator.id = 'loadingIndicator';
            loadingIndicator.innerHTML = `
                <div style="position: fixed; top: 0; left: 0; right: 0; bottom: 0; background: rgba(255,255,255,0.9); display: flex; align-items: center; justify-content: center; z-index: 9999;">
                    <div style="text-align: center;">
                        <div style="width: 40px; height: 40px; border: 4px solid #f3f3f3; border-top: 4px solid #9f2089; border-radius: 50%; animation: spin 1s linear infinite; margin: 0 auto 10px;"></div>
                        <div>Loading product details...</div>
                    </div>
                </div>
                <style>
                    @keyframes spin {
                        0% { transform: rotate(0deg); }
                        100% { transform: rotate(360deg); }
                    }
                </style>
            `;
            document.body.appendChild(loadingIndicator);
        }
        
        loadingIndicator.style.display = show ? 'block' : 'none';
    }

    updateProductDisplay() {
        // Update product title
        const productTitle = document.querySelector('.product-title');
        if (productTitle) {
            productTitle.textContent = this.productData.title;
        }

        // Update product ID
        const productId = document.querySelector('.product-id');
        if (productId) {
            productId.textContent = this.productData.id;
        }

        // Update main image
        const mainImage = document.getElementById('mainImage');
        if (mainImage && this.productData.mainImage) {
            mainImage.src = this.productData.mainImage;
        }

        // Update image dots based on available images
        this.updateImageDots();

        // Update pricing
        const currentPrice = document.querySelector('.current-price');
        const originalPrice = document.querySelector('.original-price');
        const discount = document.querySelector('.discount');
        const specialOffer = document.querySelector('.offer-text');
        const deliveryInfo = document.querySelector('.delivery-info');

        if (currentPrice) currentPrice.textContent = `₹${this.productData.currentPrice}`;
        if (originalPrice) originalPrice.textContent = `₹${this.productData.originalPrice}`;
        if (discount) discount.textContent = this.productData.discount;
        if (specialOffer) specialOffer.textContent = this.productData.specialOffer;
        if (deliveryInfo) deliveryInfo.textContent = this.productData.delivery;

        // Update breadcrumb
        this.updateBreadcrumb();

        // Update specifications
        this.updateSpecifications();

        // Update reviews
        const ratingText = document.querySelector('.rating-text');
        const ratingCount = document.querySelector('.rating-count');
        
        if (ratingText) ratingText.textContent = `${this.productData.rating.toFixed(1)} out of 5`;
        if (ratingCount) ratingCount.textContent = `(${this.productData.reviews.toLocaleString()} reviews)`;



        // Update product details section
        this.updateProductDetailsSection();
    }

    updateBreadcrumb() {
        const breadcrumbContent = document.querySelector('.breadcrumb-content');
        if (breadcrumbContent && this.productData.category && this.productData.subcategory) {
            breadcrumbContent.innerHTML = `
                <span class="breadcrumb-item">Home</span>
                <span class="breadcrumb-separator">/</span>
                <span class="breadcrumb-item">${this.productData.category}</span>
                <span class="breadcrumb-separator">/</span>
                <span class="breadcrumb-item">${this.productData.subcategory}</span>
                <span class="breadcrumb-separator">/</span>
                <span class="breadcrumb-item">${this.productData.title.substring(0, 20)}...</span>
            `;
        }
    }



    updateImageDots() {
        const imageDots = document.querySelector('.image-dots');
        if (!imageDots || !this.productData.images) return;

        // Clear existing dots
        imageDots.innerHTML = '';

        // Create dots for each available image (up to 5 total: main + 4 additional)
        const totalImages = Math.min(this.productData.images.length, 5);
        
        for (let i = 0; i < totalImages; i++) {
            const dot = document.createElement('div');
            dot.className = `dot ${i === 0 ? 'active' : ''}`;
            dot.setAttribute('data-image', i.toString());
            
            // Add click event
            dot.addEventListener('click', () => {
                this.changeImage(i);
            });
            
            imageDots.appendChild(dot);
        }
    }

    updateProductDetailsSection() {
        const detailContent = document.querySelector('.detail-content');
        if (detailContent && this.productData.description) {
            detailContent.innerHTML = `<p>${this.productData.description}</p>`;
        }
    }

    updateSpecifications() {
        const specsGrid = document.querySelector('.specs-grid');
        if (!specsGrid) return;

        const specs = this.productData.specifications;
        const specItems = Object.entries(specs).map(([key, value]) => `
            <div class="spec-item">
                <span class="spec-label">${this.capitalizeFirst(key)}:</span>
                <span class="spec-value">${value}</span>
            </div>
        `).join('');

        specsGrid.innerHTML = specItems;
    }

    capitalizeFirst(str) {
        return str.charAt(0).toUpperCase() + str.slice(1);
    }

    changeImage(index) {
        if (!this.productData.images || index < 0 || index >= this.productData.images.length) return;
        
        this.currentImageIndex = index;
        const mainImage = document.getElementById('mainImage');
        
        if (mainImage) {
            mainImage.src = this.productData.images[index];
        }

        // Update dots
        document.querySelectorAll('.dot').forEach((dot, i) => {
            dot.classList.toggle('active', i === index);
        });
    }



    toggleWishlist() {
        const wishlistBtns = document.querySelectorAll('.wishlist-btn, .wishlist-action');
        
        wishlistBtns.forEach(btn => {
            if (btn.textContent.includes('❤️')) {
                btn.textContent = btn.textContent.replace('❤️', '🤍');
                btn.style.color = '#666';
            } else {
                btn.textContent = btn.textContent.replace('🤍', '❤️');
                btn.style.color = '#e91e63';
            }
        });

        const isWishlisted = wishlistBtns[0].textContent.includes('❤️');
        alert(isWishlisted ? 'Added to wishlist!' : 'Removed from wishlist!');
    }

    shareProduct() {
        if (navigator.share) {
            navigator.share({
                title: this.productData.title,
                text: `Check out this amazing product: ${this.productData.title}`,
                url: window.location.href
            });
        } else {
            // Fallback for browsers that don't support Web Share API
            const shareUrl = `https://wa.me/?text=Check out this amazing product: ${this.productData.title} - ${window.location.href}`;
            window.open(shareUrl, '_blank');
        }
    }

    showSpecialOffers() {
        alert('Special Offers:\n\n• Get ₹16 off on orders above ₹500\n• Use code: SAVE16\n• Valid till month end\n• Free delivery on all orders');
    }

    addToCart() {
        alert('Product added to cart!\n\nYou can continue shopping or proceed to checkout.');
        
        // Update cart badge
        const cartBadge = document.querySelector('.cart-badge');
        if (cartBadge) {
            const currentCount = parseInt(cartBadge.textContent) || 0;
            cartBadge.textContent = currentCount + 1;
        }
    }

    async buyNow() {
        console.log('=== BUY NOW DEBUG START ===');
        
        try {
            // Debug: Check if productData exists
            console.log('Product Data:', this.productData);
            
            if (!this.productData) {
                console.error('No product data available');
                alert('Product data not loaded. Please refresh the page.');
                return;
            }

            // Get user data from localStorage
            const userDataString = localStorage.getItem('userData');
            console.log('User Data String from localStorage:', userDataString);
            
            if (!userDataString) {
                console.error('No user data found in localStorage');
                alert('Please login first');
                return;
            }

            const userData = JSON.parse(userDataString);
            console.log('Parsed User Data:', userData);
            
            if (!userData.userId) {
                console.error('No userId in user data');
                alert('Invalid user data. Please login again.');
                return;
            }

            // Prepare request payload
            const requestPayload = {
                user_id: userData.userId,
                product_id: this.productData.id,
                catalog_id: this.productData.catalogId || '12345',
                quantity: 1
            };
            
            console.log('Request Payload:', requestPayload);
            console.log('Request URL:', 'http://localhost:8080/api/v1/order/place');
            
            // Call the place order API
            const response = await fetch('http://localhost:8080/api/v1/order/place', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(requestPayload)
            });

            console.log('Response Status:', response.status);
            console.log('Response Headers:', response.headers);

            const result = await response.json();
            console.log('Response Body:', result);

            if (response.ok && result.success) {
                console.log('Order placed successfully:', result.data);
                // Redirect to order success page with order details
                window.location.href = `order-success.html?orderId=${result.data.order_id}&productId=${this.productData.id}&quantity=1`;
            } else {
                console.error('Failed to place order:', result);
                const errorMessage = result.error || result.message || 'Unknown error occurred';
                alert(`Failed to place order: ${errorMessage}`);
            }
        } catch (error) {
            console.error('Error placing order:', error);
            console.error('Error details:', {
                name: error.name,
                message: error.message,
                stack: error.stack
            });
            alert('Network error. Please check your connection and try again.');
        }
        
        console.log('=== BUY NOW DEBUG END ===');
    }

    generateOrderId() {
        const timestamp = Date.now().toString().slice(-6);
        const random = Math.floor(Math.random() * 1000).toString().padStart(3, '0');
        return `MEESH${timestamp}${random}`;
    }

    addDeliveryLocation() {
        const location = prompt('Enter your delivery location (city/pincode):');
        if (location) {
            alert(`Location updated: ${location}\n\nExtra discounts applied!`);
            
            // Update location banner
            const locationText = document.querySelector('.location-text');
            if (locationText) {
                locationText.textContent = `Delivery to ${location} - Extra discount applied`;
            }
        }
    }

    navigateBreadcrumb(index) {
        const breadcrumbs = ['Home', 'Home & Kitchen', 'Artificial Flora', 'Artificial Plants, Flowers And Shrubs', 'Artificial Natur...'];
        
        if (index === 0) {
            window.location.href = 'home.html';
        } else {
            alert(`Navigating to: ${breadcrumbs[index]}`);
        }
    }



    goBack() {
        // Check if we came from product listing page
        const referrer = document.referrer;
        if (referrer && referrer.includes('product-listing.html')) {
            // Go back to product listing page
            window.location.href = 'product-listing.html';
        } else if (referrer && referrer.includes('home.html')) {
            // Go back to home page
            window.location.href = 'home.html';
        } else {
            // Default fallback - go back in history
            window.history.back();
        }
    }

    openCart() {
        alert('Cart Feature:\n\n• 2 items in cart\n• Total: ₹326\n• Free delivery available\n\nCart page coming soon!');
    }

    addTouchFeedback() {
        // Add touch feedback for interactive elements
        const touchElements = document.querySelectorAll('.action-btn, .dot, .add-to-cart-btn, .buy-now-btn');
        
        console.log('Touch elements found:', touchElements.length);
        
        touchElements.forEach(element => {
            element.addEventListener('touchstart', () => {
                element.style.transform = 'scale(0.95)';
            });
            
            element.addEventListener('touchend', () => {
                element.style.transform = 'scale(1)';
            });
        });

        // Prevent zoom on double tap for inputs
        document.addEventListener('touchend', (e) => {
            if (e.target.tagName === 'INPUT' || e.target.tagName === 'BUTTON') {
                e.preventDefault();
            }
        });
    }
}

// Initialize the product details page when DOM is loaded
document.addEventListener('DOMContentLoaded', function() {
    window.productDetails = new ProductDetails();
});

// Export for use in other files
window.ProductDetails = ProductDetails; 