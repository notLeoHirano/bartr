const CONFIG = {
    API_URL: 'http://localhost:8080',
    STORAGE_KEY: 'bartr_token'
};

const state = {
    token: localStorage.getItem(CONFIG.STORAGE_KEY),
    currentUser: null,
    currentItems: [],
    currentIndex: 0
};

// Authentication Functions
function showAuthTab(tab) {
    document.querySelectorAll('.auth-tab').forEach(t => {
        t.classList.remove('bg-primary', 'text-white');
        t.classList.add('bg-gray-100', 'text-gray-600');
    });
    
    const target = event.currentTarget || document.querySelector(`button[onclick*="showAuthTab('${tab}')"]`);
    if (target) {
        target.classList.remove('bg-gray-100', 'text-gray-600');
        target.classList.add('bg-primary', 'text-white');
    }

    document.getElementById('loginForm').classList.toggle('hidden', tab !== 'login');
    document.getElementById('registerForm').classList.toggle('hidden', tab !== 'register');
    document.getElementById('authError').classList.add('hidden');
}

async function register(e) {
    e.preventDefault();
    const name = document.getElementById('registerName').value;
    const email = document.getElementById('registerEmail').value;
    const password = document.getElementById('registerPassword').value;

    try {
        const response = await fetch(`${CONFIG.API_URL}/auth/register`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ name, email, password })
        });

        const data = await response.json();

        if (!response.ok) {
            showAuthError(data.error);
            return;
        }

        state.token = data.token;
        state.currentUser = data.user;
        localStorage.setItem(CONFIG.STORAGE_KEY, state.token);
        showApp();
    } catch (error) {
        showAuthError('Registration failed. Please try again.');
    }
}

async function login(e) {
    e.preventDefault();
    const email = document.getElementById('loginEmail').value;
    const password = document.getElementById('loginPassword').value;

    try {
        const response = await fetch(`${CONFIG.API_URL}/auth/login`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ email, password })
        });

        const data = await response.json();

        if (!response.ok) {
            showAuthError(data.error);
            return;
        }

        state.token = data.token;
        state.currentUser = data.user;
        localStorage.setItem(CONFIG.STORAGE_KEY, state.token);
        showApp();
    } catch (error) {
        showAuthError('Login failed. Please try again.');
    }
}

function logout() {
    state.token = null;
    state.currentUser = null;
    localStorage.removeItem(CONFIG.STORAGE_KEY);
    document.getElementById('authScreen').classList.remove('hidden');
    document.getElementById('appScreen').classList.add('hidden');
}

function showAuthError(message) {
    const errorDiv = document.getElementById('authError');
    errorDiv.textContent = message;
    errorDiv.classList.remove('hidden');
}

async function checkAuth() {
    if (!state.token) {
        document.getElementById('authScreen').classList.remove('hidden');
        return;
    }

    try {
        const response = await fetch(`${CONFIG.API_URL}/me`, {
            headers: { 'Authorization': `Bearer ${state.token}` }
        });

        if (!response.ok) {
            logout();
            return;
        }

        state.currentUser = await response.json();
        showApp();
    } catch (error) {
        logout();
    }
}

function showApp() {
    document.getElementById('authScreen').classList.add('hidden');
    document.getElementById('appScreen').classList.remove('hidden');
    document.getElementById('userName').textContent = state.currentUser.name;
    loadSwipeItems();
}

// Navigation Functions
function showTab(tabName) {
    document.querySelectorAll('.tab').forEach(t => {
        t.classList.remove('active', 'bg-white', 'text-primary');
        t.classList.add('bg-white/20', 'text-white');
    });
    document.querySelectorAll('.card').forEach(c => c.classList.add('hidden'));

    const activeTabButton = document.querySelector(`button[onclick*="showTab('${tabName}')"]`);
    if (activeTabButton) {
        activeTabButton.classList.remove('bg-white/20', 'text-white', 'hover:bg-white/30');
        activeTabButton.classList.add('active', 'bg-white', 'text-primary');
    }
    
    document.getElementById(tabName + 'Tab').classList.remove('hidden');

    if (tabName === 'swipe') loadSwipeItems();
    if (tabName === 'items') showItemsSubTab('add');
    if (tabName === 'matches') loadMatches();
}

function showItemsSubTab(subTab) {
    document.querySelectorAll('.items-tab').forEach(t => {
        t.classList.remove('bg-primary', 'text-white', 'active');
        t.classList.add('bg-gray-100', 'text-gray-600');
    });
    
    const activeBtn = document.querySelector(`button[onclick*="showItemsSubTab('${subTab}')"]`);
    if (activeBtn) {
        activeBtn.classList.remove('bg-gray-100', 'text-gray-600');
        activeBtn.classList.add('bg-primary', 'text-white', 'active');
    }

    document.getElementById('myItems').classList.toggle('hidden', subTab !== 'mine');
    document.getElementById('addItem').classList.toggle('hidden', subTab !== 'add');

    if (subTab === 'mine') loadMyItems();
}

