// Order Success Page
class OrderSuccess {
    constructor() {
        this.orderData = {};
        this.init();
    }

    init() {
        this.setupEventListeners();
        this.loadOrderData();
        this.updateOrderDetails();
        this.startTimelineAnimation();
    }

    setupEventListeners() {
        // Home button
        const homeBtn = document.getElementById('homeBtn');
        if (homeBtn) {
            homeBtn.addEventListener('click', (e) => {
                e.preventDefault();
                this.goHome();
            });
            
            // Also add touch event for mobile
            homeBtn.addEventListener('touchend', (e) => {
                e.preventDefault();
                this.goHome();
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
                alert('Wishlist feature coming soon!');
            });
        }
        
        if (cartBtn) {
            cartBtn.addEventListener('click', () => {
                alert('Cart feature coming soon!');
            });
        }

        // Action buttons
        const trackOrderBtn = document.getElementById('trackOrderBtn');
        const continueShoppingBtn = document.getElementById('continueShoppingBtn');
        
        if (trackOrderBtn) {
            trackOrderBtn.addEventListener('click', () => {
                this.trackOrder();
            });
        }
        
        if (continueShoppingBtn) {
            continueShoppingBtn.addEventListener('click', () => {
                this.continueShopping();
            });
        }

        // Additional info buttons
        const contactSupportBtn = document.getElementById('contactSupportBtn');
        const viewOrdersBtn = document.getElementById('viewOrdersBtn');
        
        if (contactSupportBtn) {
            contactSupportBtn.addEventListener('click', () => {
                this.contactSupport();
            });
        }
        
        if (viewOrdersBtn) {
            viewOrdersBtn.addEventListener('click', () => {
                this.viewAllOrders();
            });
        }

        // Add touch feedback for mobile
        this.addTouchFeedback();
    }

    loadOrderData() {
        // Get order data from URL parameters or localStorage
        const urlParams = new URLSearchParams(window.location.search);
        const orderId = urlParams.get('orderId') || this.generateOrderId();
        const productId = urlParams.get('productId') || 's-532528444';
        
        // In a real app, you would fetch order data from API
        // For now, we'll use mock data
        this.orderData = {
            orderId: orderId,
            productId: productId,
            orderDate: new Date(),
            orderTotal: 163,
            paymentMethod: 'Online Payment',
            deliveryAddress: 'Home Address',
            expectedDelivery: '3-5 business days',
            product: {
                title: 'Artificial Natural Looking 26 Leaves Snack Rabbur Plant',
                price: 163,
                quantity: 1,
                image: '../assets/icon.png'
            },
            status: 'confirmed'
        };
    }

    generateOrderId() {
        const timestamp = Date.now().toString().slice(-6);
        const random = Math.floor(Math.random() * 1000).toString().padStart(3, '0');
        return `MEESH${timestamp}${random}`;
    }

    updateOrderDetails() {
        // Update order number
        const orderNumber = document.querySelector('.order-number');
        if (orderNumber) {
            orderNumber.textContent = `Order #${this.orderData.orderId}`;
        }

        // Update order date
        const orderDate = document.getElementById('orderDate');
        const orderTime = document.getElementById('orderTime');
        if (orderDate) {
            orderDate.textContent = this.formatDate(this.orderData.orderDate);
        }
        if (orderTime) {
            orderTime.textContent = this.formatTime(this.orderData.orderDate);
        }

        // Update order total
        const orderTotal = document.querySelector('.summary-value');
        if (orderTotal) {
            orderTotal.textContent = `₹${this.orderData.orderTotal}`;
        }

        // Update product details
        this.updateProductDetails();

        // Clear cart badge since order is placed
        const cartBadge = document.querySelector('.cart-badge');
        if (cartBadge) {
            cartBadge.textContent = '0';
        }
    }

