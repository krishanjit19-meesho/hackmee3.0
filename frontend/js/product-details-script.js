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
        if (backBtn) {
            backBtn.addEventListener('click', () => {
                this.goBack();
            });
        }

        // Header buttons
        const searchBtn = document.querySelector('.search-btn');
        const wishlistBtn = document.querySelector('.wishlist-btn');
        const cartBtn = document.querySelector('.cart-btn');
        
        if (searchBtn) {
            searchBtn.addEventListener('click', () => {
                alert('Search feature coming soon!');
            });
        }
        
        if (wishlistBtn) {
            wishlistBtn.addEventListener('click', () => {
                this.toggleWishlist();
            });
        }
        
        if (cartBtn) {
            cartBtn.addEventListener('click', () => {
                alert('Cart feature coming soon!');
            });
        }

        // Image navigation dots
        document.querySelectorAll('.dot').forEach((dot, index) => {
            dot.addEventListener('click', () => {
                this.changeImage(index);
            });
        });

        // Similar product thumbnails
        document.querySelectorAll('.similar-thumbnail').forEach((thumbnail, index) => {
            thumbnail.addEventListener('click', () => {
                this.selectSimilarProduct(index);
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

        // View reviews button
        const viewReviewsBtn = document.querySelector('.view-reviews-btn');
        if (viewReviewsBtn) {
            viewReviewsBtn.addEventListener('click', () => {
                this.viewAllReviews();
            });
        }

        // Add touch feedback for mobile
        this.addTouchFeedback();
    }

    loadProductData() {
        // Get product ID from URL or localStorage
        const urlParams = new URLSearchParams(window.location.search);
        const productId = urlParams.get('id') || 's-532528444';
        
        // In a real app, you would fetch product data from API
        // For now, we'll use mock data
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

        this.updateProductDisplay();
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

        // Update specifications
        this.updateSpecifications();

        // Update reviews
        const ratingText = document.querySelector('.rating-text');
        const ratingCount = document.querySelector('.rating-count');
        
        if (ratingText) ratingText.textContent = `${this.productData.rating} out of 5`;
        if (ratingCount) ratingCount.textContent = `(${this.productData.reviews.toLocaleString()} reviews)`;
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
        if (index < 0 || index >= this.productImages.length) return;
        
        this.currentImageIndex = index;
        const mainImage = document.getElementById('mainImage');
        
        if (mainImage) {
            mainImage.src = this.productImages[index];
        }

        // Update dots
        document.querySelectorAll('.dot').forEach((dot, i) => {
            dot.classList.toggle('active', i === index);
        });
    }

    selectSimilarProduct(index) {
        document.querySelectorAll('.similar-thumbnail').forEach((thumbnail, i) => {
            thumbnail.classList.toggle('active', i === index);
        });
        
        // In a real app, this would load a different product
        alert(`Similar product ${index + 1} selected!`);
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

    buyNow() {
        // Redirect to order success page with order details
        const orderId = this.generateOrderId();
        window.location.href = `order-success.html?orderId=${orderId}&productId=${this.productData.id}`;
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

    viewAllReviews() {
        alert('Customer Reviews:\n\n★★★★☆ John D. - "Great quality, looks very natural!"\n★★★★★ Sarah M. - "Perfect for my office desk"\n★★★★☆ Mike R. - "Good value for money"\n★★★★★ Lisa K. - "Beautiful and low maintenance"\n\nView all 1,234 reviews...');
    }

    goBack() {
        // Check if we came from product listing page
        const referrer = document.referrer;
        if (referrer && referrer.includes('product-listing.html')) {
            // Go back to product listing page
            window.location.href = 'product-listing.html';
        } else {
            // Default fallback - go back in history
            window.history.back();
        }
    }

    addTouchFeedback() {
        // Add touch feedback for interactive elements
        const touchElements = document.querySelectorAll('.action-btn, .dot, .similar-thumbnail, .add-to-cart-btn, .buy-now-btn, .view-reviews-btn');
        
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