// Items Management
async function loadMyItems() {
    try {
        const response = await fetch(`${CONFIG.API_URL}/items`, {
            headers: { 'Authorization': `Bearer ${state.token}` }
        });
        const allItems = await response.json();
        const myItems = allItems.filter(item => item.user_id === state.currentUser.id);
        
        const list = document.getElementById('myItems');
        
        if (myItems.length === 0) {
            list.innerHTML = '<div class="text-center text-gray-400 py-10">You haven\'t added any items yet.</div>';
            return;
        }

        list.innerHTML = myItems.map(item => {
            const imageHtml = item.image_url 
                ? `<img src="${item.image_url}" alt="${item.title}" class="w-full h-40 object-cover rounded-t-xl">`
                : `<div class="w-full h-40 bg-gradient-to-br from-purple-100 to-pink-100 rounded-t-xl flex items-center justify-center">
                    <span class="text-5xl">Image of ${item.category}</span>
                  </div>`;

            return `
                <div class="bg-gray-50 rounded-xl border-2 border-gray-200 mb-4 overflow-hidden">
                    ${imageHtml}
                    <div class="p-4">
                        <h3 class="font-bold text-lg mb-1 text-gray-800">${item.title}</h3>
                        <p class="text-gray-600 text-sm mb-3">${item.description || 'No description'}</p>
                        <div class="flex justify-between items-center">
                            ${item.category ? `<span class="bg-white px-2.5 py-1 rounded-xl text-xs text-primary font-semibold">${item.category}</span>` : '<span></span>'}
                            <button onclick="deleteItem(${item.id})" 
                                    class="bg-red-500 text-white px-3 py-1.5 rounded-md text-sm hover:bg-red-600 transition-colors">
                                Delete
                            </button>
                        </div>
                    </div>
                </div>
            `;
        }).join('');
    } catch (error) {
        console.error('Error loading items:', error);
    }
}

async function addItem(e) {
    e.preventDefault();
    
    const item = {
        title: document.getElementById('title').value,
        description: document.getElementById('description').value,
        category: document.getElementById('category').value,
        image_url: document.getElementById('imageUrl').value
    };

    try {
        const response = await fetch(`${CONFIG.API_URL}/items`, {
            method: 'POST',
            headers: { 
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${state.token}`
            },
            body: JSON.stringify(item)
        });

        if (response.ok) {
            document.getElementById('successMessage').innerHTML = 
                '<div class="bg-green-500 text-white p-3 rounded-lg mb-4 text-center">Item added successfully!</div>';
            e.target.reset();
            setTimeout(() => {
                document.getElementById('successMessage').innerHTML = '';
                showItemsSubTab('mine');
                loadMyItems();
            }, 1500);
        }
    } catch (error) {
        console.error('Error adding item:', error);
    }
}

async function deleteItem(id) {
    if (!confirm('Are you sure you want to delete this item?')) return;

    try {
        await fetch(`${CONFIG.API_URL}/items/${id}`, { 
            method: 'DELETE',
            headers: { 'Authorization': `Bearer ${state.token}` }
        });
        loadMyItems();
    } catch (error) {
        console.error('Error deleting item:', error);
    }
}

// Swipe Functionality
async function loadSwipeItems() {
    try {
        const response = await fetch(`${CONFIG.API_URL}/items?exclude_own=true`, {
            headers: { 'Authorization': `Bearer ${state.token}` }
        });
        state.currentItems = await response.json();
        state.currentIndex = 0;
        displayCurrentItem();
    } catch (error) {
        console.error('Error loading items:', error);
    }
}

function displayCurrentItem() {
    const content = document.getElementById('swipeContent');
    
    if (state.currentIndex >= state.currentItems.length) {
        content.innerHTML = `
            <div class="text-center text-gray-400 py-10">
                <div class="font-semibold text-lg">No more items!</div>
                <div class="text-sm mt-2">Check back later for new items.</div>
            </div>
        `;
        return;
    }

    const item = state.currentItems[state.currentIndex];
    const imageHtml = item.image_url 
        ? `<img src="${item.image_url}" alt="${item.title}" class="w-full h-64 object-cover rounded-xl mb-4">`
        : `<div class="w-full h-64 bg-gradient-to-br from-purple-100 to-pink-100 rounded-xl mb-4 flex items-center justify-center">
            <span class="text-6xl">Image of a ${item.category}</span>
          </div>`;

    content.innerHTML = `
        <div class="mb-6">
            ${imageHtml}
            <div class="text-sm text-gray-500 mb-2">Posted by ${item.owner_name}</div>
            <h2 class="text-3xl font-bold mb-3 text-gray-800">${item.title}</h2>
            ${item.category ? `<span class="inline-block bg-gray-100 px-3 py-1.5 rounded-full text-sm text-gray-600 mb-3">${item.category}</span>` : ''}
            <p class="text-gray-600 leading-relaxed">${item.description || 'No description provided'}</p>
        </div>
        <div class="flex justify-center gap-5 mt-6">
            <button onclick="swipe('left')" 
                    class="w-16 h-16 rounded-full bg-red-500 text-white text-3xl flex items-center justify-center cursor-pointer transition-all hover:scale-110 hover:shadow-xl active:scale-95">
                ✕
            </button>
            <button onclick="swipe('right')" 
                    class="w-16 h-16 rounded-full bg-green-500 text-white text-3xl flex items-center justify-center cursor-pointer transition-all hover:scale-110 hover:shadow-xl active:scale-95">
                ♥
            </button>
        </div>
    `;
}

async function swipe(direction) {
    const item = state.currentItems[state.currentIndex];
    
    try {
        await fetch(`${CONFIG.API_URL}/swipes`, {
            method: 'POST',
            headers: { 
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${state.token}`
            },
            body: JSON.stringify({
                item_id: item.id,
                direction: direction
            })
        });

        state.currentIndex++;
        displayCurrentItem();
    } catch (error) {
        console.error('Error swiping:', error);
    }
}