    updateProductDetails() {
        const productTitle = document.querySelector('.product-title');
        const productPrice = document.querySelector('.product-price');
        const productQuantity = document.querySelector('.product-quantity');
        const productImage = document.querySelector('.product-image');

        if (productTitle) {
            productTitle.textContent = this.orderData.product.title;
        }
        if (productPrice) {
            productPrice.textContent = `₹${this.orderData.product.price}`;
        }
        if (productQuantity) {
            productQuantity.textContent = `Qty: ${this.orderData.product.quantity}`;
        }
        if (productImage) {
            productImage.src = this.orderData.product.image;
        }
    }

    formatDate(date) {
        return date.toLocaleDateString('en-IN', {
            year: 'numeric',
            month: 'short',
            day: 'numeric'
        });
    }

    formatTime(date) {
        return date.toLocaleTimeString('en-IN', {
            hour: '2-digit',
            minute: '2-digit',
            hour12: true
        });
    }

    startTimelineAnimation() {
        // Animate timeline items with delay
        const timelineItems = document.querySelectorAll('.timeline-item');
        
        timelineItems.forEach((item, index) => {
            setTimeout(() => {
                item.style.opacity = '0';
                item.style.transform = 'translateX(-20px)';
                item.style.transition = 'all 0.5s ease';
                
                setTimeout(() => {
                    item.style.opacity = '1';
                    item.style.transform = 'translateX(0)';
                }, 100);
            }, index * 300);
        });
    }

    trackOrder() {
        alert(`Track Order: ${this.orderData.orderId}\n\nCurrent Status: Order Confirmed\nExpected Delivery: ${this.orderData.expectedDelivery}\n\nTracking feature coming soon!`);
    }

    continueShopping() {
        // Navigate back to home page
        window.location.href = 'home.html';
    }

    goHome() {
        // Navigate to home page
        window.location.href = 'home.html';
    }

    contactSupport() {
        const supportOptions = [
            '📞 Call Support: 1800-123-4567',
            '💬 Live Chat: Available 24/7',
            '📧 Email: support@meesho.com',
            '📱 WhatsApp: +91-98765-43210'
        ];
        
        alert('Contact Support:\n\n' + supportOptions.join('\n') + '\n\nWe\'re here to help you!');
    }

    viewAllOrders() {
        alert('View All Orders:\n\nThis would show your complete order history with tracking information for all past orders.\n\nFeature coming soon!');
    }

    addTouchFeedback() {
        // Add touch feedback for interactive elements
        const touchElements = document.querySelectorAll('.primary-btn, .secondary-btn, .info-btn, .header-btn');
        
        touchElements.forEach(element => {
            element.addEventListener('touchstart', () => {
                element.style.transform = 'scale(0.95)';
            });
            
            element.addEventListener('touchend', () => {
                element.style.transform = 'scale(1)';
            });
        });

        // Prevent zoom on double tap for buttons
        document.addEventListener('touchend', (e) => {
            if (e.target.tagName === 'BUTTON') {
                e.preventDefault();
            }
        });
    }

    // Method to simulate order processing (for demo purposes)
    simulateOrderProcessing() {
        const timelineItems = document.querySelectorAll('.timeline-item');
        let currentStep = 0;

        const processNextStep = () => {
            if (currentStep < timelineItems.length) {
                // Remove active class from all items
                timelineItems.forEach(item => item.classList.remove('active'));
                
                // Add active class to current step
                timelineItems[currentStep].classList.add('active');
                
                currentStep++;
                
                // Process next step after delay
                setTimeout(processNextStep, 2000);
            }
        };

        // Start processing after 3 seconds
        setTimeout(processNextStep, 3000);
    }
}

// Initialize the order success page when DOM is loaded
document.addEventListener('DOMContentLoaded', function() {
    window.orderSuccess = new OrderSuccess();
    
    // Optional: Simulate order processing for demo
    // window.orderSuccess.simulateOrderProcessing();
});

// Export for use in other files
window.OrderSuccess = OrderSuccess; 