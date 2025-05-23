import { populateCategories } from "/static/js/populate_categories.js";
import { fetchAllUsers } from '/static/js/chat.js';
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
    
    // Fetch categories and populate when ready
    fetchCategories().then(categories => {
        populateCategories(categories, categoriesList.id);
    }).catch(error => {
        console.error('Error fetching categories:', error);
    });

    // Create online users card
    let onlineUsersCard = document.createElement('div');
    let onlineUsersTitle = document.createElement('h3');
    let onlineUsersContent = document.createElement('div');
    onlineUsersCard.id = 'online_users_card'
    onlineUsersTitle.id = 'online_users_title';
    onlineUsersContent.id = 'online_users_content';
    onlineUsersTitle.textContent = 'Who\'s online?';
    onlineUsersCard.appendChild(onlineUsersTitle);
    smokeScreen.appendChild(onlineUsersCard);
    onlineUsersCard.appendChild(onlineUsersContent);

    fetchAllUsers();
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
