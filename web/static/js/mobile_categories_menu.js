import { populateCategories } from "/static/js/populate_categories.js";
import { fetchAllUsers, populateOnlineUsersList } from '/static/js/chat.js';
import { navigateTo } from '/static/js/routes.js';
import { renderErrorPage } from '/static/js/error.js';

export function renderMobileCategoriesMenu() {
    // Grab app
    let app = document.getElementById('app');
    
    // Create smoke screen
    let smokeScreen = document.createElement('div');
    smokeScreen.id = 'smoke_screen';
    smokeScreen.style.flexDirection = 'column';
    app.appendChild(smokeScreen);

    // Make hamburger menu invisible
    document.getElementById('hamburger_menu').style.display = 'none';

    // Create close button and add related event listener
    let closeButton = document.createElement('button');
    closeButton.id = 'mobile_categories_close_button';
    closeButton.textContent = 'X';
    closeButton.addEventListener('click', () => {
        smokeScreen.remove();
        document.getElementById('hamburger_menu').style.display = 'block';
    });
    smokeScreen.appendChild(closeButton);
    
    // Create categories container
    let categoriesContainer = document.createElement('div');
    categoriesContainer.id = 'mobile_categories_container';
    smokeScreen.appendChild(categoriesContainer);

    // Create categories title
    let categoriesTitle = document.createElement('h3');
    categoriesTitle.id = 'categories_title';
    categoriesTitle.textContent = 'Categories';
    categoriesContainer.appendChild(categoriesTitle);

    // Create categories list
    let categoriesList = document.createElement('div');
    categoriesList.id = 'mobile_category_content_container';
    categoriesContainer.appendChild(categoriesList);

    // Create online users card
    let onlineUsersCard = document.createElement('div');
    onlineUsersCard.id = 'mobile_online_users_card';
    
    let onlineUsersTitle = document.createElement('h3');
    onlineUsersTitle.id = 'online_users_title';
    onlineUsersTitle.textContent = 'Who\'s online?';
    onlineUsersCard.appendChild(onlineUsersTitle);
    
    let onlineUsersContent = document.createElement('div');
    onlineUsersContent.id = 'online_users_content';
    onlineUsersCard.appendChild(onlineUsersContent);
    
    smokeScreen.appendChild(onlineUsersCard);

    // First fetch categories
    fetchCategories().then(categories => {
        populateCategories(categories, categoriesList.id);
        // Then fetch users after categories are loaded
        return fetchAllUsers();
    }).then(users => {
        if (users && users.length > 0) {
            let onlineUsersCard = document.getElementById('mobile_online_users_card');
            let onlineUsersContent = onlineUsersCard.querySelector('#online_users_content');
            if (onlineUsersContent) {
                populateOnlineUsersList(users, onlineUsersContent);
            }
        }
    }).catch(error => {
        console.error('Error:', error);
    });
}

function fetchCategories() {
    return fetch('/api/index-data', {
        method: 'GET',
        credentials: 'include',
        headers: {
            'Accept': 'application/json'
        }
    })
    .then(response => {
        if (!response.ok) {
            if (response.status === 401) {
                navigateTo('/login');
                return Promise.reject('Unauthorized');
            } else {
                renderErrorPage(response.statusText, response.status);
                throw new Error(`HTTP error! status: ${response.status}`);
            }
        }
        return response.json();
    })
    .then(data => {
        return data.categories;
    });
}
