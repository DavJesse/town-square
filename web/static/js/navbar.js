export function renderNavBar() {
  // Create nav bar container
  let navBarContainer = document.createElement('div');
  navBarContainer.classList.add('navbar');
  navBarContainer.id = 'navbar_container';
  
  // Create logo link
  let logoHomeLink = document.createElement('a');
  logoHomeLink.href = '/';
  
  let logo = document.createElement('h1');
  logo.textContent = 'real-time forum';
  logo.classList.add('navbar__title');
  
  logoHomeLink.appendChild(logo);
  navBarContainer.appendChild(logoHomeLink);
  
  // Create center navbar cluster container
  let centerCluster = document.createElement('div');
  centerCluster.classList.add('navbar-center');
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
  
  // Create like icon
  let likeIconLink = document.createElement('a');
  likeIconLink.href = '/liked-posts';
  
  let likeIcon = document.createElement('span');
  likeIcon.classList.add('material-symbols-outlined');
  likeIcon.id = 'like_icon';
  likeIcon.textContent = 'thumb_up';
  likeIconLink.appendChild(likeIcon);
  
  centerCluster.appendChild(likeIconLink);
  
  // Create chat icon
  let chatIconLink = document.createElement('a');
  chatIconLink.href = '/chat';
  
  let chatIcon = document.createElement('span');
  chatIcon.classList.add('material-symbols-outlined');
  chatIcon.id = 'chat_icon';
  chatIcon.textContent = 'forum';
  chatIconLink.appendChild(chatIcon);
  
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
  rightCluster.classList.add('navbar-right');
  rightCluster.id = 'navbar_right_cluster';
  
  // Create search bar form
  let searchBar = document.createElement('form');
  searchBar.method = 'GET';
  searchBar.action = '/search';
  searchBar.classList.add('navbar__search-form');
  searchBar.id = 'search_bar';

   // Create right navbar cluster container
   let searchCluster = document.createElement('div');
   searchCluster.classList.add('search-cluster');
   searchCluster.id = 'search_cluster';
  
  // Create search bar input field
  let searchInput = document.createElement('input');
  searchInput.type = 'text';
  searchInput.name = 'q';
  searchInput.id = 'search_input';
  searchInput.placeholder = 'Search...';    
  searchInput.classList.add('navbar__search');
  searchCluster.appendChild(searchInput);
  
  // Create search button
  let searchButton = document.createElement('button');
  searchButton.type = 'submit';
  searchButton.id = 'search_btn';
  
  let searchIcon = document.createElement('span');
  searchIcon.classList.add('material-symbols-outlined');
  searchIcon.id ='search_icon';
  searchIcon.textContent = 'search';
  
  searchButton.appendChild(searchIcon);
  searchCluster.appendChild(searchButton);
  searchBar.appendChild(searchCluster);
  rightCluster.appendChild(searchBar);
  navBarContainer.appendChild(rightCluster);

  // Append nav bar container to the app div
  let app = document.getElementById('app')
  app.appendChild(navBarContainer);
}