// Matches & Comments
async function loadMatches() {
    try {
        const response = await fetch(`${CONFIG.API_URL}/matches`, {
            headers: { 'Authorization': `Bearer ${state.token}` }
        });
        const matches = await response.json();
        
        const list = document.getElementById('matchesList');
        
        if (matches.length === 0) {
            list.innerHTML = '<div class="text-center text-gray-400 py-10">No matches yet. Keep swiping!</div>';
            return;
        }

        list.innerHTML = matches.map(match => {
            const isUser1 = match.user1_id === state.currentUser.id;
            const yourItem = isUser1 ? match.item1_title : match.item2_title;
            const theirItem = isUser1 ? match.item2_title : match.item1_title;
            const theirName = isUser1 ? match.user2_name : match.user1_name;

            const commentsHtml = match.comments && match.comments.length > 0
                ? match.comments.map(c => `
                    <div class="bg-white/20 p-2 rounded mb-2">
                        <div class="text-xs font-semibold mb-1">${c.user_name}</div>
                        <div class="text-sm">${c.content}</div>
                    </div>
                  `).join('')
                : '';

            return `
                <div class="bg-gradient-to-r from-green-500 to-green-600 p-5 rounded-xl text-white mb-4">
                    <h3 class="font-bold text-lg mb-3">Match with ${theirName}!</h3>
                    <div class="flex justify-between items-center gap-3 mb-4">
                        <div class="flex-1 text-center">
                            <strong class="block text-sm mb-1">Your Item</strong>
                            <p class="text-sm opacity-90">${yourItem}</p>
                        </div>
                        <div class="text-2xl">⇄</div>
                        <div class="flex-1 text-center">
                            <strong class="block text-sm mb-1">Their Item</strong>
                            <p class="text-sm opacity-90">${theirItem}</p>
                        </div>
                    </div>
                    
                    <div class="border-t border-white/30 pt-3">
                        <div class="mb-2 max-h-32 overflow-y-auto">${commentsHtml || '<div class="text-sm opacity-75">No comments yet</div>'}</div>
                        <form onsubmit="addComment(event, ${match.id})" class="flex gap-2">
                            <input type="text" placeholder="Write a comment..." required
                                class="flex-1 p-2 rounded-lg text-gray-800 text-sm focus:outline-none">
                            <button type="submit" 
                                    class="bg-white text-green-600 px-4 py-2 rounded-lg text-sm font-semibold hover:bg-green-50 transition-colors">
                                Send
                            </button>
                        </form>
                    </div>
                </div>
            `;
        }).join('');
    } catch (error) {
        console.error('Error loading matches:', error);
    }
}

async function addComment(e, matchId) {
    e.preventDefault();
    const input = e.target.querySelector('input');
    const content = input.value;

    try {
        await fetch(`${CONFIG.API_URL}/comments`, {
            method: 'POST',
            headers: { 
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${state.token}`
            },
            body: JSON.stringify({
                match_id: matchId,
                content: content
            })
        });

        input.value = '';
        loadMatches();
    } catch (error) {
        console.error('Error adding comment:', error);
    }
}

// Initialize Application
checkAuth();