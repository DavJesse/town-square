import { renderProfileMenu } from '/static/js/mobile_profile_menu.js';
import { renderLogoutButton } from '/static/js/logout_button.js';

export function renderNavBar(userData) {
  // Create nav bar container
  let navBarContainer = document.createElement('div');
  navBarContainer.classList.add('navbar');
  navBarContainer.id = 'navbar_container';
  
  // Create logo link
  let leftCluster = document.createElement('div');
  let logoHomeLink = document.createElement('a');
  leftCluster.classList.add('navbar-clusters');
  leftCluster.id = 'navbar_left_cluster';
  logoHomeLink.href = '/';
  
  let logo = document.createElement('h1');
  logo.textContent = window.innerWidth < 1025 ? 'forum' : 'real-time forum';
  logo.classList.add('navbar__title');

  // Detect screen resizing and adjust logo accordinly
  window.addEventListener('resize', () => {
    logo.textContent = window.innerWidth < 1025 ? 'forum' : 'real-time forum';
  });
  
  logoHomeLink.appendChild(logo);
  leftCluster.appendChild(logoHomeLink);
  navBarContainer.appendChild(leftCluster);
  
  // Create center navbar cluster container
  let centerCluster = document.createElement('div');
  centerCluster.classList.add('navbar-clusters');
  centerCluster.id = 'navbar_center_cluster';
  
  // Create home icon
  let homeIconLink = document.createElement('a');
  homeIconLink.href = '/';
  
  let homeIcon = document.createElement('span');
  homeIcon.classList.add('material-symbols-outlined');
  homeIcon.id = 'home_icon';
  homeIcon.textContent = 'home';
  homeIconLink.appendChild(homeIcon);
  
  centerCluster.appendChild(homeIconLink);
  
  // Create chat icon
  let chatIconLink = document.createElement('a');
  chatIconLink.href = '/chat';

  let chatContainer = document.createElement('div');
  chatContainer.id = 'chat_container';
  
  let chatIcon = document.createElement('span');
  chatIcon.classList.add('material-symbols-outlined');
  chatIcon.id = 'chat_icon';
  chatIcon.textContent = 'forum';
  chatContainer.appendChild(chatIcon);
  chatIconLink.appendChild(chatContainer);
  
  centerCluster.appendChild(chatIconLink);
  
  // Create notifications icon
  let notificationIconLink = document.createElement('a');
  notificationIconLink.href = '/notifications';
  
  let notificationIcon = document.createElement('span');
  notificationIcon.classList.add('material-symbols-outlined');
  notificationIcon.id = 'notification_icon';
  notificationIcon.textContent = 'notifications';
  notificationIconLink.appendChild(notificationIcon);
  
  centerCluster.appendChild(notificationIconLink);
  navBarContainer.appendChild(centerCluster);
  
  // Create right navbar cluster container
  let rightCluster = document.createElement('div');
  rightCluster.classList.add('navbar-clusters');
  rightCluster.id = 'navbar_right_cluster';

  // Add the logout button
  const logoutButton = renderLogoutButton();
  rightCluster.appendChild(logoutButton);

  navBarContainer.appendChild(rightCluster);

  // Append nav bar container to the app div
  let app = document.getElementById('app')
  app.appendChild(navBarContainer);

  // Render profile menu on mobile
  if (window.innerWidth <= 540) {
    renderProfileMenu(userData);
  }

  // Dynamically detect mobile before rendering profile menu
  window.addEventListener('resize', () => {
    if (window.innerWidth <= 540 && !document.getElementById('profile_menu_button')) {
      renderProfileMenu(userData);
    }
  });
}

 // Create search bar form
  // let searchBar = document.createElement('form');
  // searchBar.method = 'GET';
  // searchBar.action = '/search';
  // searchBar.classList.add('navbar__search-form');
  // searchBar.id = 'search_bar';

   // Create right navbar cluster container
  //  let searchCluster = document.createElement('div');
  //  searchCluster.classList.add('search-cluster');
  //  searchCluster.id = 'search_cluster';
  
  // Create search bar input field
  // let searchInput = document.createElement('input');
  // searchInput.type = 'text';
  // searchInput.name = 'q';
  // searchInput.id = 'search_input';
  // searchInput.placeholder = 'Search...';    
  // searchInput.classList.add('navbar__search');
  // searchCluster.appendChild(searchInput);
  
  // Create search button
  // let searchButton = document.createElement('button');
  // searchButton.type = 'submit';
  // searchButton.id = 'search_btn';
  
  // let searchIcon = document.createElement('span');
  // searchIcon.classList.add('material-symbols-outlined');
  // searchIcon.id ='search_icon';
  // searchIcon.textContent = 'search';
  
  // searchButton.appendChild(searchIcon);
  // searchCluster.appendChild(searchButton);
  // searchBar.appendChild(searchCluster);
  // rightCluster.appendChild(searchBar);